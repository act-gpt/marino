package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/act-gpt/marino/config/system"
)

// 文档识别结构块
type Sugmentation struct {
	Status      bool   `json:"status"`
	Msg         string `json:"msg"`
	Filename    string `json:"filename"`     // 唯一 ID
	ContentType string `json:"content_type"` // 文本
	Data        string `json:"data"`         // 块类型
	Text        string `json:"text"`
}

func Document(filename string) (Sugmentation, error) {

	conf := system.Config.Parser
	reqUrl := conf.Host + conf.DocumentApi

	buff := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(buff)

	fileWriter, err := bodyWriter.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		fmt.Println("error writing to buffer")
		return Sugmentation{}, err
	}
	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return Sugmentation{}, err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return Sugmentation{}, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, reqUrl, buff)

	if err != nil {
		fmt.Println("error", err)
		return Sugmentation{}, err
	}
	req.Header.Add("Content-Type", contentType)
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("error", err)
		return Sugmentation{}, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error", err)
		return Sugmentation{}, err
	}

	var sugmentation Sugmentation
	err = json.Unmarshal(body, &sugmentation)
	if err != nil {
		fmt.Println("error", err)
		return Sugmentation{}, err
	}
	return sugmentation, nil
}
