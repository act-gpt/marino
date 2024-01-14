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
	Length   int    `json:"length"`
}

var MODELS = map[string]interface{}{
	// ERNIE-Bot
	// 0.012, 0.012
	"completions": &MODEL{
		Name:     "completions",
		Owner:    "baidu",
		Disabled: false,
		Length:   1024 * 4,
	},
	// ERNIE-Bot-8K
	// 0.024元/千tokens, 0.048元/千tokens
	"ernie_bot_8k": &MODEL{
		Name:     "ernie_bot_8k",
		Owner:    "baidu",
		Disabled: false,
		Length:   1024 * 8,
	},
	// ERNIE-Bot 4.0,
	// 0.12 k/tokens, 0.12 k/tokens
	"completions_pro": &MODEL{
		Name:     "completions_pro",
		Owner:    "baidu",
		Disabled: false,
		Length:   1024 * 4,
	},
	// act-gpt-001
	// 0.004 k/tokens, 0.004 k/tokens
	"act-gpt-001": &MODEL{
		Name:     "act-gpt-001",
		Owner:    "actgpt",
		Disabled: false,
		Length:   1024 * 8,
	},
	// act-gpt-002
	// 0.006 k/tokens, 0.006 k/tokens
	"act-gpt-002": &MODEL{
		Name:     "act-gpt-002",
		Owner:    "actgpt",
		Disabled: false,
		Length:   1024 * 32,
	},
	// act-gpt-003
	// 0.012 k/tokens, 0.012 k/tokens
	"act-gpt-003": &MODEL{
		Name:     "act-gpt-003",
		Owner:    "actgpt",
		Disabled: false,
		Length:   1024 * 32,
	},
	"gpt-3.5-turbo": &MODEL{
		Name:     "gpt-3.5-turbo",
		Owner:    "openai",
		Disabled: false,
		Length:   1024 * 16,
	},
	"gpt-4-turbo": &MODEL{
		Name:     "gpt-4-1106-preview",
		Owner:    "openai",
		Disabled: true,
		Length:   1024 * 128,
	},
}

func GetAvailableModel() *MODEL {
	config := system.Config.Initialled

	if config.ActGpt {
		return MODELS["act-gpt-001"].(*MODEL)
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
	item := MODELS[model]
	if item == nil {
		return nil
	}
	val := item.(*MODEL)
	val.Name = m
	return val
}
