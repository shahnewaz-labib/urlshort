package shortener

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortURL(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func ShortenURL(repo *Repository, originalURL string) (string, error) {
	existingShortURL, err := repo.FindShortURL(originalURL)
	if err != nil {
		return "", err
	}
	if existingShortURL != "" {
		return existingShortURL, nil
	}

	shortURL := generateShortURL(8)

	err = repo.SaveURL(shortURL, originalURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func ResolveURL(repo *Repository, shortURL string) (string, error) {
	originalURL, err := repo.ResolveURL(shortURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}
