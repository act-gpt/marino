package controller

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/act-gpt/marino/api"
	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/engine"
	"github.com/act-gpt/marino/engine/baidu"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"

	"dario.cat/mergo"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckDbInput struct {
	DataSource string `json:"DataSource"`
}

type CheckEmbeddingInput struct {
	Host      string `json:"Host"`
	Api       string `json:"Api"`
	Model     string `json:"Model"`
	AccessKey string `json:"AccessKey"`
}

type CheckBaiduInput struct {
	ClientId     string `json:"ClientId"`
	ClientSecret string `json:"ClientSecret"`
}

type CheckActGptInput struct {
	AccessKey string `json:"AccessKey"`
}

type CheckOpenaiInput struct {
	Type       string `json:"Type"`
	AccessKey  string `json:"AccessKey"`
	Host       string `json:"Host"`
	APIVersion string `json:"APIVersion"`
}

func CheckDb(c *gin.Context) {
	req := CheckDbInput{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	_, err = model.CheckConnection(req.DataSource)
	code := 0
	msg := "Ok"
	if err != nil {
		code = 500
		msg = "Connected to db failed"
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
	})
}

func CheckEmbedding(c *gin.Context) {
	req := CheckEmbeddingInput{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	old := system.Config.Embedding
	conf := &system.Config.Embedding
	conf.AccessKey = req.AccessKey
	conf.Host = req.Host
	conf.Model = req.Model
	conf.Api = req.Api
	res, err := engine.Embedding([]string{"hello"})
	code := 0
	msg := "Ok"
	if err != nil {
		code = 500
		msg = err.Error()
	}
	if len(res.Data) == 0 {
		code = 500
		msg = "Something error"
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
	})
	conf.AccessKey = old.AccessKey
	conf.Host = old.Host
	conf.Model = old.Model
	conf.Api = old.Api
}

func CheckResponse(c *gin.Context, llm engine.LLM, cb func(ok bool)) {
	ch := make(chan any)
	var messages []types.ChatModelMessage
	messages = append(messages, types.ChatModelMessage{
		Role:    types.ChatMessageRoleUser,
		Content: "hello",
	})
	go func() {
		defer close(ch)
		if err := llm.Completion(c.Request.Context(), messages, func(res types.ChatCompletionResponse) {
			ch <- res
		}); err != nil {
			ch <- err
		}
	}()
	for resp := range ch {
		if _, ok := resp.(types.ChatCompletionResponse); ok {
			c.JSON(http.StatusOK, ErrorResponse{
				Code:    0,
				Message: "OK",
			})
			cb(true)
		} else {
			logx.WithCallerSkip(2).Error(resp.(error).Error())
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: resp.(error).Error(),
			})
			cb(false)
		}
	}

}

