package main

import (
	"log"
	"net/http"

	"github.com/idlistic/go-backend-api-sample/internal/router"
)

func main() {
	r, cleanup, err := router.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cleanup() }()

	log.Println("🚀 Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
