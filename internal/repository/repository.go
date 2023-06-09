package repository

import (
	"database/sql"
	"task/internal/models"
)

type Repository interface {
	SaveStats(specific, min, max uint64) error
	GetStats() (result map[string]uint64, err error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) SaveStats(specific, min, max uint64) error {
	_, err := r.db.Exec(
		`
		INSERT OR REPLACE INTO statistics(endpoint,counter) 
		VALUES("specific",?),("min",?),("max",?)
		`,
		specific,
		min,
		max,
	)
	return err
}

func (r *repository) GetStats() (result map[string]uint64, err error) {
	var stats []models.Stats
	data, err := r.db.Query(`SELECT (endpoint,counter) FROM statistics`)
	if err != nil {
		return
	}
	err = data.Scan(&stats)
	if err != nil {
		return
	}
	for _, content := range stats {
		result[content.Endpoint] = content.Counter
	}
	return
}
