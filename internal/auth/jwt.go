package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func GenerateJWT(userID int64, secret []byte) (string, int, error) {
	const expiresInSeconds = 24 * 60 * 60

	header := jwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	payload := jwtPayload{
		Sub: userID,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(expiresInSeconds).Unix(),
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
	signature := signHS256(message, secret)

	return fmt.Sprintf("%s.%s", message, signature), expiresInSeconds, nil
}

func encode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func signHS256(message string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return encode(h.Sum(nil))
}
