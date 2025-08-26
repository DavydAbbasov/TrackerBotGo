package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // драйвер PostgreSQL для database/sql
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("open:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("ping:", err)
	}

	var now string
	if err := db.QueryRow(`select to_char(now(),'YYYY-MM-DD HH24:MI:SS')`).Scan(&now); err != nil {
		log.Fatal("query:", err)
	}
	fmt.Println("DB OK. Server time:", now)
}
