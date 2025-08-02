package handlers

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func UqloadHandler(w http.ResponseWriter, r *http.Request) {
	embedURL := r.URL.Query().Get("url")
	
	if embedURL == "" || !strings.Contains(embedURL, "uqload.cx/embed-") {
		http.Error(w, "URL non valide", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", embedURL, nil)
	if err != nil {
		http.Error(w, "Erreur création requête", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erreur requête embed: %v", err)
		http.Error(w, "Erreur lors de la récupération de la page", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status code: %d", resp.StatusCode)
		http.Error(w, "Erreur lors de la récupération de la page", http.StatusBadGateway)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erreur lecture contenu", http.StatusInternalServerError)
		return
	}

	bodyString := string(bodyBytes)

	videoURL := extractVideoURL(bodyString)
	if videoURL == "" {
		log.Printf("Impossible d'extraire l'URL vidéo de: %s", embedURL)
		http.Error(w, "Impossible d'obtenir le lien vidéo", http.StatusNotFound)
		return
	}

	log.Printf("✅ URL vidéo extraite: %s", videoURL)

	proxyVideo(w, r, videoURL, embedURL)
}

func extractVideoURL(html string) string {
	patterns := []string{
		`sources:\s*\["([^"]+)"\]`,                    
		`sources:\s*\[\s*"([^"]+)"\s*\]`,          
		`file:\s*"([^"]+\.mp4[^"]*)"`,        
		`src:\s*"([^"]+\.mp4[^"]*)"`,         
		`source\s*src="([^"]+\.mp4[^"]*)"`,           
		`"file":\s*"([^"]+\.mp4[^"]*)"`,        
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 && matches[1] != "" {
			videoURL := matches[1]
			videoURL = strings.ReplaceAll(videoURL, "\\", "")
			
			if strings.HasPrefix(videoURL, "http") && strings.Contains(videoURL, ".mp4") {
				return videoURL
			}
		}
	}

	return ""
}

func proxyVideo(w http.ResponseWriter, r *http.Request, videoURL, referer string) {

	req, err := http.NewRequest("GET", videoURL, nil)
	if err != nil {
		http.Error(w, "Erreur création requête vidéo", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "video/webm,video/ogg,video/*;q=0.9,application/ogg;q=0.7,audio/*;q=0.6,*/*;q=0.5")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Connection", "keep-alive")

	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableCompression: true, 
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erreur requête vidéo: %v", err)
		http.Error(w, "Erreur lors de la récupération de la vidéo", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "video/mp4")
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Erreur streaming: %v", err)
	}
}