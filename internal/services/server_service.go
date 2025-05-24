package services

import (
	"context"
	"wg_api/internal/models"
	"wg_api/internal/repository"
)

type ServerService struct {
	repo repository.ServerRepository
}

func NewServerService(repo repository.ServerRepository) *ServerService {
	return &ServerService{repo: repo}
}

func (s *ServerService) Create(ctx context.Context, server *models.Server) error {
	return s.repo.Create(ctx, server)
}

func (s *ServerService) GetByID(ctx context.Context, id uint) (*models.Server, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ServerService) GetAll(ctx context.Context) ([]models.Server, error) {
	return s.repo.GetAll(ctx)
}

func (s *ServerService) Update(ctx context.Context, server *models.Server) error {
	return s.repo.Update(ctx, server)
}

func (s *ServerService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
