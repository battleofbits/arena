package main

import (
	"github.com/kevinburke/arena"
	"log"
	"net/http"
)

func main() {
	router := arena.DoServer()
	log.Fatal(http.ListenAndServe(":8080", router))
}
