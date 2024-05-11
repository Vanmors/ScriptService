package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func NewClient(dbname, username, password, host, port string) (*sql.DB, error) {
	//connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, dbname)
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		username,
		password,
		host,
		port,
		dbname)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
