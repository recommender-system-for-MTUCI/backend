package models

import "github.com/google/uuid"

type RegistrationDTO struct {
	ID       uuid.UUID
	Email    string
	Password string
	Active   bool
}
