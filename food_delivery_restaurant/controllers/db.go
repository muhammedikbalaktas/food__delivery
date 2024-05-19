package controllers

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var dbInfo = "root:<YOUR_PASSWORD>@tcp(localhost:3306)/delivery_restaurant"

func createDb() (*sql.DB, error) {

	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		fmt.Println("error on init database")
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}
