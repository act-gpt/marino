package model

import (
	"fmt"
	"time"

	"github.com/act-gpt/marino/common"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Knowledge struct {
	Id       string `json:"id" gorm:"primaryKey"`
	FolderId string `json:"folder_id" gorm:";index"`
	BotId    string `json:"bot_id" gorm:";index"`
	OrgId    string `json:"org_id" gorm:";index"`
	UserId   string `json:"user_id" gorm:";index"`
	Name     string `json:"name" gorm:"not null;"`
	Content  string `json:"content" gorm:"type:text"`
	Ext      string `json:"ext"`
	Path     string `json:"path" gorm:"index"`
	// 0 error, 1 success, 2 processing, 3 security locked, 4 over quota
	Status     int            `json:"status" gorm:"type:int;default:1"`
	Sha        string         `json:"sha"`
	UpdateUser string         `json:"update_user"`
	Tags       datatypes.JSON `json:"tags"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (knowledge *Knowledge) BeforeCreate(tx *gorm.DB) (err error) {
	if knowledge.Id == "" {
		knowledge.Id = common.GetUUID()
	}
	return
}

func (knowledge *Knowledge) Insert() error {
	var err error
	if knowledge.Id != "" {
		knowledge.Id = common.GetUUID()
	}
	err = DB.Create(knowledge).Error
	return err
}

func (knowledge *Knowledge) Update() error {
	err := DB.Model(knowledge).Updates(knowledge).Error
	return err
}

func (knowledge *Knowledge) Delete() error {
	if knowledge.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	err := DB.Delete(knowledge).Error
	return err
}

func SearchKnowledges(keyword string, botId string, orgId string) (knowledges []*Knowledge, err error) {
	err = DB.Where("name LIKE ? AND bot_id = ? AND org_id = ?", "%"+keyword+"%", botId, orgId).Select("id", "folder_id", "bot_id", "org_id", "name").Find(&knowledges).Error
	return knowledges, err
}

func GetKnowledgeById(id string) (*Knowledge, error) {
	if id == "" {
		return nil, fmt.Errorf("id 为空！")
	}
	knowledge := Knowledge{Id: id}
	var err error = nil
	err = DB.Omit("updated_at", "deleted_at").First(&knowledge).Error
	return &knowledge, err
}

func GetKnowledgesByBotId(bot string, startIdx int, len int) ([]*Knowledge, int64, error) {
	if bot == "" {
		return nil, 0, fmt.Errorf("bot_id 为空！")
	}
	var knowledge []*Knowledge
	var err error = nil
	var total int64
	DB.Where("bot_id = ?", bot).Select("id").Find(&knowledge).Count(&total)
	err = DB.Omit("updated_at", "deleted_at").Where("bot_id = ?", bot).Limit(len).Offset(startIdx).Find(&knowledge).Error
	return knowledge, total, err
}

func GetKnowledgesByPath(bot string, path string, startIdx int, len int) ([]*Knowledge, int64, error) {
	if path == "" {
		return nil, 0, fmt.Errorf("path 为空！")
	}
	var knowledge []*Knowledge
	var err error = nil
	var total int64
	DB.Where("path LIKE ?", path+"%").Select("id").Find(&knowledge).Count(&total)
	err = DB.Omit("updated_at", "deleted_at").Where("path LIKE ?", path+"%").Limit(len).Offset(startIdx).Find(&knowledge, "bot_id = ?", bot).Error
	return knowledge, total, err
}

func DeleteKnowledgesByBot(bot string) ([]*Knowledge, error) {
	var knowledges []*Knowledge
	var err error = nil
	err = DB.Omit("updated_at", "deleted_at").Where("bot_id = ?", bot).Find(&knowledges).Error
	if err != nil {
		return nil, err
	}
	DB.Where("bot_id = ?", bot).Delete(&Knowledge{})
	return knowledges, err
}

func DeleteKnowledgesByFolder(folder string) ([]*Knowledge, error) {
	var knowledges []*Knowledge
	var err error = nil
	err = DB.Omit("updated_at", "deleted_at").Where("folder = ?", folder).Find(&knowledges).Error
	if err != nil {
		return nil, err
	}
	DB.Where("folder_id = ?", folder).Delete(&Knowledge{})
	return knowledges, err
}
