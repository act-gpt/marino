package events

import (
	"github.com/act-gpt/marino/model"
)

func CreateBot(bot model.Bot) {
	/*
		collection, err := redis.GetCollection()
		if err != nil {
			collection = "bot_1"
		}
		bot.Config["collectionName"] = collection
		bot.Config["corpus"] = bot.OrgId + ":" + bot.Id
		bot.Update()
	*/
}

func DeleteBot(bot *model.Bot) {
	knowledges, _ := model.DeleteKnowledgesByBot(bot.Id)
	model.DeleteFoldersByBot(bot.Id)
	var ids = []string{}
	for _, knowledge := range knowledges {
		ids = append(ids, knowledge.Id)
	}
	///model.DeleteSegmentByKnowledges(ids)
	destroyKnownledges(bot.Id, knowledges)
}
