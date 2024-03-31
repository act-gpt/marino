package actgpt

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"
)

const maxBufferSize = 512 * 1024
const MaxToken = 512
const Temperature = 0.01
const Model = "act-gpt-001"

type Client struct {
	AccessKey   string           `json:"access_key"`
	Model       string           `json:"model"`
	MaxToken    int              `json:"max_token"`
	Temperature float32          `json:"temperature"`
	User        string           `json:"user"`
	Bot         model.BotSetting `json:"bot"`
}

func NewClient(apiKey string, bot model.BotSetting) *Client {
	return &Client{
		AccessKey:   apiKey,
		Model:       Model,
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

func (c *Client) createRequest(request types.ChatCompletionRequest) (*http.Response, error) {
	conf := system.Config.ActGpt
	host := conf.Host
	key := conf.AccessKey
	if c.AccessKey != "" {
		key = c.AccessKey
	}
	url := strings.TrimRight(host, "/") + "/v1" + "/chat/completions"
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error with Marshal json: %v", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("error creating POST request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	if key != "" {
		req.Header.Add("Authorization", "Bearer "+key)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST completion: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		if bodyBytes, err := io.ReadAll(res.Body); err != nil {
			return nil, fmt.Errorf("failed reading error response: %w", err)
		} else {
			return nil, fmt.Errorf("%s", bodyBytes)
		}
	}
	return res, nil
}

func (c *Client) Completion(ctx context.Context, msgs []types.ChatModelMessage, cb func(types.ChatCompletionResponse)) error {

	model := config.GetModel(c.Bot)
	messages := getMessages(msgs, c.MaxToken, model.Length)

	request := types.ChatCompletionRequest{
		Model:       c.Model,
		Messages:    messages,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		User:        c.User,
	}

	res, err := c.createRequest(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed reading response: %w", err)
	}
	var response types.ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed unmarshal response: %w", err)
	}
	//fmt.Println(response.Choices[0])
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		cb(response)
		return nil
	}
}

func (c *Client) CompletionStream(ctx context.Context, msgs []types.ChatModelMessage, cb func(types.ChatCompletionStreamResponse)) error {

	model := config.GetModel(c.Bot)
	messages := getMessages(msgs, c.MaxToken, model.Length)

	request := types.ChatCompletionRequest{
		Model:       c.Model,
		Messages:    messages,
		MaxTokens:   c.MaxToken,
		Temperature: c.Temperature,
		User:        c.User,
		Stream:      true,
	}

	res, err := c.createRequest(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	isStream := strings.HasPrefix(common.Header(res.Header, "Content-Type"), "text/event-stream")
	if !isStream {
		var streamResponse types.ChatCompletionResponse
		body, _ := io.ReadAll(res.Body)
		err = json.Unmarshal(body, &streamResponse)
		if err != nil {
			return fmt.Errorf("failed reading error response:  %w", err)
		}
		return fmt.Errorf("failed event-stream code: %d, response: %s", res.StatusCode, string(body))
	}
	scanner := bufio.NewScanner(res.Body)
	// increase the buffer size to avoid running out of space
	buf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := strings.Index(string(data), "\r\n\r\n"); i >= 0 {

			return i + 4, data[0:i], nil
		}
		if i := strings.Index(string(data), "\n\n"); i >= 0 {

			return i + 2, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			// This handles the request cancellation
			return ctx.Err()
		default:
			data := scanner.Text()
			if len(data) < 6 {
				continue
			}
			data = data[6:]
			var response types.ChatCompletionStreamResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				if data == "[DONE]" {
					continue
				}
				return fmt.Errorf("error unmarshaling actgpt llm parse:  %w", err)
			}
			cb(response)
		}
	}
	return nil
}

func getMessages(msgs []types.ChatModelMessage, max int, total int) []types.ChatModelMessage {
	// TODO: Count msgs tokens and reduce it
	return msgs
}
