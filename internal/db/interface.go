package db

import (
	"context"
	pb "users/proto"
)

type UsersStorage interface {
	GetUserInfoByLogin(ctx context.Context, userLogin string) (*pb.GetUserInfoResponse, error)
	GetUserInfoByID(ctx context.Context, userID int) (*pb.GetUserInfoResponse, error)
	GetUserInfoForAuth(ctx context.Context, userLogin string) (*pb.AuthResponse, string, error)
	RegisterUser(ctx context.Context, login, name, hashedPassword string, embeddingText string, embedding []float64) (int, error)
	AddAvatar(ctx context.Context, userID int, avatarPath string) error
	DeleteUser(ctx context.Context, userID int) error
	GetUserIDByLogin(ctx context.Context, login string) (int, error)
	GetSavedEmbeddingTexts(ctx context.Context, offset, limit int) (map[int]string, error)
	SaveEmbeddingAfterRetryGenerate(ctx context.Context, rowID int, embedding []float64) error
}
