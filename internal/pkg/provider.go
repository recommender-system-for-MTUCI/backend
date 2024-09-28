package pkg

import (
	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/backend/internal/models"
)

type Provider interface {
	CreateTokenForUser(userID uuid.UUID, isAccess bool) (string, error)
	GetDataFromToken(token string) (*models.UserDataInToken, error)
}
