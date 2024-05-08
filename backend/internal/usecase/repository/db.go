package repository

import (
	"database/sql"
	"fmt"

	"github.com/emochka2007/block-accounting/internal/pkg/config"

	_ "github.com/lib/pq"
)

func ProvideDatabaseConnection(c config.Config) (*sql.DB, func(), error) {
	sslmode := "disable"
	if c.DB.EnableSSL {
		sslmode = "enable"
	}

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=%s",
		c.DB.User, c.DB.Secret, c.DB.Host, c.DB.Database, sslmode,
	)

	fmt.Println(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, func() {}, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, func() {
		db.Close()
	}, nil
}
