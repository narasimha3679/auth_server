package utils

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"auth-server/config"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	resetCodeLength = 6
	resetCodeExpiry = time.Minute * 15
)

// GenerateResetCode creates a 6-digit reset code and stores it in Redis
func GenerateResetCode(email string) (string, error) {
	// Generate 6-digit code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Store in Redis with expiration
	ctx := context.Background()
	key := fmt.Sprintf("pwd_reset:%s", code)
	if err := config.RedisClient.Set(ctx, key, email, resetCodeExpiry).Err(); err != nil {
		return "", fmt.Errorf("failed to store reset code: %v", err)
	}

	return code, nil
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
