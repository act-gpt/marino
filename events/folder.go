package events

import (
	"github.com/act-gpt/marino/model"
)

func DeleteFolder(folder *model.Folder) {
	// delete  knowledges
	knowledges, _ := model.DeleteKnowledgesByFolder(folder.Id)
	// destroy the embedding
	destroyKnownledges(folder.BotId, knowledges)
	folder.Delete()
	// sub folder
	folders, _ := model.GetSubFolders(folder.Id)
	for _, sub := range folders {
		DeleteFolder(sub)
	}
}
