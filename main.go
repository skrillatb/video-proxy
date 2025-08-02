package main

import (
	"log"
	"net/http"

	"vproxy/handlers"
)

func main() {
	http.HandleFunc("/video", handlers.MP4)


	port := ":8080"
	log.Println("Serveur lancé sur", port)
	log.Fatal(http.ListenAndServe(port, nil))
}