package model

import (
	"fmt"
	"time"

	"github.com/act-gpt/marino/common"

	"gorm.io/gorm"
)

type Folder struct {
	Id        string         `json:"id" gorm:"primaryKey"`
	BotId     string         `json:"bot_id" gorm:"index"`
	OrgId     string         `json:"org_id" gorm:"index"`
	Name      string         `json:"name" gorm:"not null;" validate:"max=60"`
	Parent    string         `json:"parent"`
	Visable   bool           `json:"visable" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (folder *Folder) BeforeCreate(tx *gorm.DB) (err error) {
	if folder.Id == "" {
		folder.Id = common.GetUUID()
	}
	return
}

func (folder *Folder) Insert() error {
	var err error
	if folder.Id != "" {
		folder.Id = common.GetUUID()
	}
	err = DB.Create(folder).Error
	return err
}

func (folder *Folder) Update() error {
	err := DB.Model(folder).Updates(folder).Error
	return err
}

func (folder *Folder) Delete() error {
	err := DB.Delete(folder).Error
	return err
}

func GetFolderById(id string) (*Folder, error) {
	if id == "" {
		return nil, fmt.Errorf("id 为空！")
	}
	folder := Folder{Id: id}
	var err error = nil
	err = DB.Omit("updated_at", "deleted_at").First(&folder, "id = ?", id).Error
	return &folder, err
}

func GetFoldersByBotId(id string, org string) ([]*Folder, error) {
	if id == "" {
		return nil, fmt.Errorf("id 为空！")
	}
	var folders []*Folder
	var err error = nil
	err = DB.Where("bot_id = ?", id).Omit("updated_at", "deleted_at").Find(&folders, "org_id = ?", org).Error
	return folders, err
}

func GetSubFolders(id string) ([]*Folder, error) {
	if id == "" {
		return nil, fmt.Errorf("id 为空！")
	}
	var folders []*Folder
	var err error = nil
	err = DB.Where("parent = ?", id).Omit("updated_at", "deleted_at").Find(&folders).Error
	return folders, err
}

func DeleteFoldersByBot(bot string) error {
	return DB.Where("bot_id = ?", bot).Delete(&Folder{}).Error
}
