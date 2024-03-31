package baidu

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/common/redis"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/engine/actgpt"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type BaiduResponse struct {
	Id               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	SentenceId       int    `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	IsTruncated      bool   `json:"is_truncated"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Result           string `json:"result"`
	ErrorCode        int    `json:"error_code"`
	ErrorMsg         string `json:"error_msg"`
	FinishReason     string `json:"finish_reason"`
	Usage            Usage  `json:"usage"`
}

type BaiduCompletionRequest struct {
	System      string                   `json:"system"`
	Temperature float32                  `json:"temperature"`
	Stream      bool                     `json:"stream"`
	User        string                   `json:"user_id"`
	Messages    []types.ChatModelMessage `json:"messages"`
}

type Authentication struct {
	Token   string `json:"access_token"`
	Expires int64  `json:"expires_in"`
}

const maxBufferSize = 512 * 1024
const MaxToken = 512
const Temperature = 0.01
const Model = "completions"
const Base = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/"
const AUTH_KEY = "BAIDU:AUTH"

type Client struct {
	actgpt.Client
}

func NewClient(apiKey string, bot model.BotSetting) *Client {
	return &Client{
		actgpt.Client{
			AccessKey:   apiKey,
			Model:       Model,
			MaxToken:    MaxToken,
			Temperature: Temperature,
			Bot:         bot,
		},
	}
}

// auth for baidu
func Auth() (string, error) {
	if system.Config.Redis.DataSource != "" {
		return redisAuth()
	}
	return dbAuth()
}

func (c *Client) createRequest(request BaiduCompletionRequest) (*http.Response, error) {
	c.buildConfig()
	reqUrl := Base
	if c.Model != "" {
		reqUrl = reqUrl + c.Model
	}

	if c.AccessKey != "" {
		reqUrl = reqUrl + "?access_token=" + c.AccessKey
	}

	body, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("error with Marshal json: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader(body))

	if err != nil {
		return nil, fmt.Errorf("error creating POST request: %v", err)
	}

	return client.Do(req)
}

func (c *Client) Completion(ctx context.Context, msgs []types.ChatModelMessage, cb func(types.ChatCompletionResponse)) error {

	system, messages := splitMessage(msgs)
	temperature := c.Temperature
	if temperature == 0 {
		temperature = 0.01
	}
	res, err := c.createRequest(BaiduCompletionRequest{
		System:      system,
		Messages:    messages,
		Temperature: temperature,
		User:        c.User,
	})

	if err != nil {
		return fmt.Errorf("POST completion: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed reading error response: %w", err)
		}
		return fmt.Errorf("%s", bodyBytes)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("failed reading response: %w", err)
	}
	var response BaiduResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("failed unmarshal response: %w", err)
	}

	if response.ErrorCode > 0 {
		return fmt.Errorf("response: %s", response.ErrorMsg)
	}

	if response.IsTruncated {
		return fmt.Errorf("$$content filtered")
	}

	usage := types.Usage(response.Usage)
	message := types.ChatModelMessage{
		Role:    types.ChatMessageRoleAssistant,
		Content: response.Result,
	}
	finishReason := "stop"
	if response.FinishReason == "" {
		if response.IsTruncated {
			finishReason = "length"
		}
	}

	choice := types.ChatCompletionChoice{
		Index:        0,
		Message:      message,
		FinishReason: finishReason,
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		cb(types.ChatCompletionResponse{
			ID:      response.Id,
			Object:  response.Object,
			Created: response.Created,
			Model:   c.Model,
			Choices: []types.ChatCompletionChoice{choice},
			Usage:   usage,
		})
		return nil
	}

}

