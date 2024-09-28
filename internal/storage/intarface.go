package storage

import "context"

type User interface {
	CreateUser(ctx context.Context, email string, password string) error
}
