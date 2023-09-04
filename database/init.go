package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.elastic.co/apm/module/apmsql/v2"
	_ "go.elastic.co/apm/module/apmsql/v2/pq"
	"log"
	"splunk-test/config"
)

func Init() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Get("TEST_DB_HOST"),
		config.Get("TEST_DB_USER"),
		config.Get("TEST_DB_PASSWORD"),
		config.Get("TEST_DB_NAME"),
		config.Get("TEST_DB_PORT"))

	db, err := apmsql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return sqlx.NewDb(db, "postgres"), nil
}
