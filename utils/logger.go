package utils

import (
	"log"
	"net/http"
	"time"
)

func LogRequest(r *http.Request, msg string) {
	fullURL := r.Host + r.URL.RequestURI()

	log.Printf("[%s] %s\n→ URL appelée: %s\n→ IP: %s\n→ Origin: %s\n→ UA: %s",
		time.Now().Format("2006-01-02 15:04:05"),
		msg,
		fullURL,
		r.RemoteAddr,
		r.Header.Get("Origin"),
		r.Header.Get("User-Agent"),
	)
}