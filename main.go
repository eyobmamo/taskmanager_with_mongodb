package main

import (
	"fmt"
	"log"
	"mongodb/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	fmt.Println("Task Manager API")
	r := router.NewRouter()
	if r != nil {
		r.Run()
	} else {
		log.Fatal("Failed to start server")
	}

}
