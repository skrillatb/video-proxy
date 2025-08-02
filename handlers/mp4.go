package handlers

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{
	Timeout: 20 * time.Second, 
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 0, 
		}).DialContext,
		MaxIdleConns:        10,
		IdleConnTimeout:     5 * time.Second,
		DisableCompression:  true, 
		ForceAttemptHTTP2:   false,
		MaxIdleConnsPerHost: 2,
	},
}

func MP4(w http.ResponseWriter, r *http.Request) {
	videoURL := r.URL.Query().Get("url")
	if videoURL == "" {
		http.Error(w, "Paramètre 'url' manquant", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(videoURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		http.Error(w, "URL invalide", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(http.MethodGet, videoURL, nil)
	if err != nil {
		http.Error(w, "Erreur création requête", http.StatusInternalServerError)
		return
	}

	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}
	// req.Header.Set("User-Agent", r.UserAgent()) 


	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, "Erreur accès distant", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	io.CopyBuffer(w, resp.Body, make([]byte, 32*1024))
}

func copyHeaders(dst, src http.Header) {
	for k, v := range src {
		for _, val := range v {
			if k == "Transfer-Encoding" {
				continue
			}
			dst.Add(k, val)
		}
	}
}
