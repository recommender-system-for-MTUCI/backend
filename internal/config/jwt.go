package config

type JWT struct {
	AccsesTokenLifetime  int `config:"access-token-lifetime"`
	RefreshTokenLifetime int `config:"refresh-token-lifetime"`
}
