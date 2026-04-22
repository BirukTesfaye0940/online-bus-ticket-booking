package grpc

import (
	"context"

	"github.com/biruk/bus-ticket/auth-service/internal/domain"
	pb "github.com/biruk/bus-ticket/auth-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, _, err := h.authService.Register(ctx, req.Email, req.Password, req.FullName, req.PhoneNumber, domain.Role(req.Role.String()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.RegisterResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}

	return &pb.LoginResponse{
		AccessToken: token,
		User:        toProtoUser(user),
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, role, err := h.authService.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID.String(),
		Role:   pb.Role(pb.Role_value[string(role)]),
	}, nil
}

func toProtoUser(u *domain.User) *pb.User {
	return &pb.User{
		Id:          u.ID.String(),
		Email:       u.Email,
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber,
		Role:        pb.Role(pb.Role_value[string(u.Role)]),
		CreatedAt:   u.CreatedAt.String(),
		UpdatedAt:   u.UpdatedAt.String(),
	}
}