func CheckBaidu(c *gin.Context) {
	req := CheckBaiduInput{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	conf := &system.Config.Baidu
	conf.ClientId = req.ClientId
	conf.ClientSecret = req.ClientSecret
	llm, err := api.Client.Engine(model.BotSetting{
		Model: "completions",
	})
	if err != nil {
		responseError(c, ErrorResponse{
			Code:       500,
			Message:    err.Error(),
			StatusCode: 500,
		})
		return
	}
	llm.SetMaxToken(1)
	CheckResponse(c, llm, func(ok bool) {
		conf.ClientId = ""
		conf.ClientSecret = ""
		conf, err := model.LoadConfig(baidu.AUTH_KEY)
		if err != nil || conf.Id == "" {
			return
		}
		conf.Setting = model.JSON{
			"access_token": "",
			"expires_in":   0,
		}
		conf.Update()
	})
}

func CheckActGpt(c *gin.Context) {
	req := CheckActGptInput{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	llm, err := api.Client.Engine(model.BotSetting{
		Model: "act-gpt-001",
	})
	if err != nil {
		responseError(c, ErrorResponse{
			Code:       500,
			Message:    err.Error(),
			StatusCode: 500,
		})
		return
	}
	conf := &system.Config.ActGpt
	if os.Getenv("DEMO") != "1" {
		conf.AccessKey = req.AccessKey
	}
	llm.SetMaxToken(1)
	CheckResponse(c, llm, func(ok bool) {
		conf.AccessKey = ""
	})
}

func CheckOpenAI(c *gin.Context) {
	req := CheckOpenaiInput{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	llm, err := api.Client.Engine(model.BotSetting{
		Model: "gpt-3.5-turbo-1106",
	})
	if err != nil {
		responseError(c, ErrorResponse{
			Code:       500,
			Message:    err.Error(),
			StatusCode: 500,
		})
		return
	}
	conf := &system.Config.OpenAi
	conf.Host = req.Host
	conf.AccessKey = req.AccessKey
	conf.Type = req.Type
	conf.APIVersion = req.APIVersion
	llm.SetMaxToken(1)

	CheckResponse(c, llm, func(ok bool) {
		conf.Host = "https://api.openai.com"
		conf.Type = "openai"
		conf.AccessKey = ""
		conf.APIVersion = "2023-05-15"
	})
}

func CheckEngine(c *gin.Context) {
	t := c.Query("type")
	if t == "baidu" {
		CheckBaidu(c)
		return
	} else if t == "actgpt" {
		CheckActGpt(c)
		return
	} else if t == "openai" {
		CheckOpenAI(c)
		return
	} else if t == "db" {
		CheckDb(c)
		return
	} else if t == "embedding" {
		CheckEmbedding(c)
		return
	}
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Bad request",
	})
}

func ShaIntegrity(c *gin.Context) {
	body := []byte(common.EbbedFile())
	url := c.Query("url")
	if url != "" {
		resp, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"code":    500,
				"error":   err.Error(),
			})
			return
		}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"code":    500,
				"error":   err.Error(),
			})
			return
		}
	}
	h := sha512.New384()
	h.Write(body)
	bs := h.Sum(nil)
	encodedStr := base64.StdEncoding.EncodeToString(bs)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data": gin.H{
			"integrity": "sha384-" + encodedStr,
		},
	})
}

func Models(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": config.MODELS,
	})
}

func GetSystemConfig(c *gin.Context) {
	conf, err := model.LoadSystemConfig()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "forbidden",
		})
		return
	}

	if c.GetInt("role") != common.RoleGuestUser {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"Initialled": conf.Initialled,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": conf,
	})
}

func SetSystemConfig(c *gin.Context) {
	config := system.SystemConfig{}
	role := c.GetInt("role")
	conf := system.Config
	inited := conf.Initialled.Db
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	// already configured and user is not root
	if inited && role != common.RoleRootUser {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "forbidden",
		})
		return
	}
	if err := mergo.Merge(&config, conf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	// set system config
	system.InitNeedSave(config)
	// Db is not connected
	if model.DB == nil {
		fmt.Println("init db")
		if err := model.InitDB(config.Db.DataSource); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}
	}
	common.WriteEnv(map[string]string{
		"DB_DATA_SOURCE": config.Db.DataSource,
	})

	if os.Getenv("DEMO") == "1" {
		config.ActGpt.AccessKey = ""
	}

	// reset to actgpt
	if config.Db.Dimension == 768 {
		config.Embedding.Host = config.ActGpt.Host
		config.Embedding.AccessKey = config.ActGpt.AccessKey
		config.Embedding.Model = "act-gpt-001"
		config.Embedding.Api = "/v1/embeddings"
	}

	// Update or Insert
	if !inited {
		conf := model.Config{
			Type:    "stystem",
			Setting: common.Struct2JSON(config),
		}
		conf.Insert()
	} else {
		if _, err := model.SaveSystemConfig(config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": config,
	})
}

type IPInfo struct {
	CountryCode string `json:"country_code"`
}

func GoDoc(c *gin.Context) {
	ip := c.ClientIP()
	redirect := "https://doc.act-gpt.com"
	cn := "https://doc.act-gpt.cn"
	url := "https://api.ip.sb/geoip/" + ip
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.Redirect(303, redirect)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		c.Redirect(303, redirect)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.Redirect(303, redirect)
		return
	}
	var info IPInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		c.Redirect(303, redirect)
		return
	}
	if info.CountryCode == "CN" {
		c.Redirect(303, cn)
		return
	}
	c.Redirect(303, redirect)
}
