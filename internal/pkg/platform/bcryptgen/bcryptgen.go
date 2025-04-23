package bcryptgen

import (
	"websocket_client/internal/pkg/core/adapter/passwordgeneratoradapter"

	"golang.org/x/crypto/bcrypt"
)

type BCryptGen struct {
}

func NewBCryptGen() passwordgeneratoradapter.Adapter {
	return &BCryptGen{}
}

func (b *BCryptGen) Generate(pass string) (string, error) {
	byteHashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(byteHashedPassword), err
}
