package service

import (
	"errors"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"team2/shuttleslot/repository"
)

type CourtService interface {
	CreateCourt(payload model.Court) (model.Court, error)
	FindAllCourts(page int, size int) ([]model.Court, dto.Paginate, error)
	FindCourtById(id string) (model.Court, error)
	UpdateCourt(id string, payload model.Court) (model.Court, error)
	DeleteCourt(id string) error
}

type courtService struct {
	courtRepository repository.CourtRepository
}

func (s *courtService) CreateCourt(payload model.Court) (model.Court, error) {
	court, err := s.courtRepository.Create(payload)
	if err != nil {
		return model.Court{}, err
	}

	return court, nil
}

func (s *courtService) FindAllCourts(page int, size int) ([]model.Court, dto.Paginate, error) {
	return s.courtRepository.FindAll(page, size)
}

func (s *courtService) FindCourtById(id string) (model.Court, error) {
	court, err := s.courtRepository.FindById(id)
	if err != nil {
		return model.Court{}, err
	}

	return court, nil
}

func (s *courtService) UpdateCourt(id string, payload model.Court) (model.Court, error) {

	court, err := s.courtRepository.FindById(id)
	if err != nil {
		return model.Court{}, err
	}

	if payload.Name == "" {
		payload.Name = court.Name
	}

	if payload.Price < 1 {
		payload.Price = court.Price
	}

	courtUpdate, err := s.courtRepository.Update(id, payload)
	if err != nil {
		return model.Court{}, err
	}

	return courtUpdate, nil
}

func (s *courtService) DeleteCourt(id string) error {
	_, err := s.courtRepository.FindById(id)
	if err != nil {
		return errors.New("court not found")
	}

	err = s.courtRepository.Deleted(id)
	if err != nil {
		return err
	}

	return nil
}

func NewCourtService(courtRepository repository.CourtRepository) CourtService {
	return &courtService{courtRepository: courtRepository}
}
