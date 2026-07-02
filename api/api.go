package api

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/engine"
	"github.com/act-gpt/marino/engine/embedding"
	"github.com/act-gpt/marino/engine/parser"
	"github.com/act-gpt/marino/splitters"

	//"github.com/act-gpt/marino/engine/reranker"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var Client *Api

var EMBEDDINGS_BATCH_SIZE = 8

type Api struct {
	Config system.SystemConfig
	ctx    context.Context
}

type Split interface {
	NewPreprocessor() *common.Preprocessor
	Preprocess(*types.Document) (map[string][]*types.Chunk, error)
}

func NewApiClient() *Api {
	ctx := context.Background()
	Client = &Api{
		Config: system.Config,
		ctx:    ctx,
	}
	return Client
}

func (api Api) Embedding(text []string, bot model.BotSetting) ([]types.Embedding, error) {
	var embeddings []types.Embedding
	if len(text) == 0 {
		return embeddings, nil
	}
	setting, err := bot.MergeSetting(api.Config)
	if err != nil {
		return nil, err
	}
	res, err := embedding.Request(text, setting.Chunk.Embedding)
	if err != nil {
		return nil, err
	}

	if len(res.Data) != len(text) {
		return nil, fmt.Errorf("embedding response count mismatch: expected %d, got %d", len(text), len(res.Data))
	}
	for _, data := range res.Data {
		if len(data.Embedding) == 0 {
			return nil, fmt.Errorf("embedding response contains empty vector")
		}
		embeddings = append(embeddings, data.Embedding)
	}
	return embeddings, nil
}

func (api Api) Parse(filename string) (types.Sugmentation, error) {
	res, err := parser.Document(filename)
	if err != nil {
		return types.Sugmentation{}, err
	}
	return res, nil
}

func (api Api) Engine(bot model.BotSetting) (engine.LLM, error) {
	return engine.New(bot)
}

// TODO: not finished
func (api Api) Filter(filename string) (types.Sugmentation, error) {
	return types.Sugmentation{}, nil
}

func indexOf(document []types.Document, text string) int {
	for i, item := range document {
		if item.Text == text {
			return i
		}
	}
	return -1
}

func (api Api) Reranker(query string, documents []types.Document, bot model.BotSetting) ([]types.Document, error) {
	var docs []string
	var items []types.Document
	if len(documents) == 0 {
		return items, nil
	}
	for _, doc := range documents {
		docs = append(docs, doc.Text)
	}
	model := api.Config.Reranker.Model

	if bot.RerankModel != "" {
		model = bot.RerankModel
	}

	res, err := engine.Reranker(query, docs, model, bot.Contexts)
	if err != nil {
		return nil, err
	}
	data := res.Data
	for i, item := range data.Documents {
		n := indexOf(documents, item)
		// 小于 0.35 质量已经很低了，过滤掉
		doc := types.Document{
			Text:  item,
			Score: data.Scores[i],
		}
		if n > -1 {
			doc.ID = documents[n].ID
			doc.DocumentID = documents[n].DocumentID
			doc.Metadata = documents[n].Metadata
		}
		//fmt.Println("rerank", doc.ID, doc.Score)
		if model == "act-embed-001" && doc.Score < 0.35 {
			continue
		}
		items = append(items, doc)
	}
	return items, nil
}

