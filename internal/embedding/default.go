package embedding

import (
	"context"
	"log/slog"
	"time"

	"users/internal/db"
)

func RunSetterDefaultEmbedding(storage db.UsersStorage) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	for {
		err := storage.UpdateDefaultEmbeddingUsers(ctx)
		if err != nil {
			slog.Warn("Не удалось обновить дефолтный эмбеддинг для пользователей", "err", err)
		}

		time.Sleep(time.Hour * 23)
	}
}
