package repository

import (
	"context"

	"github.com/biruk/bus-ticket/auth-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type userRepository struct {
	queries *Queries
}

func NewUserRepository(queries *Queries) domain.UserRepository {
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	params := CreateUserParams{
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     u.FullName,
		PhoneNumber:  u.PhoneNumber,
		Role:         string(u.Role),
	}

	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainUser(user), nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return toDomainUser(user), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return toDomainUser(user), nil
}

func (r *userRepository) UpdateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	params := UpdateUserParams{
		ID: pgtype.UUID{Bytes: u.ID, Valid: true},
		FullName: pgtype.Text{
			String: u.FullName,
			Valid:  u.FullName != "",
		},
		PhoneNumber: pgtype.Text{
			String: u.PhoneNumber,
			Valid:  u.PhoneNumber != "",
		},
		Role: pgtype.Text{
			String: string(u.Role),
			Valid:  u.Role != "",
		},
	}

	user, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainUser(user), nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func toDomainUser(u User) *domain.User {
	return &domain.User{
		ID:           uuid.UUID(u.ID.Bytes),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     u.FullName,
		PhoneNumber:  u.PhoneNumber,
		Role:         domain.Role(u.Role),
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}
