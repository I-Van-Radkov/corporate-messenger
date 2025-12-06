package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPasswordBase64(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hashPasswordStr := getHashByte(password, salt)

	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashPasswordBase64 := base64.StdEncoding.EncodeToString(hashPasswordStr)

	return fmt.Sprintf("%s.%s", saltBase64, hashPasswordBase64), nil
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to create salt: %w", err)
	}

	return salt, nil
}

func getHashByte(password string, salt []byte) []byte {
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return hashedPassword
}

func VerifyPassword(newPasswordStr, realHashedPasswordBase64 string) (bool, error) {
	parts := strings.Split(realHashedPasswordBase64, ".")

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}
	hashedRealPassword, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("failed to decode password: %w", err)
	}

	hashedNewPassword := getHashByte(newPasswordStr, salt)

	if len(hashedNewPassword) != len(hashedRealPassword) {
		return false, nil
	}
	if subtle.ConstantTimeCompare(hashedNewPassword, hashedRealPassword) == 0 {
		return false, nil
	}

	return true, nil
}
