package auth

import (
	"errors"
	"fmt"
	"imohamedsheta/gocrud/pkg/config"
	"imohamedsheta/gocrud/pkg/jwt"
	"imohamedsheta/gocrud/pkg/logger"
	"time"
)

type AuthTokenType string

const (
	RefreshToken AuthTokenType = "refresh_token"
	AccessToken  AuthTokenType = "access_token"
)

// generate new refresh token, refresh token is used to generate new access token long lived tokens (7 days)
func GenerateRefreshToken(userID int64, claims map[string]any) (string, error) {
	refreshTokenExpirationTime := config.App.Get("app.auth.refresh_token_expiry").(time.Duration)
	secret := config.App.Get("app.secret").(string)
	payload := map[string]any{
		"user_id": userID,
		"type":    RefreshToken,
	}

	for k, v := range claims {
		if k == "user_id" || k == "type" || k == "iat" || k == "exp" { // ignore these fields if added to claims they will be added to payload
			continue
		}
		payload[k] = v
	}

	return jwt.GenerateJWTToken(payload, secret, refreshTokenExpirationTime)
}

// generate new access token, access token is used to authenticate user they are short lived tokens and need to be refreshed every 15 minutes
func GenerateAccessToken(userID int64, claims map[string]any) (string, error) {
	accessTokenExpirationTime := config.App.Get("app.auth.access_token_expiry").(time.Duration)
	secret := config.App.Get("app.secret").(string)
	payload := map[string]any{
		"user_id": userID,
		"type":    AccessToken,
	}

	for k, v := range claims {
		if k == "user_id" || k == "type" || k == "iat" || k == "exp" {
			continue
		}
		payload[k] = v
	}

	return jwt.GenerateJWTToken(payload, secret, accessTokenExpirationTime)
}

// ValidateRefreshToken validates the refresh token and returns the decoded JWT token
func ValidateAuthToken(jwtToken string, tokenType AuthTokenType) (*jwt.JWT, error) {
	secret := config.App.Get("app.secret").(string)

	valid, err := jwt.Verify(jwtToken, secret)
	if err != nil || !valid {

		return nil, fmt.Errorf("invalid %s token", string(tokenType))
	}

	jwtTokenDecoded, err := jwt.DecodeJWT(jwtToken)
	if err != nil {
		logger.Log().Error(err.Error())
		return nil, errors.New("invalid  token")
	}

	extractedTokenTypeRaw, err := jwtTokenDecoded.Get("type")
	if err != nil {
		return nil, errors.New("invalid token type")
	}

	extractedTokenType, ok := extractedTokenTypeRaw.(string)
	if !ok || extractedTokenType != string(tokenType) {
		return nil, errors.New("invalid token type")
	}

	exp, err := jwtTokenDecoded.Get("exp")
	if err != nil {
		return nil, errors.New("missing expiration field")
	}

	expirationTime := time.Unix(int64(exp.(float64)), 0)
	if time.Now().After(expirationTime) {
		return nil, fmt.Errorf("expired %s token", string(tokenType))
	}

	return jwtTokenDecoded, nil
}
