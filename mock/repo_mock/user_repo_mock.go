package repomock

import (
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) CreateCustomer(payload model.User) (model.User, error) {
	args := u.Called(payload)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserRepositoryMock) CreateEmployee(payload model.User) (model.User, error) {
	args := u.Called(payload)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserRepositoryMock) CreateAdmin(payload model.User) (model.User, error) {
	args := u.Called(payload)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserRepositoryMock) FindUserByUsername(username string) (model.User, error) {
	args := u.Called(username)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserRepositoryMock) FindUserById(id string) (model.User, error) {
	args := u.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserRepositoryMock) FindUserByRole(role string, page int, size int) ([]model.User, dto.Paginate, error) {
	args := u.Called(role, page, size)
	return args.Get(0).([]model.User), args.Get(1).(dto.Paginate), args.Error(2)
}
func (u *UserRepositoryMock) UpdateUser(id string, payload model.User) (model.User, error) {
	args := u.Called(id, payload)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserRepositoryMock) DeleteUser(id string) error {
	args := u.Called(id)
	return args.Error(0)
}
