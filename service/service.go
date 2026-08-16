package service

import (
	"context"
	"errors"
	"log/slog"

	"users/crypto"
	"users/embedding"
	pb "users/proto"
)

const (
	myLoginAlias        = "my"
	allowedInterestsLen = 250
)

func (s *UsersServer) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	if req.UserLogin == "" {
		return nil, errors.New("userLogin пустой")
	}

	return GetUserInfoFromDB(s.db, req.UserLogin)
}

func (s *UsersServer) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	response, currentPassword, err := GetUserInfoForAuthFromDB(s.db, req.Login)
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

	if len(req.Interests) > allowedInterestsLen {
		return nil, errors.New("Допустимая длина интересов: 250 символов")
	}

	userID, err := RegisterUserInDB(s.db, req.Login, req.Name, hashedPassword)
	if err != nil {
		return nil, err
	}

	go func() {
		err := embedding.InsertEmbedding(s.db, userID, req.Interests, "UPDATE users SET embedding = $1::float8[] WHERE user_id = $1")
		if err != nil {
			slog.Error("Ошибка при создании и вставки эмбеддинга для пользователя", "error", err)
		}
	}()

	return &pb.RegisterResponse{}, nil
}

func (s *UsersServer) AddAvatar(ctx context.Context, req *pb.AddAvatarRequest) (*pb.AddAvatarResponse, error) {
	if err := AddUserAvatarToDB(s.db, int(req.UserID), req.AvatarPath); err != nil {
		return nil, err
	}

	return &pb.AddAvatarResponse{}, nil
}
