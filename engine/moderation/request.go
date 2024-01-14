package moderation

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"
)

func Request(content string, bot string, userId string) (int, error) {

	var conf = system.Config.Moderation
	// 未设置 url
	if !conf.CheckContent {
		return 200, nil
	}

	client := http.Client{}
	values := map[string]interface{}{
		"content": content,
	}

	jsondata, _ := json.Marshal(&values)
	req, _ := http.NewRequest(http.MethodPost, conf.Api, bytes.NewReader(jsondata))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("post err:%+v\n", err)
		return 0, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var checkResp types.WordCheckResponse
	json.Unmarshal([]byte(body), &checkResp)

	if checkResp.Code != "0" {
		return 0, err
	}

	if len(checkResp.WordList) > 0 {
		for _, item := range checkResp.WordList {
			if item.Level == "高" {
				category := item.Category
				if category == "政治" || category == "暴恐违禁" || category == "色情" || category == "不良价值观" {
					var val = 0
					var label = "Polity"
					switch category {
					case "政治":
						val = 1
						label = "Polity"
					case "暴恐违禁":
						val = 2
						label = "Terror"
					case "不良价值观":
						val = 3
						label = "Illegal"
					case "色情":
						val = 4
						label = "Porn"
					}
					blocked := model.Blocked{
						UserId:     userId,
						Content:    content,
						Reason:     label,
						BotId:      bot,
						Suggestion: "Block",
					}
					blocked.Insert()
					return val, nil
				}
			}
		}
	}
	return 200, nil
}
