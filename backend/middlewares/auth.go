package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/misshanya/secret-santa/config"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			tokenString = getTokenFromHeader(authHeader)
		} else {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString = cookie.Value
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.GetConfig().JWTSecret), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok && float64(time.Now().Unix()) > exp {
				http.Error(w, "Unauthorized: Token has expired", http.StatusUnauthorized)
				return
			}

			var userID int
			if idFloat, ok := claims["id"].(float64); ok {
				userID = int(idFloat)
			} else {
				http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		}
	})
}

func getTokenFromHeader(AuthHeader string) string {
	return strings.Split(AuthHeader, " ")[1]
}
