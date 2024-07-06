package model

import "time"

type Court struct {
	Id        string    `json:"id"`
	Name      string    `josn:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
