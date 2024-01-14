package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/act-gpt/marino/config/system"
)

type TextReq struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Data   string `json:"data"`
}

func Html2text(html string) (string, error) {

	method := "POST"
	conf := system.Config.Parser
	reqUrl := conf.Host + conf.TextApi

	item := TextReq{
		Data: html,
	}
	body, err := json.Marshal(item)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, reqUrl, bytes.NewReader(body))

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var text TextReq
	err = json.Unmarshal(body, &text)
	if err != nil {
		return "", err
	}
	return text.Data, nil
}
