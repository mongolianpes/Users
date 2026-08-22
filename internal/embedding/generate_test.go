package embedding

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"testing"
	"time"

	"users/internal/crypto"
	"users/internal/db"
)

func updateTestUsers(db *sql.DB) ([]int, error) {
	if _, err := db.Exec("DELETE FROM users"); err != nil {
		return []int{}, err
	}

	hashedPass, err := crypto.HashString("1234567890")
	if err != nil {
		return []int{}, err
	}

	rowsUserIDs, err := db.Query(`INSERT INTO users (login, name, password) VALUES
		(1, '1', $1),
		(2, '2', $1),
		(3, '3', $1)
		RETURNING user_id`, hashedPass)
	if err != nil {
		return []int{}, err
	}

	var updUsers []int
	var userID1 int
	var userID2 int
	var userID3 int
	rowsUserIDs.Next()
	rowsUserIDs.Scan(&userID1)
	rowsUserIDs.Next()
	rowsUserIDs.Scan(&userID2)
	rowsUserIDs.Next()
	rowsUserIDs.Scan(&userID3)
	updUsers = append(updUsers, userID1)
	updUsers = append(updUsers, userID2)
	updUsers = append(updUsers, userID3)

	var userIDAvatar int
	err = db.QueryRow(`INSERT INTO users (login, name, password, avatar_path) VALUES
	('withavatar', 'avatarka', $1, 'path-to-avatar') RETURNING user_id`, hashedPass).Scan(&userIDAvatar)
	if err != nil {
		return []int{}, err
	}
	updUsers = append(updUsers, userIDAvatar)

	return updUsers, err
}

func TestGenerateEmbeddingForUser(t *testing.T) {
	ollamaHost = "http://localhost:11434"

	storage, err := db.NewPostgresStorage()
	if err != nil {
		t.Error(err)
	}

	db, err := connectTestToDB()
	if err != nil {
		t.Error(err)
	}

	updUsers, err := updateTestUsers(db)
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := GenerateEmbeddingForUser(storage, ctx, updUsers[0], "embedding"); err != nil {
		t.Error(err)
	}

	var embedding []float64
	var raw sql.NullString
	err = db.QueryRow("SELECT embedding FROM users WHERE user_id=$1", updUsers[0]).Scan(&raw)
	if err != nil {
		t.Error(err)
	}

	if raw.Valid {
		parts := strings.Trim(raw.String, "[]")
		nums := strings.Split(parts, ",")
		embedding = make([]float64, len(nums))
		for i, s := range nums {
			embedding[i], _ = strconv.ParseFloat(strings.TrimSpace(s), 64)
		}
	} else {
		t.Error("Вернул невалид эмбеддинг")
	}

	err = db.QueryRow("SELECT embedding FROM users WHERE user_id=$1", updUsers[3]).Scan(&raw)
	if err != nil {
		t.Error(err)
	}

	if raw.Valid {
		parts := strings.Trim(raw.String, "[]")
		nums := strings.Split(parts, ",")
		embedding = make([]float64, len(nums))
		for i, s := range nums {
			embedding[i], _ = strconv.ParseFloat(strings.TrimSpace(s), 64)
		}
		t.Error("Вернул валид эмбеддинг для пользователя, у кого эмбеддинга нет")
	}

	ollamaHost = ""

	if err := GenerateEmbeddingForUser(storage, ctx, updUsers[1], "embedding"); err != envOSError {
		t.Error("Если переменная ссылка на ollamaHost пустая возвращает другую ошибку")
	}

	ollamaHost = "http://localhost:32542342"

	if err := GenerateEmbeddingForUser(storage, ctx, updUsers[0], "embedding"); err != cantConnectToOllamaError {
		t.Error("Если ollama не отвечает возвращает другую ошибку")
	}
}
