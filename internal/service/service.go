package service

import (
	"context"
	"errors"
	"log/slog"

	"users/internal/crypto"
	"users/internal/embedding"
	pb "users/proto"
)

const (
	myLoginAlias           = "my"
	allowedInterestsMaxLen = 250
	allowedInterestsMinLen = 10
)

func (s *UsersServer) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	if req.UserID != 0 {
		resp, err := s.storage.GetUserInfoByID(ctx, int(req.UserID))
		if err != nil {
			slog.Warn("Попытка получения информации о пользователе привела к ошибке", "error", err)
			return nil, err
		}
		return resp, nil
	}

	if req.UserLogin != "" {
		resp, err := s.storage.GetUserInfoByLogin(ctx, req.UserLogin)
		if err != nil {
			slog.Warn("Попытка получения информации о пользователе привела к ошибке", "error", err)
			return nil, err
		}
		return resp, nil
	}

	slog.Info("Успешное получение информации о пользователе", "id", req.UserID, "login", req.UserLogin)

	return nil, errors.New("userLogin и userID пустые")

}

func (s *UsersServer) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	response, currentPassword, err := s.storage.GetUserInfoForAuth(ctx, req.Login)
	if err != nil {
		slog.Warn("Ошибка при авторизации", "error", err, "userLogin", req.Login)
		return nil, err
	}

	if !crypto.VerifyHash(req.Password, currentPassword) {
		slog.Warn("Попытка авторизации под неверным паролем", "userLogin", req.Login)
		return nil, errors.New("Неверный логин или пароль")
	}

	slog.Info("Успешная авторизация", "login", req.Login)

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

	var userID int
	embeddigs, err := embedding.GenerateEmbeddingForUser(ctx, req.Interests)
	if err != nil {
		userID, err = s.storage.RegisterUser(ctx, req.Login, req.Name, hashedPassword, req.Interests, []float64{})
		if err != nil {
			slog.Warn("Ошибка при регистрации", "login", req.Login)
			return nil, err
		}
	} else {
		userID, err = s.storage.RegisterUser(ctx, req.Login, req.Name, hashedPassword, "", embeddigs)
		if err != nil {
			slog.Warn("Ошибка при регистрации", "login", req.Login)
			return nil, err
		}
	}

	slog.Info("Успешная регистрация", "login", req.Login)

	return &pb.RegisterResponse{
		UserID: int64(userID),
	}, nil
}

func (s *UsersServer) AddAvatar(ctx context.Context, req *pb.AddAvatarRequest) (*pb.AddAvatarResponse, error) {
	if err := s.storage.AddAvatar(ctx, int(req.UserID), req.AvatarPath); err != nil {
		slog.Warn("Ошибка при добавлении аватара", "userID", req.UserID)
		return nil, err
	}

	slog.Info("Успешное добавление автара", "userID", req.UserID)

	return &pb.AddAvatarResponse{}, nil
}

func (s *UsersServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.storage.DeleteUser(ctx, int(req.UserID)); err != nil {
		slog.Warn("Ошибка при удалении пользователя", "id", req.UserID)
		return nil, err
	}

	slog.Info("Успешное удаление пользователя", "id", req.UserID)

	return &pb.DeleteUserResponse{}, nil
}

func (s *UsersServer) GetUserID(ctx context.Context, req *pb.GetUserIDRequest) (*pb.GetUserIDResponse, error) {
	userID, err := s.storage.GetUserIDByLogin(ctx, req.Login)
	if err != nil {
		slog.Warn("Ошибка при получении ID пользователя", "id", req.Login)
		return nil, err
	}

	return &pb.GetUserIDResponse{
		ID: int64(userID),
	}, nil
}
