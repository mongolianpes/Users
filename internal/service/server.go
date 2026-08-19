package service

import (
	"users/internal/db"
	pb "users/internal/proto"
)

type UsersServer struct {
	pb.UnimplementedUsersServer
	storage db.UsersStorage
}

func NewUsersServer(storage db.UsersStorage) *UsersServer {
	return &UsersServer{
		storage: storage,
	}
}
