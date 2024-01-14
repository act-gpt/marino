package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/act-gpt/marino/config/system"

	"github.com/sashabaranov/go-openai"
)

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

func Request(input []string) (openai.EmbeddingResponse, error) {

	conf := system.Config.Embedding
	api := conf.Host + conf.Api
	model := conf.Model
	key := conf.AccessKey

	if api == "" || model == "" {
		return openai.EmbeddingResponse{}, fmt.Errorf("Api and Model must must be set before you use embedding model")
	}
	// for azure.com
	if strings.Contains(api, "azure.com") {
		api += "?api-version=2023-05-15"
	}

	item := EmbeddingRequest{
		Model: model,
		Input: input,
	}

	body, err := json.Marshal(item)
	if err != nil {
		return openai.EmbeddingResponse{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api, bytes.NewReader(body))

	if err != nil {
		fmt.Println(err)
		return openai.EmbeddingResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	if key != "" {
		if strings.Contains(api, "azure.com") {
			req.Header.Add("Api-Key", key)
		} else {
			req.Header.Add("Authorization", "Bearer "+key)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return openai.EmbeddingResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return openai.EmbeddingResponse{}, err
	}
	var embedding openai.EmbeddingResponse
	err = json.Unmarshal(body, &embedding)
	if err != nil {
		fmt.Println(err)
		return openai.EmbeddingResponse{}, err
	}

	return embedding, nil
}
