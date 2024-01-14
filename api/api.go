package api

import (
	"context"
	"fmt"
	"time"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/engine"
	"github.com/act-gpt/marino/engine/embedding"
	"github.com/act-gpt/marino/engine/parser"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var Client *Api

type Api struct {
	Config system.SystemConfig
	ctx    context.Context
}

func NewApiClient() *Api {
	ctx := context.Background()
	Client = &Api{
		Config: system.Config,
		ctx:    ctx,
	}
	return Client
}

func (api Api) Embedding(text []string) ([]types.Embedding, error) {
	res, err := embedding.Request(text)
	if err != nil {
		return nil, err
	}
	var embeddings []types.Embedding
	for _, data := range res.Data {
		embeddings = append(embeddings, data.Embedding)
	}
	return embeddings, nil
}

func (api Api) Parse(filename string) (parser.Sugmentation, error) {
	res, err := parser.Document(filename)
	if err != nil {
		return parser.Sugmentation{}, err
	}
	return res, nil
}

func systemPrompt() string {
	now := time.Now()
	date := now.Format("2006-01-02 15:04:05")
	prompt := system.Config.SystemPrompt
	if prompt == "" {
		prompt = config.SYSTEM_PROMPT
	}
	prompt = fmt.Sprintf(`%s Now is %s.`, prompt, date)
	return prompt
}

// for conversation
func (api Api) BuildConversion(messages []model.Message) []types.ChatModelMessage {
	var msgs []types.ChatModelMessage
	msgs = append(msgs, types.ChatModelMessage{
		Role:    types.ChatMessageRoleSystem,
		Content: systemPrompt(),
	})
	for _, val := range messages {
		msgs = append(msgs, types.ChatModelMessage{
			Role:    types.ChatMessageRoleUser,
			Content: val.Question,
		})
		msgs = append(msgs, types.ChatModelMessage{
			Role:    types.ChatMessageRoleAssistant,
			Content: val.Answer,
		})
	}
	return msgs
}

func (api Api) Engine(bot model.BotSetting) (engine.LLM, error) {
	return engine.New(bot)
}

// for knowledge query
func (api Api) BuildQuery(query string, segments []model.Segment, messages []model.Message, bot model.BotSetting) []types.ChatModelMessage {
	// 上下文
	temp := common.PromptTemplate(config.QUESTION_TEMPLATE)
	contexts, _ := temp.Render(struct {
		Contexts []model.Segment
	}{
		Contexts: segments,
	})

	temp = common.PromptTemplate(config.HISTORIES_TEMPLATE)
	histories, _ := temp.Render(struct {
		Histories []model.Message
	}{
		Histories: messages,
	})

	temp = common.PromptTemplate(config.COMPLETION_PROMPT)
	str, _ := temp.Render(struct {
		Prompt    string
		Context   string
		Histories string
		Query     string
	}{
		Prompt:    bot.Prompt,
		Context:   contexts,
		Histories: histories,
		Query:     query,
	})

	var msgs []types.ChatModelMessage
	msgs = append(msgs, types.ChatModelMessage{
		Role:    types.ChatMessageRoleSystem,
		Content: systemPrompt(),
	})
	msgs = append(msgs, types.ChatModelMessage{
		Role:    types.ChatMessageRoleUser,
		Content: str,
	})

	return msgs
}

func (api Api) Insert(document *types.Document, update bool) error {

	processor := common.NewPreprocessor(&common.PreprocessorConfig{})
	// 获取分块
	chunks, err := processor.Preprocess(document)
	if err != nil {
		return err
	}
	if update {
		api.Delete([]string{document.ID})
	}
	num := 0
	for batch := range genBatches(chunks, 10) {
		var list []string
		num += len(batch)
		for _, data := range batch {
			list = append(list, data.Text)
		}
		// 生成向量
		embeddings, _ := api.Embedding(list)
		for i, data := range batch {
			segument := &model.Segment{
				Id:          data.ID,
				KnowledgeId: data.DocumentID,
				Embedding:   embeddings[i],
				Index:       i + 1,
				Text:        data.Text,
				Corpus:      data.Metadata.Corpus,
				Source:      data.Metadata.Source,
				Url:         data.Metadata.Url,
			}
			segument.Insert()
		}
	}
	logx.Info(fmt.Sprintf("Embedding with %s, Length %d", document.ID, num))
	return nil
}

func (api Api) Get(id string) (model.Segment, error) {
	return model.FindSegment(id)
}

func (api Api) Delete(items []string) error {
	return model.DeleteSegments(items)
}

func (api Api) Query(question string, bot model.BotSetting) ([]model.Segment, error) {
	embeddings, err := api.Embedding([]string{question})
	if err != nil {
		return []model.Segment{}, err
	}
	return model.QueryEmbedding(embeddings[0], bot.Corpus, bot.Contexts, bot.Score/100)
}

func genBatches(chunks map[string][]*types.Chunk, size int) <-chan []*types.Chunk {
	ch := make(chan []*types.Chunk)
	go func() {
		var batch []*types.Chunk

		for _, chunkList := range chunks {
			for _, chunk := range chunkList {
				batch = append(batch, chunk)

				if len(batch) == size {
					// Reach the batch size, copy and send all the buffered chunks.
					temp := make([]*types.Chunk, size)
					copy(temp, batch)
					ch <- temp

					// Clear the buffer.
					batch = batch[:0]
				}
			}
		}
		// Send all the remaining chunks, if any.
		if len(batch) > 0 {
			ch <- batch
		}

		close(ch)
	}()

	return ch
}
