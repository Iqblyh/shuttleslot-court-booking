package repomock

import (
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"

	"github.com/stretchr/testify/mock"
)

type CourtRepositoryMock struct {
	mock.Mock
}

func (c *CourtRepositoryMock) Create(payload model.Court) (model.Court, error) {
	args := c.Called(payload)
	return args.Get(0).(model.Court), args.Error(1)
}
func (c *CourtRepositoryMock) FindAll(page int, size int) ([]model.Court, dto.Paginate, error) {
	args := c.Called(page, size)
	return args.Get(0).([]model.Court), args.Get(1).(dto.Paginate), args.Error(2)
}
func (c *CourtRepositoryMock) FindById(id string) (model.Court, error) {
	args := c.Called(id)
	return args.Get(0).(model.Court), args.Error(1)
}
func (c *CourtRepositoryMock) Update(id string, payload model.Court) (model.Court, error) {
	args := c.Called(id, payload)
	return args.Get(0).(model.Court), args.Error(1)
}
func (c *CourtRepositoryMock) Deleted(id string) error {
	args := c.Called(id)
	return args.Error(0)
}
