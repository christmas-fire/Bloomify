package service

import (
	"fmt"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetById(userId int) (models.User, error) {
	return s.repo.GetById(userId)
}

func (s *UserService) Delete(userId int) error {
	return s.repo.Delete(userId)
}

func (s *UserService) UpdateUsername(userId int, input models.UpdateUsernameInput) error {
	if err := validateUsername(input.NewUsername); err != nil {
		return err
	}
	return s.repo.UpdateUsername(userId, input)
}

func (s *UserService) UpdatePassword(userId int, input models.UpdatePasswordInput) error {
	if err := validatePassword(input.NewPassword); err != nil {
		return err
	}

	input.OldPassword = generatePasswordHash(input.OldPassword)
	input.NewPassword = generatePasswordHash(input.NewPassword)

	return s.repo.UpdatePassword(userId, input)
}

func validateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("username must have at least 3 characters")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must have at least 8 characters")
	}

	return nil
}
