package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestInit(t *testing.T) {
	// Create a temporary .env file for testing
	envContent := []byte("JWT_SECRET_KEY=test_secret_key\n")
	err := os.WriteFile(".env", envContent, 0644)
	if err != nil {
		t.Fatalf("Error creating .env file: %v", err)
	}
	defer os.Remove(".env")

	// Test loading .env file
	err = godotenv.Load()
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	// Test case where .env file doesn't exist
	err = os.Remove(".env")
	if err != nil {
		t.Fatalf("Error removing .env file: %v", err)
	}

	err = godotenv.Load()
	if err == nil {
		t.Error("Expected error loading .env file, got nil")
	}
}

func TestCreateToken(t *testing.T) {
	// Set JWT_SECRET_KEY for testing environment
	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET_KEY")

	// Test case 1: Valid token creation
	username := "test_user"
	token, err := createToken(username)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if token == "" {
		t.Error("Expected non-empty token")
	}

	// Test case 2: Missing JWT_SECRET_KEY
	os.Unsetenv("JWT_SECRET_KEY")
	_, err = createToken(username)
	if err == nil {
		t.Error("Expected error due to missing JWT_SECRET_KEY, got nil")
	}
}
