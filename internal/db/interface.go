package db

import (
	"context"
	pb "users/internal/proto"
)

type UsersStorage interface {
	GetUserInfoByLogin(ctx context.Context, userLogin string) (*pb.GetUserInfoResponse, error)
	GetUserInfoByID(ctx context.Context, userID int) (*pb.GetUserInfoResponse, error)
	GetUserInfoForAuth(ctx context.Context, userLogin string) (*pb.AuthResponse, string, error)
	RegisterUser(ctx context.Context, login, name, hashedPassword string) (int, error)
	AddAvatar(ctx context.Context, userID int, avatarPath string) error
	DeleteUser(ctx context.Context, userID int) error
	GetUserIDByLogin(ctx context.Context, login string) (int, error)
	SaveEmbedding(ctx context.Context, rowID int, embedding []float64) error
}
