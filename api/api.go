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
	setting, err := bot.MergeSetting(api.Config)
	if err != nil {
		return nil, err
	}
	res, err := embedding.Request(text, setting.Chunk.Embedding)
	if err != nil {
		return nil, err
	}
	var embeddings []types.Embedding
	for _, data := range res.Data {
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

func (api Api) Reranker(query string, document []types.Document, bot model.BotSetting) ([]types.Document, error) {
	var docs []string
	for _, doc := range document {
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
	var items []types.Document
	data := res.Data

	for i, item := range data.Documents {
		n := indexOf(document, item)
		// 小于 0.35 质量已经很低了，过滤掉
		doc := types.Document{
			Text:  item,
			Score: data.Scores[i],
		}
		if n > -1 {
			doc.ID = document[n].ID
			doc.DocumentID = document[n].DocumentID
			doc.Metadata = document[n].Metadata
		}
		fmt.Println("rerank", doc.ID, doc.Score)
		if doc.Score < 0.35 {
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
	switch setting.Chunk.Type {
	case "markdown":
		processor := splitters.MdPreprocessor(&splitters.PreprocessorConfig{})
		chunks, codes, err = processor.Preprocess(document, bot)
	default:
		processor := splitters.SemanticPreprocessor(&splitters.PreprocessorConfig{})
		chunks, codes, err = processor.Preprocess(document, bot)
	}
	if err != nil {
		return err
	}

	if update {
		err = model.DeleteSegments([]string{document.ID})
		if err != nil {
			return err
		}
	}

	num := 0
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
		for i, embedding := range embeddings {
			data := batch[i]
			segument := &model.Segment{
				Id:          data.ID,
				KnowledgeId: data.DocumentID,
				Embedding:   embedding,
				Index:       i + 1,
				// replace code into chunck
				Text:   replaceCode(data.Text, codes),
				Corpus: data.Metadata.Corpus,
				Source: data.Metadata.Source,
				Url:    data.Metadata.Url,
			}
			err = segument.Insert()
			if err != nil {
				fmt.Println(err)
				return err
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
