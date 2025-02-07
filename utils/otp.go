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

// Generate OTP stores the OTP using the phone number as the key
func GenerateOTP(phone string) (string, error) {
	// Generate 6-digit code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Store in Redis with expiration
	ctx := context.Background()
	key := fmt.Sprintf("phone_otp:%s", phone) // Use phone as key
	if err := config.RedisClient.Set(ctx, key, code, time.Minute*5).Err(); err != nil {
		return "", fmt.Errorf("failed to store OTP: %v", err)
	}

	return code, nil
}

// ValidateOTP checks if the OTP matches the stored value for the phone number
func ValidateOTP(phone, otp string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("phone_otp:%s", phone) // Use phone as key

	storedOTP, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("invalid or expired OTP")
	}

	if storedOTP != otp {
		return false, fmt.Errorf("incorrect OTP")
	}

	// Delete OTP after successful validation
	config.RedisClient.Del(ctx, key)

	return true, nil
}
