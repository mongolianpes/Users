package main

import (
	"log/slog"
	"net"

	pb "users/proto"
	"users/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		slog.Error("не удалось слушать порт", "error", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUsersServer(grpcServer, service.NewUsersServer())

	slog.Info("Сервер запущен")
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("Ошибка при запуске сервера", "error", err)
	}
}
