package embedding

import (
	"context"
	"database/sql"
	"time"
	"users/internal/db"
)

func RetryInsertEmbeddings(storage db.UsersStorage) {
	for {
		runInsertSavedEmbeddings(storage)
		time.Sleep(time.Hour * 4)
	}
}

func runInsertSavedEmbeddings(storage db.UsersStorage) error {
	offset := 0
	limit := 5

	for {
		offset += 5
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		embeddings, err := storage.GetSavedEmbeddingTexts(ctx, offset, limit)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return err
		}

		if len(embeddings) == 0 {
			break
		}

		for rowID, text := range embeddings {
			if err := GenerateEmbeddingForUser(storage, ctx, rowID, text); err != nil {
				return err
			}

			if err := storage.DeleteSavedEmbeddingText(ctx, rowID); err != nil {
				return err
			}
		}
	}

	return nil
}
