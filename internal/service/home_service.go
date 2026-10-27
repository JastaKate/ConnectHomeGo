package service

import (
	"context"
	"errors"
	"strings"

	"connectHome/internal/model"
	"connectHome/internal/repository"
	"connectHome/internal/request"
)

var ErrInvalidHome = errors.New("Название и адрес обязательны")

type HomeService struct {
	repository *repository.HomeRepository
}

func NewHomeService(homeRepository *repository.HomeRepository) *HomeService {
	return &HomeService{repository: homeRepository}
}

func (s *HomeService) Create(ctx context.Context, input request.HomeRequest) (model.Home, error) {
	name := strings.TrimSpace(input.Name)
	address := strings.TrimSpace(input.Address)
	if name == "" || address == "" {
		return model.Home{}, ErrInvalidHome
	}
	return s.repository.Create(ctx, name, address)
}

func (s *HomeService) List(ctx context.Context) ([]model.Home, error) {
	return s.repository.List(ctx)
}
