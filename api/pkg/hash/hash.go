package hash

type Hash interface {
	GenerateHash(password string) (string, error)
	CompareHash(hashedPassword, password []byte) error
}
