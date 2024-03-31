package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"

	"github.com/antlabs/strsim"
	"github.com/go-aie/xslices"
)

type Preprocessor struct {
	cfg *common.PreprocessorConfig
}

func NewPreprocessor(cfg *common.PreprocessorConfig) *Preprocessor {
	return &Preprocessor{
		cfg: cfg.Init(),
	}
}

func (p *Preprocessor) Preprocess(docs ...*types.Document) (map[string][]*types.Chunk, error) {
	chunkMap := make(map[string][]*types.Chunk)

	for _, doc := range docs {
		docID := doc.ID
		meta := doc.Metadata
		if docID == "" {
			docID = common.GetUUID()
		}

		var llm = "ernie-speed-128k"
		var length = 1024 * 128
		var size = length/2 - common.Encoder.Length(config.SPLIT_PROMPT) - 50

		fmt.Println("AI Spilt process: ", docID, size)
		textChunks, err := p.split(doc.Text, llm, size)
		if err != nil {
			return nil, err
		}
		//fmt.Println("Chuncks size: ", len(textChunks))
		for _, textChunk := range textChunks {
			id := common.GetUUID()
			// return chunks
			chunkMap[docID] = append(chunkMap[docID], &types.Chunk{
				ID:         id,
				DocumentID: docID,
				Text:       textChunk,
				Metadata:   meta,
			})
		}
	}

	return chunkMap, nil
}

func remian(data []common.SplitData) int {
	s := ""
	for _, item := range data {
		txt := item.Segment
		if txt == "" {
			txt = item.Question + "\n" + item.Answer
		}
		s += txt
	}
	return len(s)
}

func MergeData(list []common.SplitData, max int, min int) []string {
	var cs []string
	var str []string
	var length = len(list)
	for j := 0; j < length; {
		item := list[j]
		txt := item.Segment
		if txt == "" {
			txt = item.Question + "\n" + item.Answer
		}

		// append
		str = append(str, txt)
		// calculate length
		leng := common.TokensLength(str)
		// if end
		if j == length-1 {
			val := strings.Join(str[:], "\n\n")
			cs = append(cs, val)
			break
		}
		j++
		// contine if less then min
		if leng < min {
			continue
		}

		// if remain content is less then max
		if remian(list[j:]) < max-min {
			continue
		}
		// if is MAX
		if leng > max {
			// Only one element in list
			if len(str) == 1 {
				str = str[:]
			} else {
				str = str[:len(str)-1]
				j--
			}
		}
		val := strings.Join(str[:], "\n\n")
		cs = append(cs, val)
		str = []string{}
	}
	return cs
}

// spilt text to chunks
func (p *Preprocessor) split(text string, llm string, size int) ([]string, error) {
	if text == "" || strings.TrimSpace(text) == "" {
		return nil, nil
	}

	var cs []string
	// get chuncks
	var s = chunks(text, size)
	for i := 0; i < len(s); i++ {
		doc := s[i]
		same := 0.90
		// max request
		var MAX_REQUEST = 3
		var list []common.SplitData
		for i := 0; i < MAX_REQUEST; i++ {
			// get split data
			item, err := p.request(doc, llm, size, false)
			if err != nil {
				fmt.Println(err)
				continue
			}
			dist := compare(doc, merge(item))
			// if we got it
			if dist >= same {
				list = item
				fmt.Printf("AI preprocessor: %d times, levenshtein Distance: %f\n", i, dist)
				break
			}
			// get best result
			if dist > compare(doc, merge(list)) {
				fmt.Printf("Levenshtein Distance: %f\n", dist)
				list = item
			}
		}
		// per chunk max length
		//var MAX = 600
		// per chunk min length
		//var MIN = 400
		items := MergeData(list, p.cfg.MaxChunkNum, p.cfg.MinChunkCharNum)
		cs = append(cs, items...)
	}
	return cs, nil
}

// request to llm
func (p *Preprocessor) request(doc string, llm string, length int, jsonPrompt bool) ([]common.SplitData, error) {
	engine, err := New(model.BotSetting{
		Model: "completions_pro",
	})
	if err != nil {
		return nil, err
	}
	engine.SetModel(llm)
	engine.SetTemperature(0)
	engine.SetMaxToken(length)

	var prompts []types.ChatModelMessage

	prompts = append(prompts, types.ChatModelMessage{
		Role:    types.ChatMessageRoleSystem,
		Content: "You are an AI assistant that helps people to finish task.",
	})

	temp := common.PromptTemplate(config.SPLIT_PROMPT)
	if jsonPrompt {
		temp = common.PromptTemplate(config.SPLIT_PROMPT_JSON)
	}
	str, _ := temp.Render(struct {
		Document string
	}{
		Document: doc,
	})
	prompts = append(prompts, types.ChatModelMessage{
		Role:    types.ChatMessageRoleUser,
		Content: str,
	})

	ch := make(chan any)
	go func() {
		defer close(ch)
		if err := engine.Completion(context.Background(), prompts, func(res types.ChatCompletionResponse) {
			ch <- res
		}); err != nil {
			ch <- err
		}
	}()

	text := ""

	for resp := range ch {
		if res, ok := resp.(types.ChatCompletionResponse); ok {
			for _, choice := range res.Choices {
				text += choice.Message.Content
			}
			if res.Choices[0].FinishReason != "stop" {
				fmt.Printf("Not finished yet, reason: %s", res.Choices[0].FinishReason)
			}
		} else {
			err := resp.(error)
			return nil, err
		}
	}

	if jsonPrompt {
		var transactions []common.SplitData
		text = strings.TrimSpace(text)
		m1 := regexp.MustCompile(`\t`)
		text = m1.ReplaceAllString(text, "")
		b := []byte(text)
		e := json.Unmarshal(b, &transactions)
		if e != nil {
			return nil, e
		}
		return transactions, nil
	}
	transactions := common.FormatSplitText(text)
	return transactions, nil

}

// merge data to string
func merge(list []common.SplitData) string {
	text := ""
	for _, item := range list {
		segment := item.Segment
		if segment == "" {
			segment = item.Question + "\n" + item.Answer + "\n"
		}
		text += segment
	}
	return text
}

// calculate string similarity
func compare(text string, text1 string) float64 {
	val := strsim.Compare(text, text1)
	return val
}

func chunks(s string, chunkSize int) []string {
	runes := []rune(s)
	if len(runes) == 0 {
		return nil
	}
	var chunks []string
	var i int
	//overlap := 200
	for i < len(runes) {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes) + 1
		}
		chunkRunes, err := common.Encoder.Encode(runes[i:end], chunkSize)
		if err != nil {
			chunks = append(chunks, string(runes[i:end]))
			continue
		}
		text := string(chunkRunes)
		var lastPuncIdx = -1
		for _, punc := range []rune{
			'.', '?', '!', '\n',
			'。', '？', '！', '；',
		} {

			lastPuncIdx = xslices.Max(lastPuncIdx, lastRuneIndex(text, punc))
		}
		if lastPuncIdx != -1 {
			last := lastPuncIdx + 1
			text = string([]rune(text)[:last])
			i += last
			if len(runes)-i < 20 {
				text = string(runes[i-last:])
				i = len(runes)
			}
		} else {
			i += utf8.RuneCountInString(text)
		}
		chunks = append(chunks, text)
	}
	return chunks
}

func lastRuneIndex(s string, r rune) int {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == r {
			return i
		}
	}
	return -1
}
