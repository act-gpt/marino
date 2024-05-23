package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"

	openai "github.com/sashabaranov/go-openai"
)

const MaxToken = 512
const Temperature = 0.01
const APIVersion = "2023-05-15"

type Client struct {
	AccessKey   string           `json:"access_key"`
	Model       string           `json:"model"`
	MaxToken    int              `json:"max_token"`
	Temperature float32          `json:"temperature"`
	User        string           `json:"user"`
	Bot         model.BotSetting `json:"bot"`
	APIVersion  string           `json:"api_version"`
}

func NewClient(apiKey string, bot model.BotSetting) *Client {
	return &Client{
		AccessKey:   apiKey,
		Model:       bot.Model,
		MaxToken:    MaxToken,
		Temperature: Temperature,
		Bot:         bot,
	}
}

func (c *Client) SetModel(model string) {
	c.Model = model
}

// SetMaxToken 设置最大token数
func (c *Client) SetMaxToken(maxToken int) {
	c.MaxToken = maxToken
}

// SetTemperature 设置响应灵活程度
func (c *Client) SetTemperature(temperature float32) {
	c.Temperature = temperature
}

func (c *Client) SetUser(user string) {
	c.User = user
}

func (c *Client) Completion(ctx context.Context, msgs []types.ChatModelMessage, cb func(types.ChatCompletionResponse)) error {

	conf := c.buildConfig()
	cli := openai.NewClientWithConfig(conf)

	m := config.GetModel(c.Bot)
	messages := getMessages(msgs, c.MaxToken, m.Length)

	request := openai.ChatCompletionRequest{
		Model:       c.Model,
		Messages:    messages,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		User:        c.User,
	}

	res, err := cli.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s error with request: %v", c.Model, err)
	}

	var choices []types.ChatCompletionChoice
	for _, data := range res.Choices {
		reason := data.FinishReason
		if reason == "" {
			reason = "stop"
		}
		choices = append(choices, types.ChatCompletionChoice{
			Index:        data.Index,
			Message:      types.ChatModelMessage(data.Message),
			FinishReason: reason,
		})
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		cb(types.ChatCompletionResponse{
			ID:      res.ID,
			Object:  res.Object,
			Created: res.Created,
			Model:   c.Model,
			Choices: choices,
			Usage:   types.Usage(res.Usage),
		})
		return nil
	}
}

func (c *Client) CompletionStream(ctx context.Context, msgs []types.ChatModelMessage, cb func(types.ChatCompletionStreamResponse)) error {

	conf := c.buildConfig()
	cli := openai.NewClientWithConfig(conf)

	m := config.GetModel(c.Bot)
	messages := getMessages(msgs, c.MaxToken, m.Length)
	promptTokens := common.NumTokensFromMessages(msgs, c.Model)
	reply := ""
	finished := ""

	request := openai.ChatCompletionRequest{
		Model:       c.Model,
		Messages:    messages,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		User:        c.User,
		Stream:      true,
	}
	stream, err := cli.CreateChatCompletionStream(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s error with request: %v", c.Model, err)
	}
	defer stream.Close()

	for {
		select {
		case <-ctx.Done():
			// This handles the request cancellation
			return ctx.Err()
		default:
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return err
			}
			var choices []types.ChatCompletionStreamChoice
			for _, data := range response.Choices {
				choices = append(choices, types.ChatCompletionStreamChoice{
					Index:        data.Index,
					Delta:        types.ChatCompletionStreamChoiceDelta(data.Delta),
					FinishReason: data.FinishReason,
				})
				reply += data.Delta.Content
				if finished == "" {
					finished = data.FinishReason
				}
			}
			var usage types.Usage
			if finished != "" {
				completionTokens := common.NumTokensFromMessages([]types.ChatModelMessage{
					{
						Role:    types.ChatMessageRoleAssistant,
						Content: reply,
					},
				}, c.Model)
				usage = types.Usage{
					PromptTokens:     promptTokens,
					CompletionTokens: completionTokens,

					TotalTokens: promptTokens + completionTokens,
				}
			}
			res := types.ChatCompletionStreamResponse{
				Created: response.Created,
				ID:      response.ID,
				Model:   c.Model,
				Object:  response.Object,
				Choices: choices,
				Usage:   &usage,
			}
			cb(res)
		}
	}
}

func (c *Client) buildConfig() openai.ClientConfig {
	conf := &system.Config.OpenAi
	config := openai.DefaultConfig(conf.AccessKey)
	// for other OpenAI compatible api
	if conf.Type == "openai" {
		// trim last slash
		config.BaseURL = strings.TrimRight(conf.Host, "/") + "/v1"
	}
	// change to azure
	if conf.Type == "azure" {
		config = openai.DefaultAzureConfig(conf.AccessKey, conf.Host)
		config.AzureModelMapperFunc = func(model string) string {
			azureModelMapping := map[string]string{
				"gpt-3.5-turbo":        "gpt-35-turbo",
				"gpt-4":                "gpt-4",
				"gpt-4o":               "gpt-4o",
				"gpt-4-turbo":          "gpt-4-turbo",
				"gpt-4-vision-preview": "gpt-4-vision",
				"gpt-4-visio":          "gpt-4-vision",
			}
			return azureModelMapping[model]
		}
		if conf.APIVersion != "" {
			config.APIVersion = conf.APIVersion
		}
	}
	fmt.Println("config", config)
	return config
}

func getMessages(msgs []types.ChatModelMessage, max int, total int) []openai.ChatCompletionMessage {
	// TODO: Count msgs tokens and reduce it
	var messages []openai.ChatCompletionMessage
	for _, message := range msgs {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return messages
}
