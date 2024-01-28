package db

import (
	"context"

	"github.com/Pyotr23/the-box/bot/internal/pkg/model"
)

type repository interface {
	InsertJob(ctx context.Context, data model.JobSettingsChatID) error
	GetJobs(ctx context.Context) ([]model.JobSettingsChatID, error)
}

type Service struct {
	repo repository
}

func NewDbService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetJobs(ctx context.Context) ([]model.JobSettingsChatID, error) {
	return s.repo.GetJobs(ctx)
}

func (s *Service) AddJob(ctx context.Context, data model.JobSettingsChatID) error {
	return s.repo.InsertJob(ctx, data)
}
