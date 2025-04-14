package services

import (
	"context"
	"wg_api/internal/models"
	"wg_api/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