// for knowledge query
func (api Api) BuildQuery(query string, docs []types.Document, messages []model.Message, bot model.BotSetting) []types.ChatModelMessage {
	// 上下文
	temp := common.PromptTemplate(config.QUESTION_TEMPLATE)
	contexts, _ := temp.Render(struct {
		Contexts []types.Document
	}{
		Contexts: docs,
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
	var system = systemPrompt()
	msgs = append(msgs, types.ChatModelMessage{
		Role:    types.ChatMessageRoleSystem,
		Content: system,
	})
	msgs = append(msgs, types.ChatModelMessage{
		Role:    types.ChatMessageRoleUser,
		Content: str,
	})
	return msgs
}

func (api Api) Insert(document types.Document, update bool, bot model.BotSetting) error {

	var chunks map[string][]*types.Chunk
	var codes []string
	var err error

	setting, err := bot.MergeSetting(api.Config)
	if err != nil {
		return err
	}
	logx.Info(fmt.Sprintf("Chunk type %s", setting.Chunk.Type))
	switch setting.Chunk.Type {
	case "markdown":
		processor := splitters.MdPreprocessor(&splitters.PreprocessorConfig{
			MaxTokens: setting.Chunk.MaxTokens,
			MinTokens: setting.Chunk.MinTokens,
			Overlap:   setting.Chunk.Overlap,
		})
		chunks, codes, err = processor.Preprocess(document, bot)
	case "semantic":
		processor := splitters.SemanticPreprocessor(&splitters.PreprocessorConfig{
			MaxTokens: setting.Chunk.MaxTokens,
			MinTokens: setting.Chunk.MinTokens,
			Overlap:   setting.Chunk.Overlap,
		})
		chunks, codes, err = processor.Preprocess(document, bot)
	default:
		processor := splitters.TextPreprocessor(&splitters.PreprocessorConfig{
			MaxTokens: setting.Chunk.MaxTokens,
			MinTokens: setting.Chunk.MinTokens,
			Overlap:   setting.Chunk.Overlap,
		})
		chunks, codes, err = processor.Preprocess(document, bot)
	}

	if err != nil {
		logx.Error(fmt.Sprintf("Processor error: %s", err.Error()))
		return err
	}
	logx.Info(fmt.Sprintf("Chunks length: %d", len(chunks)))
	if update {
		err = model.DeleteSegments([]string{document.ID})
		if err != nil {
			return err
		}
	}

	num := 0
	i := 1
	for batch := range genBatches(chunks, EMBEDDINGS_BATCH_SIZE) {
		var list []string
		num += len(batch)
		for _, data := range batch {
			// replace in embedding
			list = append(list, strings.ReplaceAll(replaceCode(data.Text, []string{}), "\n", " "))
		}
		// 生成向量
		embeddings, err := api.Embedding(list, bot)
		if err != nil {
			return err
		}
		for m, embedding := range embeddings {
			data := batch[m]
			segument := &model.Segment{
				Id:          data.ID,
				KnowledgeId: data.DocumentID,
				Embedding:   embedding,
				Index:       i,
				// replace code into chunck
				Text:   replaceCode(data.Text, codes),
				Corpus: data.Metadata.Corpus,
				Source: data.Metadata.Source,
				Url:    data.Metadata.Url,
			}
			i++
			err = segument.Insert()
			if err != nil {
				logx.Error(fmt.Sprintf("Insert %s's segment failed: %s", document.ID, err.Error()))
			}
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
	embeddings, err := api.Embedding([]string{question}, bot)
	if err != nil {
		return []model.Segment{}, err
	}
	if len(embeddings) == 0 {
		return []model.Segment{}, fmt.Errorf("no embeddings returned for query")
	}
	num := bot.Retrieval
	if num == 0 {
		num = 20
	}
	return model.QueryEmbedding(embeddings[0], bot.Corpus, num, bot.Score/100)
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

func replaceCode(text string, codes []string) string {
	re := regexp.MustCompile(`\[code_block_(\d+)\]`)
	if len(codes) == 0 {
		return re.ReplaceAllString(text, "")
	}
	items := re.FindAllStringSubmatch(text, -1)
	if len(items) > 0 {
		for _, item := range items {
			reg := regexp.MustCompile(regexp.MustCompile("\\[").ReplaceAllString(item[0], "\\["))
			i, err := strconv.Atoi(item[1])
			if err != nil || i > len(codes)-1 {
				text = reg.ReplaceAllString(text, "")
			} else {
				code := codes[i]
				text = reg.ReplaceAllString(text, code)
			}
		}
	}
	return text
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
