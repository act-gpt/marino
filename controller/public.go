package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/act-gpt/marino/api"
	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/engine/moderation"
	"github.com/act-gpt/marino/events"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"

	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"gorm.io/datatypes"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/exp/slices"
)

type QueryRequest struct {
	Conversation string  `json:"conversation"`
	Prompt       string  `json:"prompt"`
	Stream       bool    `json:"stream"`
	Temperature  float32 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
	User         string  `json:"user"`
}

type SegemntResponse struct {
	Id    string  `json:"id"`
	Doc   string  `json:"doc"`
	Score float64 `json:"score"`
	Text  string  `json:"text"`
	Sha   string  `json:"sha"`
}

type ResultResponse struct {
	types.ChatCompletionResponse
	Source []SegemntResponse `json:"source,omitempty"`
}

type StreamResultResponse struct {
	types.ChatCompletionStreamResponse
	Source []SegemntResponse `json:"source,omitempty"`
}

type ErrorResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

var FREE = "FREE"

func responseError(c *gin.Context, res ErrorResponse) {
	c.JSON(res.StatusCode, res)
}

func Redirect(c *gin.Context) {
	url := c.Query("target")
	logx.Info("Redirect: " + url)
	c.Redirect(http.StatusMovedPermanently, url)
}

func Sign(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	nonce := common.GetUUID()
	timestamp := now.UnixMilli()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nonce":     nonce,
		"timestamp": timestamp,
		"bot":       id,
		"exp":       now.Add(time.Hour * 4).Unix(),
	})
	token, _ := at.SignedString([]byte(system.Config.Secret))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data": map[string]interface{}{
			"nonce":     nonce,
			"sign":      token,
			"timestamp": timestamp,
		},
	})
}

func ConversationHistories(c *gin.Context) {
	bot := c.Param("id")
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(c.Query("size"))
	if size <= 0 {
		size = common.ItemsPerPage
	}
	conversation := c.Query("cvs")
	data, err := model.GetMessagesByConversation(conversation, bot, (page-1)*size, size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    data,
	})
}

// get segment detail by id
func SegmentDetail(c *gin.Context) {

	bot, err := model.GetBotById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	// not found
	if bot.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    403,
			"message": "NotFound",
		})
		return
	}

	result, err := api.Client.Get(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    result,
	})
}

func Like(c *gin.Context) {

	type Like struct {
		Like int `json:"like"`
	}
	req := Like{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	num, err := model.LikeMessage(c.Param("id"), req.Like)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    num,
	})
}

func ValidateQuery(c *gin.Context) (*model.Bot, *model.Organization, QueryRequest, error) {

	req := QueryRequest{}
	err := c.ShouldBindJSON(&req)
	switch {
	case errors.Is(err, io.EOF):
		responseError(c, ErrorResponse{
			Code:       400,
			Message:    "missing request body",
			StatusCode: 400,
		})
		return &model.Bot{}, &model.Organization{}, req, err
	case err != nil:
		responseError(c, ErrorResponse{
			Code:       400,
			Message:    err.Error(),
			StatusCode: 400,
		})
		return &model.Bot{}, &model.Organization{}, req, err
	}
	if len(req.Prompt) > 2000 {
		responseError(c, ErrorResponse{
			Code:       400,
			Message:    "too many content",
			StatusCode: 400,
		})
		return &model.Bot{}, &model.Organization{}, req, err
	}

	_bot, ok := c.Get("bot")
	if !ok {
		responseError(c, ErrorResponse{
			Code:       404,
			Message:    "not found",
			StatusCode: 404,
		})
		return &model.Bot{}, &model.Organization{}, req, err
	}

	bot := _bot.(*model.Bot)
	_org, ok := c.Get("org")
	if !ok {
		responseError(c, ErrorResponse{
			Code:       404,
			Message:    "not found",
			StatusCode: 404,
		})
		return &model.Bot{}, &model.Organization{}, req, err
	}
	org := _org.(*model.Organization)

	if system.Config.Moderation.CheckContent {
		code, err := moderation.Request(req.Prompt, req.User, bot.Id)
		if err != nil || code != 200 {
			var block = config.BlockedMessages[code]
			responseError(c, ErrorResponse{
				Code:       7700 + code,
				Message:    block,
				StatusCode: 404,
			})
			events.Emmiter.Emit("green", bot.Id, code)
			if err == nil {
				err = errors.New("moderation failed")
			}
			return bot, org, req, err
		}
	}
	return bot, org, req, err
}

