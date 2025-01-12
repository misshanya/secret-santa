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
		tokenString, err := getTokenFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := parseAndValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTokenFromRequest(r *http.Request) (string, error) {
	var tokenString string

	authHeader := r.Header.Get("Authorization")

	if authHeader != "" {
		tokenString = getTokenFromHeader(authHeader)
	} else {
		cookie, err := r.Cookie("token")
		if err != nil {
			return "", fmt.Errorf("Unauthorized: no token")
		}
		tokenString = cookie.Value
	}

	return tokenString, nil
}

func getTokenFromHeader(AuthHeader string) string {
	return strings.Split(AuthHeader, " ")[1]
}

func parseAndValidateToken(tokenString string) (int, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.GetConfig().JWTSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok && float64(time.Now().Unix()) > exp {
			return 0, fmt.Errorf("Unauthorized: Token has expired")
		}

		if idFloat, ok := claims["id"].(float64); ok {
			return int(idFloat), nil
		}
		return 0, fmt.Errorf("Unauthorized: Invalid token claims")
	}

	return 0, fmt.Errorf("Unauthorized: Invalid token")
}
