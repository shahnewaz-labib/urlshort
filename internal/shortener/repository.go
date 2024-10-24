package shortener

import (
	"fmt"
	"os"
	"shahnewaz-labib/urlshort/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type URL struct {
	ID          uint   `gorm:"primaryKey"`
	ShortUrl    string `gorm:"uniqueIndex"`
	OriginalUrl string
}

type Repository struct {
	db *gorm.DB
}

func NewRepository() (*Repository, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&URL{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	logger.Info("Connected to PostgreSQL and AutoMigrate successfully")

	return &Repository{db: db}, nil
}

func (repo *Repository) SaveURL(shortURL string, originalURL string) error {
	url := &URL{
		ShortUrl:    shortURL,
		OriginalUrl: originalURL,
	}
	return repo.db.Create(url).Error
}

func (repo *Repository) ResolveURL(shortURL string) (string, error) {
	var url URL

	err := repo.db.First(&url, "short_url = ?", shortURL).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("short URL not found")
		}
		return "", fmt.Errorf("failed to retrieve original URL: %w", err)
	}

	return url.OriginalUrl, nil
}

func (repo *Repository) FindShortURL(originalURL string) (string, error) {
	var url URL
	err := repo.db.First(&url, "original_url = ?", originalURL).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", fmt.Errorf("failed to find short URL: %w", err)
	}
	return url.ShortUrl, nil
}
