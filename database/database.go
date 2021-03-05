package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	// "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("vim-go")
}

func InitDatabase() (db sql.DB) {

	log.Println("Creating database.db...")
	// dbfile, err := os.Create("sqlite-database.db") // Create SQLite file
	dbfile, err := os.OpenFile("sqlite.db", os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	dbfile.Close()

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                            // Defer Closing the database
	createTable(sqliteDatabase)                             // Create Database Tables

	return *sqliteDatabase
}

func createTable(db *sql.DB) {
	createHostTableSQL := `CREATE TABLE hosts IF NOT EXISTS (
		"host" TEXT NOT NULL PRIMARY KEY,
		"source" TEXT,
		"date" DATETIME,
	  );`

	statement, err := db.Prepare(createHostTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("host table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertHost(db *sql.DB, host string, source string, date time.Time) {
	insertHostSQL := `INSERT INTO hosts(host, source, date) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertHostSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(host, source, date)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func getHost(db *sql.DB) {
	row, err := db.Query("SELECT host FROM hosts WHERE host = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
}
