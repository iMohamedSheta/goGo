package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/iMohamedSheta/xapp/app/enums"
	"github.com/iMohamedSheta/xapp/pkg/auth"
	"github.com/iMohamedSheta/xapp/pkg/logger"
	"github.com/iMohamedSheta/xapp/pkg/response"

	"go.uber.org/zap/zapcore"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				response.ErrorJson(w, "Missing or invalid Authorization header", "missing_auth_header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the access token
			jwtToken, err := auth.ValidateAuthToken(token, auth.AccessToken)
			if err != nil {
				response.ErrorJson(w, err.Error(), "invalid_token", http.StatusUnauthorized)
				return
			}

			// Now you have the decoded JWT token, and you can retrieve the user_id or other information
			userID, err := jwtToken.Get("user_id")
			if err != nil || userID == nil {
				response.ErrorJson(w, "Invalid token", "decode_token_error", http.StatusUnauthorized)
				logger.Log().Error(err.Error(), zapcore.Field{Key: "user_id", String: userID.(string)})
				return
			}

			// Set user_id in the request context for use in downstream handlers
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
