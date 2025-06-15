package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"sports-clubs-api/db"
	"sports-clubs-api/models"

	"github.com/labstack/echo/v4"
)

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
