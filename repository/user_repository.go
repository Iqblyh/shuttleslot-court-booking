package repository

import (
	"database/sql"
	"math"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"time"
)

type UserRepository interface {
	CreateCustomer(payload model.User) (model.User, error)
	CreateEmployee(payload model.User) (model.User, error)
	CreateAdmin(payload model.User) (model.User, error)
	FindUserByUsername(username string) (model.User, error)
	FindUserById(id string) (model.User, error)
	FindUserByRole(role string, page int, size int) ([]model.User, dto.Paginate, error)
	UpdateUser(id string, payload model.User) (model.User, error)
	DeleteUser(id string) error
}

type userRepository struct {
	DB *sql.DB
}

func (r *userRepository) CreateCustomer(payload model.User) (model.User, error) {
	var customer model.User

	err := r.DB.QueryRow("INSERT INTO users (name, phone_number, email, username, password, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, phone_number, email, username, password, points, role, created_at, updated_at", payload.Name, payload.PhoneNumber, payload.Email, payload.Username, payload.Password, payload.Role).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Email, &customer.Username, &customer.Password, &customer.Point, &customer.Role, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return model.User{}, err
	}

	return customer, nil
}

func (r *userRepository) CreateEmployee(payload model.User) (model.User, error) {
	var employee model.User

	err := r.DB.QueryRow("INSERT INTO users (name, phone_number, email, username, password, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, phone_number, email, username, password, points, role, created_at, updated_at", payload.Name, payload.PhoneNumber, payload.Email, payload.Username, payload.Password, payload.Role).Scan(&employee.Id, &employee.Name, &employee.PhoneNumber, &employee.Email, &employee.Username, &employee.Password, &employee.Point, &employee.Role, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		return model.User{}, err
	}

	return employee, nil
}

func (r *userRepository) CreateAdmin(payload model.User) (model.User, error) {
	var admin model.User

	err := r.DB.QueryRow("INSERT INTO users (name, phone_number, email, username, password, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, phone_number, email, username, password, points, role, created_at, updated_at", payload.Name, payload.PhoneNumber, payload.Email, payload.Username, payload.Password, payload.Role).Scan(&admin.Id, &admin.Name, &admin.PhoneNumber, &admin.Email, &admin.Username, &admin.Password, &admin.Point, &admin.Role, &admin.CreatedAt, &admin.UpdatedAt)
	if err != nil {
		return model.User{}, err
	}

	return admin, nil
}

func (r *userRepository) FindUserByUsername(username string) (model.User, error) {
	var user model.User

	err := r.DB.QueryRow("SELECT id, name, phone_number, email, username, password, points, role, created_at, updated_at FROM users WHERE username = $1", username).Scan(&user.Id, &user.Name, &user.PhoneNumber, &user.Email, &user.Username, &user.Password, &user.Point, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *userRepository) FindUserById(id string) (model.User, error) {
	var user model.User

	err := r.DB.QueryRow("SELECT id, name, phone_number, email, username, password, points, role, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.PhoneNumber, &user.Email, &user.Username, &user.Password, &user.Point, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *userRepository) FindUserByRole(role string, page int, size int) ([]model.User, dto.Paginate, error) {
	var users []model.User

	// rumus pagination
	offset := (page - 1) * size

	rows, err := r.DB.Query("SELECT id, name, phone_number, email, username, password, points, role, created_at, updated_at FROM users WHERE role = $1 LIMIT $2 OFFSET $3", role, size, offset)
	if err != nil {
		return []model.User{}, dto.Paginate{}, err
	}

	totalRows := 0
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Id, &u.Name, &u.PhoneNumber, &u.Email, &u.Username, &u.Password, &u.Point, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return []model.User{}, dto.Paginate{}, err
		}
		users = append(users, u)
		totalRows++
	}

	paginate := dto.Paginate{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return users, paginate, nil
}

func (r *userRepository) UpdateUser(id string, payload model.User) (model.User, error) {
	var user model.User

	err := r.DB.QueryRow("UPDATE users SET name = $1, phone_number = $2, email = $3, username = $4, password = $5, updated_at = $6 WHERE id = $7 RETURNING id, name, phone_number, email, username, password, points, role, created_at, updated_at", payload.Name, payload.PhoneNumber, payload.Email, payload.Username, payload.Password, time.Now(), id).Scan(&user.Id, &user.Name, &user.PhoneNumber, &user.Email, &user.Username, &user.Password, &user.Point, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteUser(id string) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}
