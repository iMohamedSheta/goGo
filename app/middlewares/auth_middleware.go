package middlewares

import (
	"context"
	"imohamedsheta/gocrud/app/enums"
	"imohamedsheta/gocrud/pkg/config"
	"imohamedsheta/gocrud/pkg/jwt"
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/pkg/response"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				response.ErrorJson(w, "Missing or invalid Authorization header", "missing_auth_header", http.StatusUnauthorized)
				return
			}

			secret := config.App.Get("app.secret").(string)

			token := strings.TrimPrefix(authHeader, "Bearer ")
			valid, err := jwt.Verify(token, secret)
			if err != nil || !valid {
				response.ErrorJson(w, "Invalid token", "invalid_token", http.StatusUnauthorized)
				return
			}

			jwtToken, err := jwt.DecodeJWT(token)
			if err != nil {
				logger.Log().Error(err.Error(), zap.String("token", token))
				response.ErrorJson(w, "Invalid token", "decode_token_error", http.StatusUnauthorized)
				return
			}

			userID, err := jwtToken.Get("user_id")

			if err != nil {
				logger.Log().Error(err.Error(), zap.Any("jwt_token", jwtToken))

				response.ErrorJson(w, "Invalid token", "decode_token_error", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), enums.ContextKeyUserId, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// IsAuthenticated checks if the user is authenticated by looking for user_id in the context.
func IsAuthenticated(r *http.Request) bool {
	userID := r.Context().Value(enums.ContextKeyUserId)
	return userID != nil
}
