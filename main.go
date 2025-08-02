package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"vproxy/handlers"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/video", handlers.MP4)
	
	http.HandleFunc("/health", healthHandler)

	port := ":8080"
	log.Println("Serveur lancé sur", port)
	log.Fatal(http.ListenAndServe(port, nil))
}