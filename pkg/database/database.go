package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // _ to make sure there's no ref required
)

var (
	Ctx   context.Context
	DBCon *sql.DB // Create a globla variable to reach the DBConnection
)

// Init DB connection and store to global variable.
func Init() {
	db, err := sql.Open("mysql", "root:secret@tcp(127.0.0.1:3306)/imperial") // Should make this a flag to insert when running for details...
	if err != nil {
		log.Fatal(err.Error())
	}
	DBCon = db
	initTables() // Start creation of tables, should just make a flag if new tables needs to be created like in Laravel
}

// Creates the tables at start up
// TO DO: Add a flag to initiate this
func initTables() {
	stmt, err := DBCon.Prepare("CREATE Table ships(id int NOT NULL AUTO_INCREMENT, name varchar(50), class varchar(50), crew int, image varchar(300), value float, status varchar(50), armament varchar(30), PRIMARY KEY(id));")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close() // Always close resource after no longer usage else memory leak..
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Table created successfully")
	}
}
