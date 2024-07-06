package model

import "time"

type User struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Point       int       `json:"point"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (u User) IsValidRole() bool {
	return u.Role == "admin" || u.Role == "employee" || u.Role == "customer"
}
