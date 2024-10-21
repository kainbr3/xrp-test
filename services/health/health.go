package health

import (
	"context"
	r "crypto-braza-tokens-api/repositories"
	l "crypto-braza-tokens-api/utils/logger"

	"go.uber.org/zap"
)

type HealthService struct {
	repo *r.Repository
}

func NewHealthService(repo *r.Repository) *HealthService {
	return &HealthService{repo}
}

func (hs *HealthService) Liveness() string {
	return "running"
}

func (hs *HealthService) Readiness(ctx context.Context) (string, error) {
	if err := hs.repo.CheckHealth(ctx); err != nil {
		l.Logger.Error("error checking health", zap.Error(err))
		return "", err
	}

	return "healthy", nil
}
