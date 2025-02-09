package service

import (
	"errors"
	"team2/shuttleslot/model"
	"team2/shuttleslot/model/dto"
	"team2/shuttleslot/repository"
	"team2/shuttleslot/util"
)

type UserService interface {
	CreateAdmin(payload model.User) (model.User, error)
	CreateCustomer(payload model.User) (model.User, error)
	CreateEmployee(payload model.User) (model.User, error)
	FindUserByRole(role string, page int, size int) ([]model.User, dto.Paginate, error)
	FindUserByUsername(username string) (model.User, error)
	FindUserById(id string) (model.User, error)
	UpdatedUser(id string, payload model.User) (model.User, error)
	DeletedUser(id string) error
	Login(payload dto.LoginRequest) (dto.LoginResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
	auth           AuthService
	util           util.UtilInterface
}

// Login implements UserService.
func (s *userService) Login(payload dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.userRepository.FindUserByUsername(payload.Username)
	if err != nil {
		return dto.LoginResponse{}, errors.New("username or password invalid! ")
	}

	err = s.util.ComparePasswordHash(user.Password, payload.Password)
	if err != nil {
		return dto.LoginResponse{}, errors.New("username or password invalid! ")
	}

	user.Password = ""
	token, err := s.auth.GenerateToken(user)
	if err != nil {
		return dto.LoginResponse{}, errors.New("failed create token! ")
	}
	return token, nil
}

// CreateAdmin implements UserService.
func (s *userService) CreateAdmin(payload model.User) (model.User, error) {
	passwordHash, err := s.util.EncryptPassword(payload.Password)
	if err != nil {
		return model.User{}, err
	}
	payload.Password = passwordHash
	payload.Role = "admin"

	return s.userRepository.CreateAdmin(payload)
}

// CreateCustomer implements UserService.
func (s *userService) CreateCustomer(payload model.User) (model.User, error) {
	passwordHash, err := s.util.EncryptPassword(payload.Password)
	if err != nil {
		return model.User{}, err
	}
	payload.Password = passwordHash
	payload.Role = "customer"

	return s.userRepository.CreateCustomer(payload)
}

// CreateEmployee implements UserService.
func (s *userService) CreateEmployee(payload model.User) (model.User, error) {
	passwordHash, err := s.util.EncryptPassword(payload.Password)
	if err != nil {
		return model.User{}, err
	}
	payload.Password = passwordHash
	payload.Role = "employee"

	return s.userRepository.CreateEmployee(payload)
}

// FindUserByRole implements UserService.
func (s *userService) FindUserByRole(role string, page int, size int) ([]model.User, dto.Paginate, error) {
	return s.userRepository.FindUserByRole(role, page, size)
}

// FindUserByUsername implements UserService.
func (s *userService) FindUserByUsername(username string) (model.User, error) {
	user, err := s.userRepository.FindUserByUsername(username)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// FindUserById implements UserService.
func (s *userService) FindUserById(id string) (model.User, error) {
	user, err := s.userRepository.FindUserById(id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// UpdateUser implements UserService.
func (s *userService) UpdatedUser(id string, payload model.User) (model.User, error) {

	user, err := s.userRepository.FindUserById(id)
	if err != nil {
		return model.User{}, errors.New("user not found")
	}

	passwordHash := ""

	if payload.Name == "" {
		payload.Name = user.Name
	}

	if payload.Email == "" {
		payload.Email = user.Email
	}

	if payload.Username == "" {
		payload.Username = user.Username
	}

	if payload.PhoneNumber == "" {
		payload.PhoneNumber = user.PhoneNumber
	}

	if payload.Password == "" {
		passwordHash = user.Password
	} else {
		passwordHash, err = s.util.EncryptPassword(payload.Password)
		if err != nil {
			return model.User{}, errors.New("error in encrypting password")
		}
	}

	payload.Password = passwordHash

	return s.userRepository.UpdateUser(id, payload)
}

// DeleteUser implements UserService.
func (s *userService) DeletedUser(id string) error {
	_, err := s.userRepository.FindUserById(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepository.DeleteUser(id)
}

func NewUserService(userRepository repository.UserRepository, authService AuthService, util util.UtilInterface) UserService {
	return &userService{
		userRepository: userRepository,
		auth:           authService,
		util:           util,
	}
}
