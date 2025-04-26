package jwt

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

type JWT struct {
	Header    JWTHeader      `json:"header"`
	Payload   map[string]any `json:"payload"` // Claims
	Signature string         `json:"signature"`
}

// JWTHeader represents the header part of the JWT
type JWTHeader struct {
	Alg  string `json:"alg"`
	Type string `json:"typ"`
}

// GenerateJWTToken generates a JWT token with user data
func GenerateJWTToken(payload map[string]any, secret string, expirationTime time.Duration) (string, error) {
	// JWT Header
	header := JWTHeader{
		Alg:  "HS256", // HMAC-SHA256 is the default algorithm and the only one supported
		Type: "JWT",
	}

	// Add standard claims like 'iat' (issued at) and 'exp' (expiration)
	payload["iat"] = time.Now().Unix()
	payload["exp"] = time.Now().Add(expirationTime).Unix()

	// Encode Header and Payload
	encodedHeader, err := encodeToBase64(header)
	if err != nil {
		return "", err
	}

	encodedPayload, err := encodeToBase64(payload)
	if err != nil {
		return "", err
	}

	// Create the Signature
	message := encodedHeader + "." + encodedPayload
	signature, err := signHMAC256(message, secret)
	if err != nil {
		return "", err
	}

	// Build the final JWT token
	token := encodedHeader + "." + encodedPayload + "." + signature
	return token, nil
}

// Get allows accessing any key from JWT payload dynamically
func (jwt *JWT) Get(key string) (any, error) {
	value, exists := jwt.Payload[key]
	if !exists {
		return nil, fmt.Errorf("key '%s' not found in JWT payload", key)
	}
	return value, nil
}

func (jwt *JWT) Set(key string, value any) {
	jwt.Payload[key] = value
}

func Verify(token string, secret string) (bool, error) {
	// Verify the token's signature using HMAC-SHA256 (the only supported algorithm)
	isVerified, err := verifyHMAC256(token, secret)
	if err != nil {
		return false, err
	}

	return isVerified, nil
}

// DecodeJWT decodes a JWT token into a JWT struct
func DecodeJWT(token string) (*JWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format: expected 3 parts")
	}

	// Decode and Unmarshal the payload
	decodedPayload, err := decodeFromBase64(parts[1])
	if err != nil {
		return nil, err
	}

	var payloadData map[string]any
	err = json.Unmarshal(decodedPayload, &payloadData)
	if err != nil {
		return nil, errors.New("failed to unmarshal payload: " + err.Error())
	}

	// Decode and Unmarshal the header (optional)
	decodedHeader, err := decodeFromBase64(parts[0])
	if err != nil {
		return nil, err
	}

	var headerData JWTHeader
	err = json.Unmarshal(decodedHeader, &headerData)
	if err != nil {
		return nil, errors.New("failed to unmarshal header: " + err.Error())
	}

	// Return the full JWT data (header and payload)
	return &JWT{
		Header:    headerData,
		Payload:   payloadData,
		Signature: parts[2],
	}, nil
}

/*
---------------------------------------------------------
| Helpers Functions
---------------------------------------------------------
*/

// SignHMAC256 generates an HMAC-SHA256 signature
func signHMAC256(message, secret string) (string, error) {
	hmac := hmac.New(sha256.New, []byte(secret))
	hmac.Write([]byte(message))
	signature := hmac.Sum(nil)
	encodedSignature, err := encodeToBase64(signature)

	if err != nil {
		return "", err
	}

	return encodedSignature, nil
}

// VerifyHMAC256 verifies the HMAC-SHA256 signature of the JWT
func verifyHMAC256(token, secret string) (bool, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, errors.New("invalid token format: expected 3 parts")
	}

	header, payload, signature := parts[0], parts[1], parts[2]
	message := header + "." + payload

	expectedSignature, err := signHMAC256(message, secret)

	if err != nil {
		return false, err
	}

	return hmac.Equal([]byte(signature), []byte(expectedSignature)), nil
}

// EncodeToBase64 encodes data to a Base64 string
func encodeToBase64(data any) (string, error) {
	encodedData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encodedData), nil
}

// DecodeFromBase64 decodes a Base64 string to bytes
func decodeFromBase64(encodedData string) ([]byte, error) {
	decodedData, err := base64.RawURLEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, err
	}
	return decodedData, nil
}
