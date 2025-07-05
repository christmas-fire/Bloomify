package service

import (
	"log/slog"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/repository"
)

type UserService struct {
	repo   repository.User
	logger *slog.Logger
}

func NewUserService(repo repository.User, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
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

func (s *UserService) UpdateUsername(userId int, oldUsername, newUsername string) error {
	return s.repo.UpdateUsername(userId, oldUsername, newUsername)
}

func (s *UserService) UpdatePassword(userId int, username, oldPassword, newPassword string) error {
	oldPasswordHash := generatePasswordHash(oldPassword)
	newPasswordHash := generatePasswordHash(newPassword)

	return s.repo.UpdatePassword(userId, username, oldPasswordHash, newPasswordHash)
}
