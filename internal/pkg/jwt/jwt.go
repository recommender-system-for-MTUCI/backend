package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/backend/internal/config"
	"github.com/recommender-system-for-MTUCI/backend/internal/pkg/jwt"
)

type Provider struct {
	publicKey            []byte
	privateKey           []byte
	accessTokenLifetime  int
	refreshTokenLifetime int
}

func NewProvider(cfg *config.JWT) *Provider {
	return &Provider{
		publicKey:  make([]byte, 0),
		privateKey: make([]byte, 0),
	}
}

func (provider *Provider) readKeyFunc(token *jwt.Token) (interface{}, error) {

	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return provider.publicKey, nil
}

func (provider *Provider) CreateTokenForUser(userID uuid.UUID, isAccess bool) (string, error) {
	now := time.Now()

	var add time.Duration

	if isAccess {
		add = time.Duration(provider.refreshTokenLifetime) * time.Minute
	} else {
		add = time.Duration(provider.refreshTokenLifetime) * time.Minute
	}
	claims := jwt.StandardClaims{
		Issuer:    userID.String(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(add).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.SignedString()
}

func (provider *Provider) GetDataFromToken(token string) (uuid.UUID, error) {
	parsedToken, err := jwt.Parse(token, provider.readKeyFunc)
	if err != nil {
		return uuid.Nil, err
	}
	if !parsedToken.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	if claims, ok := parsedToken.Claims(jwt.StandardClaims); ok {
		return uuid.Parse(claims.Issuer)
	}
	return uuid.Nil, fmt.Errorf("invalid token")
}