// Query bot
func Query(c *gin.Context) {

	bot, _, req, err := ValidateQuery(c)
	if err != nil {
		return
	}

	id := c.Param("id")
	from := "web"
	start := time.Now()

	isStream := strings.HasPrefix(common.Header(c.Request.Header, "Accept"), "text/event-stream")
	if req.Stream {
		isStream = true
	}

	if c.Query("source") == "" {
		match, _ := regexp.MatchString("Dev.", req.Conversation)
		if match {
			from = "dev"
		}
	}
	if strings.Contains(c.Request.URL.Path, "/v1/chat") {
		from = "api"
	}
	var setting model.BotSetting
	s, _ := json.Marshal(bot.Setting)
	err = json.Unmarshal(s, &setting)

	if err != nil {
		logx.Error(fmt.Errorf("failed to parse json: %w", err))
	}

	// request llm
	link := setting.Link
	mod := setting.Model
	if link != "" {
		mod = link
	}
	// 数据统计
	// messsage
	msg := model.Message{
		Id:             common.GetUUID(),
		ConversationId: req.Conversation,
		User:           req.User,
		Model:          mod,
		Question:       req.Prompt,
		BotId:          id,
		Source:         from,
		Status:         "start",
		Ip:             c.ClientIP(),
		CostTime:       0,
		LLMTime:        0,
		LLMFirstTime:   0,
		Answer:         "",
	}

	var source []SegemntResponse
	var messages []types.ChatModelMessage
	// normal for baidu, finished for act-gpt
	STOPS := []string{"stop", "length", "function_call", "tool_calls", "normal", "finished"}

	if setting.Contexts == 0 {
		msgs, err := model.GetMessagesByConversation(req.Conversation, id, 0, setting.Histories)
		messages = api.Client.BuildConversion(msgs)
		if err != nil {
			responseError(c, ErrorResponse{
				Code:       500,
				Message:    err.Error(),
				StatusCode: 500,
			})
			return
		}
	} else {
		msgs, err := model.GetMessagesByConversation(req.Conversation, id, 0, setting.Histories)
		if err != nil {
			responseError(c, ErrorResponse{
				Code:       500,
				Message:    err.Error(),
				StatusCode: 500,
			})
			return
		}
		segments, err := api.Client.Query(req.Prompt, setting)
		if err != nil {
			responseError(c, ErrorResponse{
				Code:       500,
				Message:    err.Error(),
				StatusCode: 500,
			})
			return
		}
		messages = api.Client.BuildQuery(req.Prompt, segments, msgs, setting)
		for _, val := range segments {
			source = append(source, SegemntResponse{
				Id:    val.Id,
				Doc:   val.KnowledgeId,
				Score: val.Score,
				Text:  val.Text,
				Sha:   val.Sha,
			})
		}
	}

	defer func() {
		// insert message
		msg.CostTime = time.Since(start).Seconds()
		if msg.Usage != nil {
			msg.TotalTokens = msg.Usage.TotalTokens
			msg.PromptTokens = msg.Usage.PromptTokens
			msg.CompletionTokens = msg.Usage.CompletionTokens
		}
		msg.Insert()
		events.Emmiter.Emit("query", msg)
	}()

	temperature := req.Temperature
	tokens := req.MaxTokens
	if temperature == 0 {
		temperature = setting.Temperature
	}
	if tokens == 0 {
		tokens = 512
	}
	llm, err := api.Client.Engine(setting)
	llm.SetModel(mod)
	llm.SetTemperature(temperature)
	llm.SetMaxToken(tokens)
	llm.SetUser(req.User)

	if err != nil {
		responseError(c, ErrorResponse{
			Code:       500,
			Message:    err.Error(),
			StatusCode: 500,
		})
		return
	}

	ch := make(chan any)
	llmStart := time.Now()
	go func() {
		defer close(ch)
		start := false
		if !isStream {
			if err := llm.Completion(c.Request.Context(), messages, func(res types.ChatCompletionResponse) {
				if !start {
					start = true
					msg.LLMFirstTime = time.Since(llmStart).Seconds()
				}
				ch <- res
			}); err != nil {
				msg.Status = "error"
				ch <- err
			}
		} else {
			if err := llm.CompletionStream(c.Request.Context(), messages, func(res types.ChatCompletionStreamResponse) {
				if !start {
					start = true
					msg.LLMFirstTime = time.Since(llmStart).Seconds()
				}
				ch <- res
			}); err != nil {
				msg.Status = "error"
				ch <- err
			}
		}
		msg.LLMTime = time.Since(llmStart).Seconds()
	}()
	// not stream
	if !isStream {
		for resp := range ch {

			if res, ok := resp.(types.ChatCompletionResponse); ok {

				item := ResultResponse{
					ChatCompletionResponse: res,
					Source:                 source,
				}
				c.JSON(http.StatusOK, item)
				msg.Usage = &item.Usage
				reason := ""
				for _, val := range item.Choices {
					if slices.Contains(STOPS, val.FinishReason) {
						reason = val.FinishReason
						msg.Answer = val.Message.Content
						break
					}
				}
				if reason == "normal" {
					reason = "stop"
				}
				msg.Status = reason

			} else {
				msg.Status = "error"
				err := resp.(error)
				code := 500
				msg := err.Error()
				// Just for chinese content filter
				if strings.HasPrefix(msg, "$$") {
					code = 7704
					msg = config.BlockedMessages[code]
				}
				responseError(c, ErrorResponse{
					Code:       code,
					Message:    msg,
					StatusCode: 500,
				})
				return
			}
		}
		return
	}

	c.Status(200)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	// Sets the proxy unbuffered streaming
	c.Header("X-Accel-Buffering", "no")
	c.Writer.Flush()

	num := 0
	answer := ""
	c.Stream(func(w io.Writer) bool {
		val, ok := <-ch
		if !ok {
			return false
		}
		// ChatCompletionStreamResponse
		if item, ok := val.(types.ChatCompletionStreamResponse); ok {
			num += 1
			reason := ""
			for _, val := range item.Choices {
				answer += val.Delta.Content
				if slices.Contains(STOPS, val.FinishReason) {
					reason = val.FinishReason
					break
				}
			}

			if reason == "normal" || reason == "finished" {
				reason = "stop"
			}

			msg.Status = reason
			if reason != "" {
				msg.Usage = item.Usage
				item.Usage = nil
				msg.Answer = answer
				c.SSEvent("", StreamResultResponse{
					ChatCompletionStreamResponse: item,
					Source:                       source,
				})
				c.SSEvent("", "[DONE]")
				return true
			}
			item.Usage = nil
			c.SSEvent("", StreamResultResponse{
				ChatCompletionStreamResponse: item,
			})
			return true
		}

		// something error
		item := val.(error)
		logx.WithCallerSkip(1).Error(item.Error())
		code := 500
		msg.Status = "error"
		err := item.Error()
		// Just for chinese content filter
		if strings.HasPrefix(err, "$$") {
			code = 7004
			err = config.BlockedMessages[4]
		}
		c.SSEvent("", ErrorResponse{
			Code:    code,
			Message: err,
		})
		return false
	})
}
