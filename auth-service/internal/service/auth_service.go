package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/biruk/bus-ticket/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID   `json:"user_id"`
	Role   domain.Role `json:"role"`
}

type authService struct {
	repo          domain.UserRepository
	jwtSecret     string
	tokenDuration time.Duration
}

func NewAuthService(repo domain.UserRepository, jwtSecret string, tokenDuration time.Duration) domain.AuthService {
	return &authService{
		repo:          repo,
		jwtSecret:     jwtSecret,
		tokenDuration: tokenDuration,
	}
}

func (s *authService) Register(ctx context.Context, email, password, fullName, phoneNumber string, role domain.Role) (*domain.User, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		PhoneNumber:  phoneNumber,
		Role:         role,
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.generateToken(createdUser)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return createdUser, token, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (uuid.UUID, domain.Role, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, "", fmt.Errorf("invalid token: %w", err)
	}

	if c, ok := token.Claims.(*claims); ok && token.Valid {
		return c.UserID, c.Role, nil
	}

	return uuid.Nil, "", errors.New("invalid token claims")
}

func (s *authService) generateToken(user *domain.User) (string, error) {
	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(s.jwtSecret))
}
