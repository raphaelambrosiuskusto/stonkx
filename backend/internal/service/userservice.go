package service

import (
	"context"
	"stonkx/internal/repository"
)

type UserService struct {
	Store *repository.Store
}

func NewUser(s *repository.Store) *UserService {
	return &UserService{
		Store: s,
	}
}
func (s *UserService) Serve(ctx context.Context) error {
	return nil
}
