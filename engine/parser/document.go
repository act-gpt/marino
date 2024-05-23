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
	"github.com/act-gpt/marino/types"
)

func Document(filename string) (types.Sugmentation, error) {

	conf := system.Config.Parser
	reqUrl := conf.Host + conf.DocumentApi

	buff := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(buff)
	var sugmentation types.Sugmentation
	fileWriter, err := bodyWriter.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		fmt.Println("error writing to buffer")
		return sugmentation, err
	}
	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return sugmentation, err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return sugmentation, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, reqUrl, buff)

	if err != nil {
		fmt.Println("error", err)
		return sugmentation, err
	}
	req.Header.Add("Content-Type", contentType)
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("error", err)
		return sugmentation, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error", err)
		return sugmentation, err
	}
	err = json.Unmarshal(body, &sugmentation)
	if err != nil {
		fmt.Println("error", err)
		return sugmentation, err
	}
	return sugmentation, nil
}
