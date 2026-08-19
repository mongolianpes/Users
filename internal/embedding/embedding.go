package embedding

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"users/internal/db"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var client = &http.Client{
	Timeout: 30 * time.Second,
}

const userAdaptationRate = 0.9

type ollamaJsonResponse struct {
	Embedding []float64
}

type ollamaJsonRequest struct {
	Model  string
	Prompt string
}

var ollamaHost = os.Getenv("OLLAMA_HOST")

func GenerateEmbeddingForUser(storage db.UsersStorage, ctx context.Context, rowID int, text string) error {
	if ollamaHost == "" {
		return errors.New("Переменная OLLAMA_HOST должна иметь значение: адрес локальной нейросети ollama")
	}

	req := ollamaJsonRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := client.Post(ollamaHost+"/api/embeddings", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Ollama вернул статус код не 200")
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var result ollamaJsonResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if err := storage.SaveEmbedding(ctx, rowID, result.Embedding); err != nil {
		return err
	}

	return nil
}
