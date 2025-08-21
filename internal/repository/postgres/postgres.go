package postgres

import (
	"github.com/IceMAN2377/thfc/internal/models"
	"github.com/IceMAN2377/thfc/internal/repository"
	"github.com/jmoiron/sqlx"
	"log"
)

type postgres struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) repository.Repository {
	return &postgres{
		db: db,
	}
}

func (p *postgres) PostText(record *models.Record) error {

	stmt, err := p.db.Preparex(`INSERT INTO records (title, content) VALUES ($1, $2)`)
	if err != nil {
		log.Printf("DB error: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.Title, record.Content)
	if err != nil {
		log.Printf("some error: %v", err)
		return err
	}
	return nil
}

func (p *postgres) GetByTitle(title string) (*models.Record, error) {

	stmt, err := p.db.Preparex(`SELECT title, content FROM records WHERE title=$1`)
	if err != nil {
		log.Printf("Error preparing stmt:%v", err)

		return nil, err
	}

	var record models.Record

	if err := stmt.Get(&record, title); err != nil {
		log.Printf("error retrieving the record:%v", err)

		return nil, err
	}

	return &record, nil
}
