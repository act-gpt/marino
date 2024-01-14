package events

import (
	"github.com/act-gpt/marino/api"
	"github.com/act-gpt/marino/model"
)

func destroyKnownledges(id string, knowledges []*model.Knowledge) {

	var ids = []string{}
	//for
	for _, knowledge := range knowledges {
		ids = append(ids, knowledge.Id)
	}
	err := api.Client.Delete(ids)
	if err != nil {
		return
	}
}
