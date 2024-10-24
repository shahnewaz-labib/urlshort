package shortener

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"shahnewaz-labib/urlshort/internal/logger"

	"go.uber.org/zap"
)

func ShortenURLHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var input struct {
			URL string `json:"url"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&input); err != nil {
			logger.Error("Invalid JSON input", zap.Error(err))
			http.Error(w, "Invalid JSON input", http.StatusBadRequest)
			return
		}

		shortURL, err := ShortenURL(repo, input.URL)
		if err != nil {
			logger.Error("Failed to shorten URL", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		domainName := os.Getenv("DOMAIN_NAME")
		if domainName == "" {
			domainName = "http://localhost:8080"
		}

		fullShortURL := fmt.Sprintf("%s/%s", domainName, shortURL)

		response := struct {
			ShortURL string `json:"short_url"`
		}{
			ShortURL: fullShortURL,
		}
		json.NewEncoder(w).Encode(response)
		logger.Info("URL shortened successfully", zap.String("original_url", input.URL), zap.String("short_url", fullShortURL))
	}
}

func ResolveURLHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var input struct {
			URL string `json:"url"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&input); err != nil {
			logger.Error("Invalid JSON input", zap.Error(err))
			http.Error(w, "Invalid JSON input", http.StatusBadRequest)
			return
		}

		if input.URL == "" {
			logger.Error("Missing 'url' in JSON input")
			http.Error(w, "Missing 'url' in JSON input", http.StatusBadRequest)
			return
		}

		originalURL, err := ResolveURL(repo, input.URL)
		if err != nil {
			logger.Error("Failed to resolve URL", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			OriginalURL string `json:"original_url"`
		}{
			OriginalURL: originalURL,
		}
		json.NewEncoder(w).Encode(response)
		logger.Info("URL resolved successfully", zap.String("short_url", input.URL), zap.String("original_url", originalURL))
	}
}

func RedirectHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.URL.Path[1:] // Remove the leading slash
		originalURL, err := ResolveURL(repo, shortURL)
		if err != nil {
			logger.Error("Failed to resolve URL", zap.Error(err))
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
		logger.Info("URL redirected successfully", zap.String("short_url", shortURL), zap.String("original_url", originalURL))
	}
}
