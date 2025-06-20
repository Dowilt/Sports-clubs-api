package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func ConnectDB() (*pgx.Conn, error) {
	connStr := "user=myuser password=mypassword dbname=sports_clubs sslmode=disable"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}
