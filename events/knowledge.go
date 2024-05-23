package events

import (
	"fmt"

	"github.com/act-gpt/marino/api"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/engine"
	"github.com/act-gpt/marino/engine/moderation"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"
)

func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

func update(knowledge *model.Knowledge, insert bool) {
	// for bot
	bot, err := model.GetSetting(knowledge.BotId, "")
	if err != nil {
		return
	}
	// parse html 2 text
	meta := &types.Metadata{
		Corpus: bot.Corpus,
	}

	text, err := engine.Html2Markdown(knowledge.Content)
	if err != nil {
		knowledge.Status = 0
		knowledge.Update()
		fmt.Println(err)
		return
	}

	document := types.Document{
		ID:       knowledge.Id,
		Text:     text,
		Metadata: *meta,
	}

	// content audit for Chinese
	if system.Config.Moderation.CheckContent {
		for _, item := range Chunks(text, 5000) {
			code, err := moderation.Request(item, knowledge.OrgId+":"+knowledge.UserId, knowledge.BotId)
			if err != nil || code != 200 {
				knowledge.Status = 3
				knowledge.Update()
				return
			}
		}
	}
	// embedding document
	err = api.Client.Insert(document, insert, bot)
	if err != nil {
		knowledge.Status = 0
		fmt.Println(err)
	} else {
		knowledge.Status = 1
	}
	knowledge.Update()
}

func CreateKnowledge(knowledge *model.Knowledge) {
	update(knowledge, false)
}

func UpdateKnowledge(knowledge *model.Knowledge) {
	update(knowledge, true)
}

func DeleteKnowledge(knowledge *model.Knowledge) {
	destroyKnownledges(knowledge.BotId, []*model.Knowledge{knowledge})
}
