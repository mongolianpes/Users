package embedding

import (
	"database/sql"
	"fmt"
	"testing"
	"users/internal/db"
)

func connectTestToDB() (*sql.DB, error) {
	host := "localhost"
	port := "5432"
	user := "postgres"
	password := "123"
	dbname := "project_farm"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}

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
