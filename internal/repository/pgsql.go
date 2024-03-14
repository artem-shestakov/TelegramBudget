package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewPgDb(logger *logrus.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			"localhost", 5432, "postgres", "postgres", "budget"),
	)
	if err != nil {
		panic(err)
	}

	logger.Infof("Successfully connected to DB")
	return db, nil
}
