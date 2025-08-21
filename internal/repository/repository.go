package repository

import "github.com/IceMAN2377/thfc/internal/models"

type Repository interface {
	PostText(record *models.Record) error
	GetByTitle(title string) (*models.Record, error)
}
