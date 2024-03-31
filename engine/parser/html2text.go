package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader(body))

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var text TextReq
	err = json.Unmarshal(body, &text)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return text.Data, nil
}
