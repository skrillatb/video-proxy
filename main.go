package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"vproxy/handlers"

	"github.com/joho/godotenv"
)

var corsOrigins []string

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env non trouvé ou non chargé")
	}

	prodFrontend := os.Getenv("PROD_FRONTEND")
	devFrontend := os.Getenv("DEV_FRONTEND")

	if prodFrontend != "" {
		corsOrigins = append(corsOrigins, prodFrontend)
	}
	if devFrontend != "" {
		corsOrigins = append(corsOrigins, devFrontend)
	}

	log.Printf("CORS Origins: %v", corsOrigins)
}

func isOriginAllowed(origin string) bool {
	for _, allowed := range corsOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	if origin == "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Range, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")
		return true
	}

	if isOriginAllowed(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Range, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")
		return true
	}

	log.Printf("Origin refusée: %s", origin)
	log.Printf("Origines autorisées: %v", corsOrigins)
	return false
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if !setCORSHeaders(w, r) {
		http.Error(w, "Origin non autorisée", http.StatusForbidden)
		return
	}
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func main() {
	loadEnv()
	
	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		if !setCORSHeaders(w, r) {
			http.Error(w, "Origin non autorisée", http.StatusForbidden)
			return
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		handlers.MP4(w, r)
	})
	
	http.HandleFunc("/uqload", func(w http.ResponseWriter, r *http.Request) {
		if !setCORSHeaders(w, r) {
			http.Error(w, "Origin non autorisée", http.StatusForbidden)
			return
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		handlers.UqloadHandler(w, r)
	})
	
	http.HandleFunc("/health", healthHandler)

	port := ":8080"
	log.Println("Serveur lancé sur", port)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erreur serveur: %v", err)
	}
}