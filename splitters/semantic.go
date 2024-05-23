package splitters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"
)

// Preprocessor splits a list of documents into chunks.
type SemanticPreprocess struct {
	cfg *PreprocessorConfig
}

func SemanticPreprocessor(cfg *PreprocessorConfig) *SemanticPreprocess {
	return &SemanticPreprocess{
		cfg: cfg.Init(),
	}
}

func (p *SemanticPreprocess) Preprocess(doc types.Document, bot model.BotSetting) (map[string][]*types.Chunk, []string, error) {
	chunkMap := make(map[string][]*types.Chunk)
	docID := doc.ID
	meta := doc.Metadata
	if docID == "" {
		docID = common.GetUUID()
	}
	setting, err := bot.MergeSetting(system.Config)
	if err != nil {
		return nil, nil, err
	}

	textChunks, codes, err := p.request(doc.Text, setting)
	if err != nil {
		return nil, nil, err
	}
	for _, textChunk := range textChunks {
		id := common.GetUUID()
		// return chunks
		chunkMap[docID] = append(chunkMap[docID], &types.Chunk{
			ID:         id,
			DocumentID: docID,
			Text:       textChunk,
			Metadata:   meta,
		})
	}
	return chunkMap, codes, nil
}

// request to llm
func (p *SemanticPreprocess) request(doc string, setting model.BotSetting) ([]string, []string, error) {
	conf := system.Config.Chunk
	reqUrl := conf.Host + conf.Api
	key := conf.AccessKey
	var chunk = setting.Chunk
	llm := chunk.Embedding
	max := chunk.MaxTokens
	min := chunk.MinTokens
	overlap := chunk.Overlap
	semantic := chunk.Semantic
	if max > 1000 {
		overlap = 400
	}
	item := struct {
		Embedding string `json:"embedding"`
		MaxTokens int    `json:"max_tokens"`
		MinTokens int    `json:"min_tokens"`
		Input     string `json:"input"`
		Overlap   int    `json:"overlap"`
		Semantic  bool   `json:"semantic"`
	}{
		Embedding: llm,
		MaxTokens: max,
		Overlap:   overlap,
		MinTokens: min,
		Input:     doc,
		Semantic:  semantic,
	}
	body, err := json.Marshal(item)
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader(body))

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+key)
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	type Res struct {
		Data struct {
			Segments []string `json:"segments"`
			Codes    []string `json:"codes"`
		} `json:"data"`
	}
	var resp Res
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("error", err)
		return nil, nil, err
	}
	return resp.Data.Segments, resp.Data.Codes, nil
}
