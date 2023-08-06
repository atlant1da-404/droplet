package hash

import "golang.org/x/crypto/bcrypt"

type bcryptor struct {
}

func NewHash() Hash {
	return &bcryptor{}
}

func (b *bcryptor) GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (b *bcryptor) CompareHash(hashedPassword, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return err
	}
	return nil
}
