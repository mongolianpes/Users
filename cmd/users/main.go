package main

import (
	"log/slog"
	"net"

	"users/internal/db"
	pb "users/internal/proto"
	"users/internal/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	storage, err := db.NewPostgresStorage()
	if err != nil {
		panic(err)
	}
	pb.RegisterUsersServer(grpcServer, service.NewUsersServer(storage))

	slog.Info("Сервер запущен")
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("Ошибка при запуске сервера", "error", err)
	}
}
