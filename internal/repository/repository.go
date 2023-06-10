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
	result = make(map[string]uint64)
	rows, err := r.db.Query(`SELECT * FROM statistics`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var temp models.Stats
		err = rows.Scan(&temp.Endpoint, &temp.Counter)
		if err != nil {
			return
		}
		result[temp.Endpoint] = temp.Counter
	}
	return
}
