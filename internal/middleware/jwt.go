package middleware

import (
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v4"
    "os"
    "fmt"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		if len(jwtKey) == 0 {
			fmt.Println("JWT_SECRET_KEY is not set")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Println("Missing authorization header")
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			fmt.Println("Malformed token")
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			fmt.Println("Error parsing token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			fmt.Println("Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}
