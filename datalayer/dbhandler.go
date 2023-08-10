package datalayer

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type SQLhandler struct {
	db *sql.DB
}

func CreateDbConnection(connstring string) (*SQLhandler, error) {
	db, err := sql.Open("mysql", connstring)
	if err != nil {
		return nil, err
	}

	return &SQLhandler{db: db}, nil
}
