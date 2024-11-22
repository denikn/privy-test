package main

import (
    "beprivytest/internal/adapter/http"
    "beprivytest/internal/adapter/storage"
    "beprivytest/internal/application"
    "beprivytest/internal/middleware"
    "fmt"
    serve "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/golang-jwt/jwt/v4"
    "time"
)

func init() {
    if err := godotenv.Load(); err != nil {
        fmt.Println("Error loading .env file")
    }
}

// just for testing
func createToken(username string) (string, error) {
    secretKey := os.Getenv("JWT_SECRET_KEY")
    if secretKey == "" {
        return "", fmt.Errorf("JWT_SECRET_KEY is not set")
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &jwt.StandardClaims{
        Subject:   username,
        ExpiresAt: jwt.NewNumericDate(expirationTime).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

func main() {
    baseDir := "./uploads"
    os.MkdirAll(baseDir, os.ModePerm)

    fileStorage := storage.NewFileStorage(baseDir)
    imageService := application.NewImageService(fileStorage)
    imageHandler := http.NewImageHandler(imageService)

    router := mux.NewRouter()

    // Middleware JWT just for /processed-image
    router.HandleFunc("/upload", imageHandler.UploadImage).Methods("POST")
    router.HandleFunc("/faces-count", imageHandler.GetFacesCount).Methods("GET")
    router.HandleFunc("/processed-image", middleware.JWTMiddleware(imageHandler.ServeProcessedImage)).Methods("GET")

    fmt.Println("Server is starting on port 8080...")
    // generate token to access 3rd endpoint
    token, err := createToken(os.Getenv("ALLOWED_USERNAME"))
    if err != nil {
        fmt.Println("Error creating token:", err)
        return
    }
    fmt.Println("Generated Token:", token)

    err = serve.ListenAndServe(":8080", router)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}
