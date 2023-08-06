package auth

type Authenticator interface {
	GenerateToken(username string) (string, error)
	ParseToken(accessToken string) (bool, error)
}
