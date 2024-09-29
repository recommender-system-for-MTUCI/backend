package models

import "github.com/google/uuid"

type UserDataInToken struct {
	ID       uuid.UUID
	IsAccess bool
}

type RegistrationResponse struct {
	ID           uuid.UUID
	AccessToken  string
	RefreshToken string
}
