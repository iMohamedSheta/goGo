package controllers

import (
	"encoding/json"
	"imohamedsheta/gocrud/app/models"
	"imohamedsheta/gocrud/app/requests"
	"imohamedsheta/gocrud/pkg/auth"
	"imohamedsheta/gocrud/pkg/encrypt"
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/pkg/response"
	"imohamedsheta/gocrud/pkg/validate"
	"imohamedsheta/gocrud/query"
	"net/http"
	"strings"
	"time"
)

type AuthController struct{}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log().Error(err.Error())
		response.ErrorJson(w, "Invalid JSON format", "invalid_request", http.StatusBadRequest)
		return
	}

	ok, validationErrors := validate.ValidateRequest(&req)
	if !ok {
		response.ValidationErrorJson(w, validationErrors)
		return
	}

	hashedPassword, err := encrypt.HashPassword(req.Password)
	if err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	// Create a new user object (without the ID)
	user := &models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert the user into the database
	if err := query.UsersTable().Insert(user); err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	// Retrieve the user from the database to get the generated ID
	result, err := query.Table("users").Where("username", "=", user.Username).First()
	if err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	if result == nil {
		response.ServerErrorJson(w)
		logger.Log().Error("User not found in the database after successful registration!")
		return
	}

	user.Id = result["id"].(string)

	jwtPayload := map[string]any{
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"id":        user.Id,
	}

	// Generate access token
	accessToken, err := auth.GenerateAccessToken(user.Id, jwtPayload)
	if err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	// Generate refresh token
	refreshToken, err := auth.GenerateRefreshToken(user.Id, jwtPayload)
	if err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	// Respond with the generated tokens
	response.Json(w, "User registered successfully", map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, http.StatusCreated)
}

func (controller *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	req := &requests.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	ok, validationErrors := validate.ValidateRequest(req)
	if !ok {
		response.ValidationErrorJson(w, validationErrors)
		return
	}

	result, err := query.Table("users").Where("username", "=", req.Username).First()
	if err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	if result == nil {
		response.ValidationErrorJson(w, map[string]string{
			"username": "Invalid username",
		})
		return
	}

	// Check if the provided password matches the stored hash
	if !encrypt.CheckPasswordHash(req.Password, result["password"].(string)) {
		response.ValidationErrorJson(w, map[string]string{
			"password": "Invalid password",
		})
		return
	}

	// Prepare the payload for JWT tokens
	jwtPayload := map[string]any{
		"username":  result["username"].(string),
		"id":        result["id"].(string),
		"email":     result["email"].(string),
		"firstName": result["first_name"].(string),
		"lastName":  result["last_name"].(string),
	}

	// Generate access token using the user data
	accessToken, err := auth.GenerateAccessToken(result["id"].(string), jwtPayload)
	if err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	// Generate refresh token using the user data
	refreshToken, err := auth.GenerateRefreshToken(result["id"].(string), jwtPayload)
	if err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	// Respond with the generated tokens
	response.Json(w, "Login successful", map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, http.StatusOK)
}

func (c *AuthController) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from the request header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		response.ErrorJson(w, "Missing or invalid Authorization header", "missing_auth_header", http.StatusUnauthorized)
		return
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate and decode the refresh token
	token, err := auth.ValidateAuthToken(refreshToken, auth.RefreshToken)
	if err != nil {
		response.ErrorJson(w, err.Error(), "invalid_refresh_token", http.StatusUnauthorized)
		return
	}

	userIdRaw, err := token.Get("user_id")

	if err != nil {
		response.ErrorJson(w, err.Error(), "invalid_refresh_token", http.StatusUnauthorized)
		return
	}

	userId, ok := userIdRaw.(string)
	if !ok {
		response.ErrorJson(w, "Invalid user ID", "invalid_user_id", http.StatusUnauthorized)
		return
	}

	// Generate a new access token using the user ID and claims
	accessToken, err := auth.GenerateAccessToken(userId, token.Payload)
	if err != nil {
		response.ErrorJson(w, "Failed to generate new access token", "access_token_error", http.StatusInternalServerError)
		return
	}

	// Respond with the new access token
	response.Json(w, "New access token generated", map[string]any{
		"access_token": accessToken,
	}, http.StatusOK)
}
