package service

import (
	"context"
	"errors"
	"log/slog"

	"users/internal/crypto"
	"users/internal/embedding"
	pb "users/internal/proto"
)

const (
	myLoginAlias           = "my"
	allowedInterestsMaxLen = 250
	allowedInterestsMinLen = 10
)

func (s *UsersServer) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	if req.UserID != 0 {
		return s.storage.GetUserInfoByID(ctx, int(req.UserID))
	}

	if req.UserLogin != "" {
		return s.storage.GetUserInfoByLogin(ctx, req.UserLogin)
	}

	return nil, errors.New("userLogin и userID пустые")

}

func (s *UsersServer) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	response, currentPassword, err := s.storage.GetUserInfoForAuth(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	if !crypto.VerifyHash(req.Password, currentPassword) {
		return nil, errors.New("Неверный логин или пароль")
	}

	return response, nil
}

func (s *UsersServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Login == myLoginAlias {
		return nil, errors.New("Вы не можете зарегистрироваться с данным логином")
	}

	if len(req.Password) <= 9 {
		return nil, errors.New("Длина пароля должна быть больше 9 символов")
	}

	hashedPassword, err := crypto.HashString(req.Password)
	if err != nil {
		return nil, err
	}

	if len(req.Interests) < allowedInterestsMinLen {
		return nil, errors.New("Допустимая длина интересов: от 10 символов")
	}

	if len(req.Interests) > allowedInterestsMaxLen {
		return nil, errors.New("Допустимая длина интересов: до 250 символов")
	}

	userID, err := s.storage.RegisterUser(ctx, req.Login, req.Name, hashedPassword)
	if err != nil {
		return nil, err
	}

	go func() {
		err := embedding.GenerateEmbeddingForUser(s.storage, ctx, userID, req.Interests)
		if err != nil {
			slog.Error("Ошибка при создании и вставки эмбеддинга для пользователя", "error", err)
		}
	}()

	return &pb.RegisterResponse{
		UserID: int32(userID),
	}, nil
}

func (s *UsersServer) AddAvatar(ctx context.Context, req *pb.AddAvatarRequest) (*pb.AddAvatarResponse, error) {
	if err := s.storage.AddAvatar(ctx, int(req.UserID), req.AvatarPath); err != nil {
		return nil, err
	}

	return &pb.AddAvatarResponse{}, nil
}

func (s *UsersServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.storage.DeleteUser(ctx, int(req.UserID)); err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{}, nil
}

func (s *UsersServer) GetUserID(ctx context.Context, req *pb.GetUserIDRequest) (*pb.GetUserIDResponse, error) {
	userID, err := s.storage.GetUserIDByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserIDResponse{
		ID: int32(userID),
	}, nil
}
