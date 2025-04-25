package controllers

import (
	"encoding/json"
	"imohamedsheta/gocrud/app/models"
	"imohamedsheta/gocrud/app/requests"
	"imohamedsheta/gocrud/pkg/config"
	"imohamedsheta/gocrud/pkg/encrypt"
	"imohamedsheta/gocrud/pkg/jwt"
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/pkg/response"
	"imohamedsheta/gocrud/pkg/validate"
	"imohamedsheta/gocrud/query"
	"net/http"
	"time"
)

type AuthController struct{}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	// Create a new user
	// Generate new token
	// Return token

	var req requests.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log().Error(err.Error())
		response.ErrorJson(w, "Invalid JSON format", http.StatusBadRequest)
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

	user := &models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := query.UsersTable().Insert(user); err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	jwtPayload := map[string]any{
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	}

	jwtSecret, _ := config.App.Get("app.secret").(string)
	jwtExpiry := time.Duration(config.App.Get("app.jwt_expiry").(int)) * time.Second

	token, err := jwt.GenerateJWTToken(jwtPayload, jwtSecret, jwtExpiry)

	if err != nil {
		logger.Log().Error(err.Error())
		response.ServerErrorJson(w)
		return
	}

	response.Json(w, "User created successfully", map[string]any{
		"token": token,
	}, http.StatusCreated)
}
