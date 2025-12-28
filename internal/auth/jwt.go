package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtPayload struct {
	Sub int64 `json:"sub"`
	Iat int64 `json:"iat"`
	Exp int64 `json:"exp"`
}

type Jwt struct {
	secret []byte
}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

func NewJwt(secret string) *Jwt {
	return &Jwt{secret: []byte(secret)}
}

func (j *Jwt) GenerateJWT(userID int64) (string, int, error) {
	const expiresInSeconds = 24 * 60 * 60

	header := jwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	payload := jwtPayload{
		Sub: userID,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Second * time.Duration(expiresInSeconds)).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", 0, err
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", 0, err
	}

	headerEncoded := encode(headerJSON)
	payloadEncoded := encode(payloadJSON)

	message := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)
	signature := signHS256(message, j.secret)

	return fmt.Sprintf("%s.%s", message, signature), expiresInSeconds, nil
}

func (j *Jwt) ValidateJWT(tokenString string) (int64, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return 0, ErrInvalidToken
	}

	message := parts[0] + "." + parts[1]
	signature, err := decode(parts[2])
	if err != nil {
		return 0, ErrInvalidToken
	}

	mac := hmac.New(sha256.New, j.secret)
	mac.Write([]byte(message))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return 0, ErrInvalidToken
	}

	payloadJSON, err := decode(parts[1])
	if err != nil {
		return 0, ErrInvalidToken
	}

	var payload jwtPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return 0, ErrInvalidToken
	}
	if time.Now().Unix() > payload.Exp {
		return 0, ErrExpiredToken
	}

	return payload.Sub, nil
}

func encode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func signHS256(message string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return encode(h.Sum(nil))
}

func decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}
