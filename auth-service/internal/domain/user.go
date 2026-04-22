package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RolePassenger       Role = "ROLE_PASSENGER"
	RoleStaff           Role = "ROLE_STAFF"
	RoleBusWorkerAdmin  Role = "ROLE_BUS_WORKER_ADMIN"
	RoleSuperAdmin      Role = "ROLE_SUPER_ADMIN"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	PhoneNumber  string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type AuthService interface {
	Register(ctx context.Context, email, password, fullName, phoneNumber string, role Role) (*User, string, error)
	Login(ctx context.Context, email, password string) (*User, string, error)
	ValidateToken(ctx context.Context, token string) (uuid.UUID, Role, error)
}
