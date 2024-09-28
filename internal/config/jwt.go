package config

type JWT struct {
	AccsesTokenLifetime  int    `config:"access-token-lifetime"`
	RefreshTokenLifetime int    `config:"refresh-token-lifetime"`
	PublicKeyPath        string `config:"public-key-path"`
	PrivateKeyPath       string `config:"private-key-path"`
}
