package jwt

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/backend/internal/config"
	"go.uber.org/zap"
)

type Provider struct {
	publicKey            *rsa.PublicKey
	privateKey           *rsa.PrivateKey
	accessTokenLifetime  int
	refreshTokenLifetime int
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

	claims := jwt.RegisteredClaims{
		Issuer:    userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(add)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(provider.privateKey)
}

func (provider *Provider) GetDataFromToken(token string) (uuid.UUID, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, provider.readKeyFunc)
	if err != nil {
		return uuid.Nil, err
	}
	if !parsedToken.Valid {
		return uuid.Nil, fmt.Errorf("invalid token: not valid")
	}
	if claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims); ok {
		return uuid.Parse(claims.Issuer)
	}
	return uuid.Nil, fmt.Errorf("invalid token: cannot parse claims")
}
