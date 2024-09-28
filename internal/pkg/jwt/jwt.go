package jwt

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/backend/internal/config"
	"github.com/recommender-system-for-MTUCI/backend/internal/models"
	"go.uber.org/zap"
)

type Provider struct {
	publicKey            *rsa.PublicKey
	privateKey           *rsa.PrivateKey
	accessTokenLifetime  int
	refreshTokenLifetime int
}

type CustomClaims struct {
	jwt.RegisteredClaims
	IsAccess bool
}

func NewProvider(cfg *config.JWT, log *zap.Logger) (*Provider, error) {
	prvKey, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading private key")
	}
	privetKey, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return nil, fmt.Errorf("error while parse private key")
	}
	pubKey, err := os.ReadFile(cfg.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading public key")
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return nil, fmt.Errorf("error while parse public key")
	}
	log.Info("public key", zap.ByteString("public key", pubKey))
	provider := &Provider{
		publicKey:            publicKey,
		privateKey:           privetKey,
		accessTokenLifetime:  cfg.AccsesTokenLifetime,
		refreshTokenLifetime: cfg.RefreshTokenLifetime,
	}
	return provider, nil
}

func (provider *Provider) readKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return provider.publicKey, nil
}

func (provider *Provider) writeKeyfunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return provider.privateKey, nil
}

func (provider *Provider) CreateTokenForUser(userID uuid.UUID, isAccess bool) (string, error) {
	now := time.Now()

	var add time.Duration
	if isAccess {
		add = time.Duration(provider.accessTokenLifetime) * time.Minute
	} else {
		add = time.Duration(provider.refreshTokenLifetime) * time.Minute
	}

	claims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(add)),
		},
		IsAccess: isAccess,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(provider.privateKey)
}

func (provider *Provider) GetDataFromToken(token string) (*models.UserDataInToken, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, provider.readKeyFunc)
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token: not valid")
	}
	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token: cannot parse claim")
	}
	var parsedID uuid.UUID

	parsedID, err = uuid.Parse(claims.Issuer)
	if err != nil {
		return nil, err
	}
	return &models.UserDataInToken{
		ID:       parsedID,
		IsAccess: claims.IsAccess,
	}, nil
}
