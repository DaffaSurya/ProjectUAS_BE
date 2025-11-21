package config
import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectPG() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=12345678 dbname=project_uas sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database unreachable: ", err)
	}

	fmt.Println("Database connected ✅")
	return db
}