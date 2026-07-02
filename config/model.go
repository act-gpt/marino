package config

import (
	"strings"

	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/model"
)

type MODEL struct {
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Disabled bool   `json:"disabled"`
	Show     bool   `json:"show"`
	Length   int    `json:"length"`
}

var MODELS = map[string]interface{}{
	// act-gpt-001
	"act-gpt-001": &MODEL{
		Name:     "deepseek-chat",
		Owner:    "actgpt",
		Disabled: false,
		Show:     true,
		Length:   1024 * 64,
	},
	// act-gpt-002
	"act-gpt-002": &MODEL{
		Name:     "qwen-chat",
		Owner:    "actgpt",
		Disabled: false,
		Show:     true,
		Length:   1024 * 32,
	},
	"gpt-3.5-turbo": &MODEL{
		Name:     "gpt-3.5-turbo",
		Owner:    "openai",
		Disabled: false,
		Show:     false,
		Length:   1024 * 16,
	},
	"gpt-4-turbo": &MODEL{
		Name:     "gpt-4-turbo",
		Owner:    "openai",
		Disabled: false,
		Show:     false,
		Length:   1024 * 128,
	},
}

func GetAvailableModel() *MODEL {
	config := system.Config.Initialled
	if config.ActGpt {
		return MODELS["act-gpt-002"].(*MODEL)
	}
	if config.Baidu {
		return MODELS["completions"].(*MODEL)
	}
	return MODELS["gpt-3.5-turbo"].(*MODEL)
}

func GetModel(bot model.BotSetting) *MODEL {
	link := bot.Link
	model := bot.Model
	if link != "" {
		model = link
	}
	m := model
	// for azure and open compatibly
	if strings.HasPrefix(model, "gpt") {
		if strings.HasPrefix(model, "gpt-4") {
			model = "gpt-4-turbo"
		}
		if strings.HasPrefix(model, "gpt-3.5") {
			model = "gpt-3.5-turbo"
		}
	}
	if strings.HasPrefix(model, "deepseek") {
		model = "act-gpt-001"
	}
	if strings.HasPrefix(model, "qwen") {
		model = "act-gpt-002"
	}
	item := MODELS[model]
	if item == nil {
		return nil
	}
	val := item.(*MODEL)
	val.Name = m
	return val
}
