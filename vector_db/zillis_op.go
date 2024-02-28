package vector_db

import (
	"fmt"
	"github.com/o1egl/paseto"
)

func NewPasetoMaker(symmetricKey []byte) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: %d, it must be exactly %d characters", len(symmetricKey), chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}

	return maker, nil
}
