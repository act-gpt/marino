package engine

import (
	"context"
	"fmt"

	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/engine/actgpt"
	"github.com/act-gpt/marino/engine/baidu"
	"github.com/act-gpt/marino/engine/embedding"
	"github.com/act-gpt/marino/engine/moderation"
	"github.com/act-gpt/marino/engine/openai"
	"github.com/act-gpt/marino/engine/parser"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"

	copenai "github.com/sashabaranov/go-openai"
)

type LLM interface {
	SetModel(string)
	SetMaxToken(int)
	SetTemperature(float32)
	SetUser(string)
	Completion(context.Context, []types.ChatModelMessage, func(types.ChatCompletionResponse)) error
	CompletionStream(context.Context, []types.ChatModelMessage, func(types.ChatCompletionStreamResponse)) error
}

func New(bot model.BotSetting) (LLM, error) {
	conf := config.GetModel(bot)
	config := system.Config
	if conf == nil {
		return nil, fmt.Errorf("model does not exist: %s", bot.Model)
	}
	name := conf.Name
	owner := conf.Owner
	switch owner {
	case "openai":
		return openai.NewClient(config.OpenAi.AccessKey, bot), nil
	case "baidu":
		return baidu.NewClient("", bot), nil
	case "actgpt":
		return actgpt.NewClient(config.ActGpt.AccessKey, bot), nil
	default:
		return nil, fmt.Errorf("model does not support: %s", name)
	}
}

func Embedding(input []string) (copenai.EmbeddingResponse, error) {
	return embedding.Request(input)
}

func Moderation(content string, bot string, userId string) (int, error) {
	return moderation.Request(content, bot, userId)
}

func Document(filename string) (parser.Sugmentation, error) {
	return parser.Document(filename)
}

func Html2text(html string) (string, error) {
	return parser.Html2text(html)
}
