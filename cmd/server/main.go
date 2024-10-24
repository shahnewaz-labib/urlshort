package main

import (
	"net/http"
	"os"
	"shahnewaz-labib/urlshort/internal/logger"
	"shahnewaz-labib/urlshort/internal/shortener"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}

	logger.Init()
	defer logger.Sync()

	repo, err := shortener.NewRepository()
	if err != nil {
		logger.Fatal("Failed to create repository", zap.Error(err))
	}

	http.HandleFunc("/shorten", shortener.ShortenURLHandler(repo))
	http.HandleFunc("/resolve", shortener.ResolveURLHandler(repo))
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/", shortener.RedirectHandler(repo))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server is running", zap.String("port", port))
	logger.Fatal("Server stopped", zap.Error(http.ListenAndServe(":"+port, nil)))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
