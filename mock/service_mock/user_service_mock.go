package servicemock

import (
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"

	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (u *UserServiceMock) CreateAdmin(payload model.User) (model.User, error) {
	args := u.Called(payload)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserServiceMock) CreateCustomer(payload model.User) (model.User, error) {
	args := u.Called(payload)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserServiceMock) CreateEmployee(payload model.User) (model.User, error) {
	args := u.Called(payload)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserServiceMock) FindUserByRole(role string, page int, size int) ([]model.User, dto.Paginate, error) {
	args := u.Called(role, page, size)
	return args.Get(0).([]model.User), args.Get(1).(dto.Paginate), args.Error(2)
}
func (u *UserServiceMock) FindUserByUsername(username string) (model.User, error) {
	args := u.Called(username)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserServiceMock) FindUserById(id string) (model.User, error) {
	args := u.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserServiceMock) UpdatedUser(id string, payload model.User) (model.User, error) {
	args := u.Called(id, payload)
	return args.Get(0).(model.User), args.Error(1)
}
func (u *UserServiceMock) DeletedUser(id string) error {
	args := u.Called(id)
	return args.Error(0)
}
func (u *UserServiceMock) Login(payload dto.LoginRequest) (dto.LoginResponse, error) {
	args := u.Called(payload)
	return args.Get(0).(dto.LoginResponse), args.Error(1)
}
