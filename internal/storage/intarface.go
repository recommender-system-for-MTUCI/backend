package storage

import (
	"context"

	"github.com/google/uuid"
)

type User interface {
	CreateUser(ctx context.Context, id uuid.UUID, email string, password string, active bool) error
	GetUserEmail(ctx context.Context, email string) error
}

type Movie interface {
}

type Storage interface {
	User() User
	Movie() Movie
}

func CreateUser(ctx context.Context, id uuid.UUID, email string, password string, active bool) error {
	type xz struct {
		id       uuid.UUID
		email    string
		password string
		active   bool
	}
	s := []xz{}
	s = append(s, xz{
		id:       id,
		email:    email,
		password: password,
		active:   active,
	})
	return nil
}
