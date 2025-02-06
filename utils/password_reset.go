package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"auth-server/config"
)

const (
	resetCodeLength = 32
	resetCodeExpiry = time.Minute * 15
)

// GenerateResetCode creates a random reset code and stores it in Redis
func GenerateResetCode(email string) (string, error) {
	// Generate random bytes
	bytes := make([]byte, resetCodeLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate reset code: %v", err)
	}

	resetCode := hex.EncodeToString(bytes)

	// Store in Redis with expiration
	ctx := context.Background()
	key := fmt.Sprintf("pwd_reset:%s", resetCode)
	if err := config.RedisClient.Set(ctx, key, email, resetCodeExpiry).Err(); err != nil {
		return "", fmt.Errorf("failed to store reset code: %v", err)
	}

	return resetCode, nil
}

// ValidateResetCode checks if the reset code is valid and returns the associated email
func ValidateResetCode(resetCode string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("pwd_reset:%s", resetCode)

	email, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("invalid or expired reset code")
	}

	// Delete the code after successful validation
	config.RedisClient.Del(ctx, key)

	return email, nil
}
