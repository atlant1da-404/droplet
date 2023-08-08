package auth

type Authenticator interface {
	GenerateToken(options *GenerateTokenClaimsOptions) (string, error)
	ParseToken(accessToken string) (*ParseTokenClaimsOutput, error)
}

type GenerateTokenClaimsOptions struct {
	UserId   string
	UserName string
}

type ParseTokenClaimsOutput struct {
	UserId   string
	Username string
}
