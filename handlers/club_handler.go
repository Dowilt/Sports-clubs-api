package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"sports-clubs-api/db"
	"sports-clubs-api/models"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

// GetClubs обрабатывает GET-запрос для получения списка клубов с фильтрацией и сортировкой
func GetClubs(c echo.Context) error {
	query := c.QueryParam("q")
	titlesMin := c.QueryParam("titles_min")
	titlesMax := c.QueryParam("titles_max")
	ageMin := c.QueryParam("age_min")
	ageMax := c.QueryParam("age_max")
	sortField := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")

	var whereClauses []string
	var args []interface{}
	argID := 1

	if query != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(c.name ILIKE $%d OR ci.name ILIKE $%d)", argID, argID))
		args = append(args, "%"+query+"%")
		argID++
	}

	if titlesMin != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.titles_count >= $%d", argID))
		args = append(args, titlesMin)
		argID++
	}

	if titlesMax != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.titles_count <= $%d", argID))
		args = append(args, titlesMax)
		argID++
	}

	if ageMin != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.avg_age >= $%d", argID))
		args = append(args, ageMin)
		argID++
	}

	if ageMax != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.avg_age <= $%d", argID))
		args = append(args, ageMax)
		argID++
	}

	baseQuery := `
        SELECT c.id, c.name, ci.name, c.titles_count, c.avg_age 
        FROM clubs c 
        JOIN cities ci ON c.city_id = ci.id
    `

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if sortField != "" && sortOrder != "" {
		validOrder := "ASC"
		if sortOrder == "desc" {
			validOrder = "DESC"
		}

		switch sortField {
		case "titles_count":
			baseQuery += fmt.Sprintf(" ORDER BY c.titles_count %s", validOrder)
		case "avg_age":
			baseQuery += fmt.Sprintf(" ORDER BY c.avg_age %s", validOrder)
		}
	}

	conn, err := db.ConnectDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), baseQuery, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Query execution failed"})
	}
	defer rows.Close()

	var clubs []models.Club
	for rows.Next() {
		var club models.Club
		if err := rows.Scan(&club.ID, &club.Name, &club.City, &club.Titles, &club.AvgAge); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Row scan failed"})
		}
		clubs = append(clubs, club)
	}

	return c.JSON(http.StatusOK, clubs)
}

// getOrCreateCityID находит ID города или создаёт новый, если его нет
func getOrCreateCityID(ctx context.Context, conn *pgx.Conn, cityName string) (int, error) {
	var cityID int
	err := conn.QueryRow(ctx, "SELECT id FROM cities WHERE name ILIKE $1", cityName).Scan(&cityID)
	if err == nil {
		return cityID, nil
	}

	// Если такого города нет, создаём
	err = conn.QueryRow(ctx, "INSERT INTO cities (name) VALUES ($1) RETURNING id", cityName).Scan(&cityID)
	if err != nil {
		return 0, err
	}
	return cityID, nil
}

// CreateClub обрабатывает POST-запрос для добавления нового клуба
func CreateClub(c echo.Context) error {
	log.Printf("Received POST /clubs request")

	var club models.Club
	if err := c.Bind(&club); err != nil {
		log.Printf("Error binding request body: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	log.Printf("Received club data: %+v", club)

	// Валидация данных
	if club.Titles < 0 || club.AvgAge <= 0 || club.AvgAge >= 100 {
		log.Printf("Invalid club data: Titles=%d, AvgAge=%d", club.Titles, club.AvgAge)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid club data"})
	}

	conn, err := db.ConnectDB()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}
	defer conn.Close(context.Background())

	cityID, err := getOrCreateCityID(context.Background(), conn, club.City)
	if err != nil {
		log.Printf("Error getting or creating city: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process city"})
	}

	log.Printf("City ID for '%s': %d", club.City, cityID)

	_, err = conn.Exec(context.Background(),
		"INSERT INTO clubs (name, city_id, titles_count, avg_age) VALUES ($1, $2, $3, $4)",
		club.Name, cityID, club.Titles, club.AvgAge,
	)
	if err != nil {
		log.Printf("Error inserting club into database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create club"})
	}

	log.Printf("Club created successfully: %+v", club)
	return c.JSON(http.StatusCreated, map[string]string{"message": "Клуб успешно создан"})
}

// UpdateClub обрабатывает PUT-запрос для обновления существующего клуба
func UpdateClub(c echo.Context) error {
	id := c.Param("id")

	var req models.UpdateClubRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	conn, err := db.ConnectDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}
	defer conn.Close(context.Background())

	// Получаем текущие данные
	var club models.Club
	err = conn.QueryRow(context.Background(),
		"SELECT name, city_id, titles_count, avg_age FROM clubs WHERE id = $1",
		id).Scan(&club.Name, &club.CityID, &club.Titles, &club.AvgAge)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch club data"})
	}

	// Обновляем только те поля, которые были переданы
	if req.Name != nil {
		club.Name = *req.Name
	}
	if req.City != nil {
		cityID, err := getOrCreateCityID(context.Background(), conn, *req.City)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process city"})
		}
		club.CityID = cityID
	}
	if req.Titles != nil {
		if *req.Titles < 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Titles cannot be negative"})
		}
		club.Titles = *req.Titles
	}
	if req.AvgAge != nil {
		if *req.AvgAge <= 0 || *req.AvgAge >= 100 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Avg age must be between 1 and 99"})
		}
		club.AvgAge = *req.AvgAge
	}

	// Обновляем запись в БД
	_, err = conn.Exec(context.Background(),
		"UPDATE clubs SET name = $1, city_id = $2, titles_count = $3, avg_age = $4 WHERE id = $5",
		club.Name, club.CityID, club.Titles, club.AvgAge, id,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update club"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Клуб успешно обновлен"})
}

// DeleteClub обрабатывает DELETE-запрос для удаления клуба
func DeleteClub(c echo.Context) error {
	id := c.Param("id")

	conn, err := db.ConnectDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "DELETE FROM clubs WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete club"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Клуб успешно удален"})
}
