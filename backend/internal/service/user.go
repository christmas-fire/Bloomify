package service

import (
	"context"
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

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetById(ctx context.Context, userId int) (models.User, error) {
	return s.repo.GetById(ctx, userId)
}

func (s *UserService) Delete(ctx context.Context, userId int) error {
	return s.repo.Delete(ctx, userId)
}

func (s *UserService) UpdateUsername(ctx context.Context, userId int, oldUsername, newUsername string) error {
	return s.repo.UpdateUsername(ctx, userId, oldUsername, newUsername)
}

func (s *UserService) UpdatePassword(ctx context.Context, userId int, username, oldPassword, newPassword string) error {
	oldPasswordHash := generatePasswordHash(oldPassword)
	newPasswordHash := generatePasswordHash(newPassword)

	return s.repo.UpdatePassword(ctx, userId, username, oldPasswordHash, newPasswordHash)
}
