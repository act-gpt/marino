package common

import (
	"fmt"
	"log"
	"strings"

	"github.com/act-gpt/marino/types"

	"github.com/pkoukk/tiktoken-go"
)

var Encoder = newDummyTokenizer()

type dummyTokenizer struct {
	//encoder *tokenizer.Encoder
	encoder *tiktoken.Tiktoken
}

func newDummyTokenizer() *dummyTokenizer {
	//encoder, err := tokenizer.NewEncoder()
	// tiktoken
	encoder, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
	if err != nil {
		// We assume that there's no error.
		panic(err)
	}
	return &dummyTokenizer{encoder: encoder}
}

// Encode iterates through runes and returns a slice of the leading runes, which
// consume at most tokenNum number of tokens.
func (t *dummyTokenizer) Encode(items []rune, tokenNum int) ([]rune, error) {
	runes := items
	if len(runes) > tokenNum {
		leng := tokenNum + 180
		if leng > len(runes) {
			leng = len(runes)
		}
		runes = runes[:leng]
	}
	for i := len(runes) - 1; i >= 0; i-- {
		rs := runes[:i]
		tokens := t.encoder.Encode(string(rs), nil, nil)
		if len(tokens) <= tokenNum {
			return rs, nil
		}
	}
	return runes, nil
}

func (t *dummyTokenizer) Length(str string) int {
	tokens := t.encoder.Encode(str, nil, nil)
	return len(tokens)
}

func NumTokensFromMessages(messages []types.ChatModelMessage, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return
	}
	var tokensPerMessage, tokensPerName int
	switch model {
	case "gpt-3.5-turbo-0613",
		"gpt-3.5-turbo-16k-0613",
		"gpt-4-0314",
		"gpt-4-32k-0314",
		"gpt-4-0613",
		"gpt-4-32k-0613":
		tokensPerMessage = 3
		tokensPerName = 1
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
		tokensPerName = -1   // if there's a name, the role is omitted
	default:
		if strings.Contains(model, "gpt-3.5-turbo") {
			return NumTokensFromMessages(messages, "gpt-3.5-turbo-0613")
		} else if strings.Contains(model, "gpt-4") {
			return NumTokensFromMessages(messages, "gpt-4-0613")
		} else {
			err = fmt.Errorf("num_tokens_from_messages() is not implemented for model %s. See https://github.com/openai/openai-python/blob/main/chatml.md for information on how messages are converted to tokens.", model)
			log.Println(err)
			return
		}
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>
	return numTokens
}