func (c *Client) CompletionStream(ctx context.Context, msgs []types.ChatModelMessage, cb func(types.ChatCompletionStreamResponse)) error {

	system, messages := splitMessage(msgs)
	temperature := c.Temperature
	if temperature == 0 {
		temperature = 0.01
	}
	res, err := c.createRequest(BaiduCompletionRequest{
		System:      system,
		Messages:    messages,
		Temperature: temperature,
		User:        c.User,
		Stream:      true,
	})
	if err != nil {
		return fmt.Errorf("POST completion: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed reading error response: %w", err)
		}
		return fmt.Errorf("%s", bodyBytes)
	}

	isStream := strings.HasPrefix(common.Header(res.Header, "Content-Type"), "text/event-stream")
	if !isStream {
		var streamResponse BaiduResponse
		body, _ := io.ReadAll(res.Body)
		err = json.Unmarshal(body, &streamResponse)
		if err != nil {
			return fmt.Errorf("failed reading error response:  %w", err)
		}
		return fmt.Errorf("failed event-stream code: %d, response: %s", streamResponse.ErrorCode, streamResponse.ErrorMsg)
	}

	scanner := bufio.NewScanner(res.Body)
	// increase the buffer size to avoid running out of space
	buf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
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
			var response BaiduResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				if data == "[DONE]" {
					continue
				}
				return fmt.Errorf("error unmarshaling baidu llm parse:  %w", err)
			}

			end := response.IsEnd
			cause := response.FinishReason

			if response.IsTruncated {
				return fmt.Errorf("$$content filtered")
			}
			if cause == "normal" || cause == "stop" {
				cause = ""
			}
			if end {
				cause = "stop"
			}

			var choices []types.ChatCompletionStreamChoice

			choices = append(choices, types.ChatCompletionStreamChoice{
				Index: 0,
				Delta: types.ChatCompletionStreamChoiceDelta{
					Role:    types.ChatMessageRoleAssistant,
					Content: response.Result,
				},
				FinishReason: cause,
			})
			usage := types.Usage(response.Usage)
			cb(types.ChatCompletionStreamResponse{
				Created: response.Created,
				ID:      response.Id,
				Model:   c.Model,
				Object:  response.Object,
				Usage:   &usage,
				Choices: choices,
			})
		}
	}
	return nil
}

func (c *Client) buildConfig() {
	auth, err := Auth()
	if err != nil {
		fmt.Println(fmt.Errorf("\033[31;1;4mbaidu config error with: %v\033[0m", err))
	}
	c.AccessKey = auth
}

// turn message system into string for baidu api
func splitMessage(msgs []types.ChatModelMessage) (string, []types.ChatModelMessage) {
	system := ""
	var items []types.ChatModelMessage
	for _, msg := range msgs {
		if msg.Role == "system" {
			system = msg.Content
			continue
		}
		items = append(items, msg)
	}
	return system, items
}

// rquest accsee token for baidu
func authRequest(key string) (Authentication, error) {
	conf := system.Config.Baidu
	id := conf.ClientId
	secret := conf.ClientSecret
	// 未设置
	if id == "" || secret == "" {
		return Authentication{}, fmt.Errorf("ClientId and ClientSecret must be set before you use baidu model")
	}
	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials", id, secret)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(``))
	if err != nil {
		return Authentication{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return Authentication{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Authentication{}, err
	}

	var auth Authentication
	err = json.Unmarshal(body, &auth)
	if err != nil {
		logx.Error(fmt.Sprintf("Baidu auth error with: %s", url))
		return Authentication{}, err
	}
	now := time.Now()
	futureTime := now.Add(time.Second * time.Duration(auth.Expires))
	unixTime := futureTime.Unix()
	auth.Expires = unixTime - 1000*5
	return auth, nil
}

func redisAuth() (string, error) {
	val, _ := redis.Get(AUTH_KEY)
	now := time.Now()
	if val == "" {
		auth, err := authRequest(AUTH_KEY)
		if err != nil {
			return "", err
		}
		_body, err := json.Marshal(auth)
		if err == nil {
			redis.Set(AUTH_KEY, string(_body), 0)
		}
		return auth.Token, nil
	}
	var auth Authentication
	err := json.Unmarshal([]byte(val), &auth)
	if err != nil {
		return "", err
	}
	if now.Unix() > auth.Expires {
		auth, _ = authRequest(AUTH_KEY)
		_body, err := json.Marshal(auth)
		if err == nil {
			redis.Set(AUTH_KEY, string(_body), 0)
		}
	}
	return auth.Token, nil
}

func dbAuth() (string, error) {
	now := time.Now()
	conf, _ := model.LoadConfig(AUTH_KEY)
	/*
		if conf.Id == "" && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	*/
	// empty
	if conf.Id == "" {
		auth, err := authRequest(AUTH_KEY)
		if err != nil {
			logx.Error("Baidu access denied")
			return "", err
		}
		setting := common.Struct2JSON(auth)
		conf := model.Config{
			Type:    AUTH_KEY,
			Setting: setting,
		}
		conf.Insert()
		return auth.Token, nil
	}
	var auth Authentication
	str, _ := json.Marshal(conf.Setting)
	json.Unmarshal([]byte(str), &auth)
	// token expires
	if now.Unix()-5*100 > int64(auth.Expires) {
		auth, err := authRequest(AUTH_KEY)
		if err != nil {
			logx.Error("Baidu access denied")
			return "", err
		}
		// Update
		conf.Setting = common.Struct2JSON(auth)
		conf.Update()
		return auth.Token, nil
	}
	return auth.Token, nil
}
