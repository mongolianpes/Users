package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	pb "users/internal/proto"
)

type PostgresStorage struct {
	db *sql.DB
}

func (s *PostgresStorage) GetUserInfoByLogin(ctx context.Context, userLogin string) (*pb.GetUserInfoResponse, error) {
	result := &pb.GetUserInfoResponse{}
	if err := s.db.QueryRowContext(ctx, "SELECT name, user_id, avatar_path FROM users WHERE login = $1", userLogin).Scan(&result.Name, &result.UserID, &result.AvatarPath); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, userNotFoundError
		default:
			return nil, err
		}
	}

	result.Login = userLogin

	return result, nil
}

func (s *PostgresStorage) GetUserInfoByID(ctx context.Context, userID int) (*pb.GetUserInfoResponse, error) {
	result := &pb.GetUserInfoResponse{}
	if err := s.db.QueryRowContext(ctx, "SELECT login, name, avatar_path FROM users WHERE user_id = $1", userID).Scan(&result.Login, &result.Name, &result.AvatarPath); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, userNotFoundError
		default:
			return nil, err
		}
	}

	result.UserID = int32(userID)

	return result, nil
}

func (s *PostgresStorage) GetUserInfoForAuth(ctx context.Context, userLogin string) (*pb.AuthResponse, string, error) {
	result := &pb.AuthResponse{}
	var currentPassword string
	if err := s.db.QueryRowContext(ctx, "SELECT user_id, password, name, avatar_path FROM users WHERE login = $1", userLogin).Scan(&result.UserID, &currentPassword, &result.UserName, &result.AvatarPath); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, "", incorrectLoginOrPasswordError
		default:
			return nil, "", err
		}
	}

	return result, currentPassword, nil
}

func (s *PostgresStorage) RegisterUser(ctx context.Context, login, name, hashedPassword string) (int, error) {
	var userID int
	if err := s.db.QueryRowContext(ctx, "INSERT INTO users (login, name, password) VALUES ($1, $2, $3) RETURNING user_id", login, name, hashedPassword).Scan(&userID); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return 0, userNotFoundError
			default:
				return 0, unknownError
			}
		} else {
			return 0, unknownError
		}
	}

	return userID, nil
}

func (s *PostgresStorage) AddAvatar(ctx context.Context, userID int, avatarPath string) error {
	result, err := s.db.ExecContext(ctx, "UPDATE users SET avatar_path = $1 WHERE user_id = $2", avatarPath, userID)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return userNotFoundError
	}
	return nil
}

func (s *PostgresStorage) DeleteUser(ctx context.Context, userID int) error {
	if _, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", userID); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) GetUserIDByLogin(ctx context.Context, login string) (int, error) {
	var userID int
	if err := s.db.QueryRowContext(ctx, "SELECT user_id FROM users WHERE login = $1", login).Scan(&userID); err != nil {
		return userID, err
	}

	return userID, nil
}

func (s *PostgresStorage) SaveEmbedding(ctx context.Context, rowID int, embedding []float64) error {
	if _, err := s.db.ExecContext(ctx, "UPDATE users SET embedding = $1::float8[] WHERE user_id = $2", embedding, rowID); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) SaveEmbeddingText(ctx context.Context, rowID int, text string) error {
	if _, err := s.db.ExecContext(ctx, "INSERT INTO embeddings_users (user_id, text) VALUES ($1, $2)", rowID, text); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) GetSavedEmbeddingTexts(ctx context.Context, offset, limit int) (map[int]string, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT user_id, text FROM embeddings_users OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		return nil, err
	}

	result := make(map[int]string)
	var id int
	var text string
	for rows.Next() {
		if err := rows.Scan(&id, &text); err != nil {
			continue
		}
		result[id] = text
	}

	return result, nil
}

func (s *PostgresStorage) DeleteSavedEmbeddingText(ctx context.Context, rowID int) error {
	if _, err := s.db.ExecContext(ctx, "DELETE FROM embeddings_users WHERE user_id = $1", rowID); err != nil {
		return err
	}
	return nil
}
