package model

import (
	"fmt"
	"time"

	"github.com/act-gpt/marino/common"

	"gorm.io/gorm"
)

type Conversation struct {
	Id        string         `json:"id" gorm:"primaryKey"`
	User      string         `json:"user"`
	BotId     string         `json:"bot_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (conversation *Conversation) BeforeCreate(tx *gorm.DB) (err error) {
	if conversation.Id == "" {
		conversation.Id = common.GetUUID()
	}
	return
}

func (conversation *Conversation) Insert() error {
	err := DB.Create(conversation).Error
	return err
}

func (conversation *Conversation) Update() error {
	// This can update zero values
	err := DB.Model(conversation).Updates(conversation).Error
	return err
}

func (conversation *Conversation) Delete() error {
	if conversation.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	err := DB.Delete(conversation).Error
	return err
}
