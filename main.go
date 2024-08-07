package main

import (
	"fmt"
	"log"
	"mongodb/router"
)

func main() {
	fmt.Println("Task Manager API")
	r := router.NewRouter()
	if r != nil {
		r.Run()
	} else {
		log.Fatal("Failed to start server")
	}

}
