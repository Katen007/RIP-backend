package main

import (
	"log"

	"lab1_rip/internal/api"
)

func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Applocation terminated!")
}
