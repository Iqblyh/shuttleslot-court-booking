package repository

import (
	"database/sql"
	"math"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"time"
)

type CourtRepository interface {
	Create(payload model.Court) (model.Court, error)
	FindAll(page int, size int) ([]model.Court, dto.Paginate, error)
	FindById(id string) (model.Court, error)
	Update(id string, payload model.Court) (model.Court, error)
	Deleted(id string) error
}

type courtRepository struct {
	DB *sql.DB
}

func (r *courtRepository) Create(payload model.Court) (model.Court, error) {
	var court model.Court

	err := r.DB.QueryRow("INSERT INTO courts (name, price) VALUES ($1, $2) RETURNING *", payload.Name, payload.Price).Scan(&court.Id, &court.Name, &court.Price, &court.CreatedAt, &court.UpdatedAt)

	if err != nil {
		return model.Court{}, err
	}

	return court, nil
}

func (r *courtRepository) FindAll(page int, size int) ([]model.Court, dto.Paginate, error) {
	var courts []model.Court

	// rumus pagination
	offset := (page - 1) * size

	rows, err := r.DB.Query("SELECT id, name, price, created_at, updated_at FROM courts LIMIT $1 OFFSET $2", size, offset)
	if err != nil {
		return []model.Court{}, dto.Paginate{}, err
	}

	totalRows := 0
	for rows.Next() {
		var c model.Court
		if err := rows.Scan(&c.Id, &c.Name, &c.Price, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return []model.Court{}, dto.Paginate{}, err
		}
		courts = append(courts, c)
		totalRows++
	}

	paginate := dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return courts, paginate, nil
}

func (r *courtRepository) FindById(id string) (model.Court, error) {
	var court model.Court

	err := r.DB.QueryRow("SELECT id, name, price, created_at, updated_at FROM courts WHERE id = $1", id).Scan(&court.Id, &court.Name, &court.Price, &court.CreatedAt, &court.UpdatedAt)
	if err != nil {
		return model.Court{}, err
	}

	return court, nil
}

func (r *courtRepository) Update(id string, payload model.Court) (model.Court, error) {
	var court model.Court

	err := r.DB.QueryRow("UPDATE courts SET name = $1, price = $2, updated_at = $3 WHERE id = $4 RETURNING *", payload.Name, payload.Price, time.Now(), id).Scan(&court.Id, &court.Name, &court.Price, &court.CreatedAt, &court.UpdatedAt)
	if err != nil {
		return model.Court{}, err
	}

	return court, nil
}

func (r *courtRepository) Deleted(id string) error {
	_, err := r.DB.Exec("DELETE FROM courts WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func NewCourtRepository(db *sql.DB) CourtRepository {
	return &courtRepository{
		DB: db,
	}
}
