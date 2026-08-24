package embedding

import (
	"database/sql"
	"testing"
	"users/internal/db"
)

func addSavedEmbeddings(db *sql.DB, users []int) error {
	if _, err := db.Exec("DELETE FROM embeddings_users"); err != nil {
		return err
	}

	_, err := db.Exec(`INSERT INTO embeddings_users (user_id, text) VALUES
		($1, 'text'),
		($2, 'text')`, users[0], users[1])
	if err != nil {
		return err
	}

	return nil
}

func TestRunInsertSavedEmbeddings(t *testing.T) {
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

	if err := addSavedEmbeddings(db, updUsers); err != nil {
		t.Error(err)
	}

	if err := runInsertSavedEmbeddings(storage); err != nil {
		t.Error(err)
	}

	var text string
	if err := db.QueryRow("SELECT text FROM embeddings_users WHERE user_id = $1", updUsers[0]).Scan(&text); err != nil {
		t.Error(err)
	}

	var embedding []float64
	if err := db.QueryRow("SELECT embedding FROM users WHERE user_id = $1", updUsers[2]).Scan(&embedding); err == nil {
		t.Error(err)
	}
}
