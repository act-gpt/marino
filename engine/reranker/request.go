package reranker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/act-gpt/marino/config/system"
)

type RerankerRequest struct {
	Model     string   `json:"model"`
	Top       int      `json:"top"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
}

type RerankerResponse struct {
	Data struct {
		Scores    []float64 `json:"scores"`
		Ids       []int     `json:"ids"`
		Documents []string  `json:"documents"`
	} `json:"data"`
}

func Reranker(query string, documents []string, top int) (RerankerResponse, error) {

	conf := system.Config.Reranker
	reqUrl := conf.Host + conf.Api

	model := conf.Model
	key := conf.AccessKey

	item := RerankerRequest{
		Model:     model,
		Top:       top,
		Query:     query,
		Documents: documents,
	}

	body, err := json.Marshal(item)
	if err != nil {
		return RerankerResponse{}, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader(body))

	if err != nil {
		fmt.Println(err)
		return RerankerResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	if key != "" {
		req.Header.Add("Authorization", "Bearer "+key)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return RerankerResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return RerankerResponse{}, err
	}
	var resp RerankerResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println(err)
		return RerankerResponse{}, err
	}
	return resp, nil
}
