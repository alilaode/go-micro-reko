package user

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	cryptFormat = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
)

func (ur *UserRepo) GenerateUserHash(ctx context.Context, password string) (hash string, err error) {

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	argonHash := argon2.IDKey([]byte(password), salt, ur.time, ur.memory, ur.threads, ur.keylen)

	b64Hash := ur.encrypt(ctx, argonHash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	encodeHash := fmt.Sprintf(cryptFormat, argon2.Version, ur.memory, ur.time, ur.threads, b64Salt, b64Hash)

	return encodeHash, nil
}

func (ur *UserRepo) encrypt(ctx context.Context, text []byte) string {

	nonce := make([]byte, ur.gcm.NonceSize())

	ciphertext := ur.gcm.Seal(nonce, nonce, text, nil)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func (ur *UserRepo) decrypt(ctx context.Context, chiphertext string) ([]byte, error) {

	decode, err := base64.StdEncoding.DecodeString(chiphertext)
	if err != nil {
		return nil, err
	}

	if len(decode) < ur.gcm.NonceSize() {
		return nil, errors.New("invalid nonce size")
	}

	return ur.gcm.Open(nil,
		decode[:ur.gcm.NonceSize()],
		decode[ur.gcm.NonceSize():],
		nil)

}

func (ur *UserRepo) comparePassword(ctx context.Context, password, hash string) (bool, error) {

	parts := strings.Split(hash, "$")
	if len(parts) < 6 {
		return false, nil
	}

	var memory, time uint32
	var parallelisem uint8

	switch parts[1] {
	case "argon2id":

		_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelisem)
		if err != nil {
			return false, err
		}

		salt, err := base64.RawStdEncoding.DecodeString(parts[4])
		if err != nil {
			return false, err
		}

		hash := parts[5]

		decryptedHash, err := ur.decrypt(ctx, hash)
		if err != nil {
			return false, err
		}

		var keyLen = uint32(len(decryptedHash))

		comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, parallelisem, keyLen)

		return subtle.ConstantTimeCompare(comparisonHash, decryptedHash) == 1, nil

	}

	return false, nil

}
