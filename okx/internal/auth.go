package internal

type Auth struct {
	ApiKey     string
	Passphrase string
	SecretKey  string
}

func NewAuth(apiKey, passphrase, secretKey string) *Auth {
	return &Auth{ApiKey: apiKey, Passphrase: passphrase, SecretKey: secretKey}
}
