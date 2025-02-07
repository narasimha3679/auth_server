package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"auth-server/config"
	"auth-server/models"
	"auth-server/utils"
)

type RegisterInput struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	CountryCode *string `json:"country_code,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RequestResetInput struct {
	Email string `json:"email" binding:"required,email"`
}

type PhoneLoginInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	CountryCode string `json:"country_code" binding:"required"`
}

type VerifyPhoneOtpInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	CountryCode string `json:"country_code" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

type ResetPasswordInput struct {
	ResetCode string `json:"reset_code" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:     input.Email,
		Password:  string(hashedPassword),
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	if input.PhoneNumber != nil {
		user.PhoneNumber = *input.PhoneNumber
	}
	if input.CountryCode != nil {
		user.CountryCode = *input.CountryCode
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(c *gin.Context) {
	var input RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.ValidateToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Verify it's a refresh token
	if refresh, ok := claims["refresh"].(bool); !ok || !refresh {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	userId := uint(claims["user_id"].(float64))
	newAccessToken, err := utils.GenerateToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

func RequestPasswordReset(c *gin.Context) {
	var input RequestResetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify email exists
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	// Generate reset code
	resetCode, err := utils.GenerateResetCode(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset code"})
		return
	}

	// In a production environment, you would send this code via email
	// For demo purposes, we'll return it in the response
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset code generated",
		"code":    resetCode,
	})
}

func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate reset code and get email
	email, err := utils.ValidateResetCode(input.ResetCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update user password
	result := config.DB.Model(&models.User{}).Where("email = ?", email).Update("password", string(hashedPassword))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func PhoneLogin(c *gin.Context) {
	var input PhoneLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format full phone number with country code
	fullPhone := input.CountryCode + input.PhoneNumber

	// Check if user exists with this phone number
	var user models.User
	if err := config.DB.Where("phone_number = ? AND country_code = ?",
		input.PhoneNumber, input.CountryCode).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No user found with this phone number"})
		return
	}

	// Generate OTP
	otp, err := utils.GenerateOTP(fullPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	// In production, send OTP via SMS
	// For development, return it in response
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
		"otp":     otp, // Remove this in production
	})
}

func VerifyPhoneOtp(c *gin.Context) {
	var input VerifyPhoneOtpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format full phone number
	fullPhone := input.CountryCode + input.PhoneNumber

	// Validate OTP
	valid, err := utils.ValidateOTP(fullPhone, input.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Verify the stored phone matches the submitted phone
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	// Find user
	var user models.User
	if err := config.DB.Where("phone_number = ? AND country_code = ?",
		input.PhoneNumber, input.CountryCode).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Phone verified successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
