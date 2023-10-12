package main

import (
	"fmt"
	"funcopts/db"
	"log"
)

func main() {
	conn, err := db.Open("https://127.0.0.1:8080", db.WithCache())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conn)
}
