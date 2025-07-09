package service

import (
	"errors"

	"github.com/JoaoVitorFerreiro/golang-start/internal/application/dto"
	"github.com/JoaoVitorFerreiro/golang-start/internal/domain/entity"
	"github.com/JoaoVitorFerreiro/golang-start/internal/infra/repository"
)

type UserService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}

func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("email already exists")
    }

    newUser, err := entity.NewUser(req.Name, req.Email)
    if err != nil {
        return nil, err
    }

		if err := s.userRepo.Save(newUser); err != nil {
        return nil, err
    }

    return s.toUserResponse(newUser), nil
}

func (s *UserService) GetUserByID(id string) (*dto.UserResponse, error) {
    if id == "" {
        return nil, errors.New("id is required")
    }

    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    return s.toUserResponse(user), nil
}

func (s *UserService) GetAllUsers() ([]*dto.UserResponse, error) {
    users, err := s.userRepo.FindAll()
    if err != nil {
        return nil, err
    }

    var responses []*dto.UserResponse
    for _, u := range users {
        responses = append(responses, s.toUserResponse(u))
    }

    return responses, nil
}

func (s *UserService) UpdateUser(id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
    if id == "" {
        return nil, errors.New("id is required")
    }

    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    if req.Email != "" && req.Email != user.Email {
        existingUser, err := s.userRepo.FindByEmail(req.Email)
        if err != nil {
            return nil, err
        }
        if existingUser != nil {
            return nil, errors.New("email already exists")
        }
        
        if err := user.UpdateEmail(req.Email); err != nil {
            return nil, err
        }
    }

		if req.Name != "" && req.Name != user.Name {
        if err := user.UpdateName(req.Name); err != nil {
            return nil, err
        }
    }

		if err := s.userRepo.Save(user); err != nil {
        return nil, err
    }

    return s.toUserResponse(user), nil
}

func (s *UserService) DeleteUser(id string) error {
    if id == "" {
        return errors.New("id is required")
    }

    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return err
    }
    if user == nil {
        return errors.New("user not found")
    }

    return s.userRepo.Delete(id)
}

// Converter entity para DTO de resposta
func (s *UserService) toUserResponse(u *entity.User) *dto.UserResponse {
    return &dto.UserResponse{
        ID:        u.ID,
        Name:      u.Name,
        Email:     u.Email,
        CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }
}