package service

import (
	"context"
	"errors"

	"github.com/haakaashs/todos-backend/internal/model"
)

var ErrNotFound = errors.New("todo not found")

type Repository interface {
	Create(context.Context, string) (model.Todo, error)
	Get(context.Context, string) (model.Todo, error)
	Update(context.Context, *model.Todo) (model.Todo, error)
	Delete(context.Context, string) error
	List(context.Context) ([]model.Todo, error)
}

type Service struct {
	repo Repository
}

func NewTodosService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, title string) (model.Todo, error) {
	// business logic here
	return s.repo.Create(ctx, title)
}

func (s *Service) Get(ctx context.Context, id string) (model.Todo, error) {
	// business logic here
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]model.Todo, error) {
	// business logic here
	return s.repo.List(ctx)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	// business logic here
	return s.repo.Delete(ctx, id)
}

func (s *Service) Update(ctx context.Context, t *model.Todo) (model.Todo, error) {
	// business logic here
	return s.repo.Update(ctx, t)
}
