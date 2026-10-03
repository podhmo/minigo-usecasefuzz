package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Println("open:", err)
		return
	}
	defer db.Close()
	var n int
	if err := db.QueryRow("select 1+1").Scan(&n); err != nil {
		fmt.Println("query:", err)
		return
	}
	fmt.Println("result:", n)
}
