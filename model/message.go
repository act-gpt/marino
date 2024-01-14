package model

import (
	"fmt"
	"time"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/types"

	"gorm.io/gorm"
)

type Message struct {
	Id               string         `json:"id" gorm:"primaryKey"`
	Source           string         `json:"source"`
	User             string         `json:"user"`
	ConversationId   string         `json:"conversation_id" gorm:"index"`
	BotId            string         `json:"bot_id" gorm:"index"`
	Status           string         `json:"status"`
	Question         string         `json:"question"`
	Answer           string         `json:"answer"`
	Model            string         `json:"model"`
	Ip               string         `json:"ip"`
	Like             int            `json:"like" gorm:"type:int;default:0"`
	Dislike          int            `json:"dislike" gorm:"type:int;default:0"`
	CostTime         float64        `json:"cost_time" gorm:"type:float"`
	LLMTime          float64        `json:"llm_time" gorm:"type:float"`
	LLMFirstTime     float64        `json:"llm_first_time" gorm:"type:float"`
	PromptTokens     int            `json:"prompt_tokens"`
	CompletionTokens int            `json:"completion_tokens"`
	TotalTokens      int            `json:"total_tokens"`
	Usage            *types.Usage   `json:"usage" gorm:"-:all"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (msg *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if msg.Id == "" {
		msg.Id = common.GetUUID()
	}
	return
}

func (msg *Message) Insert() error {
	err := DB.Create(msg).Error
	return err
}

func (msg *Message) Update() error {
	// This can update zero values
	err := DB.Model(msg).Updates(msg).Error
	return err
}

func (msg *Message) Delete() error {
	if msg.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	err := DB.Delete(msg).Error
	return err
}

func GetMessagesByConversation(id string, bot string, startIdx int, len int) (messages []Message, err error) {
	err = DB.Omit("updated_at", "deleted_at").Where("conversation_id = ? and bot_id = ? and status != ?", id, bot, "error").Order("created_at desc").Limit(len).Offset(startIdx).Find(&messages).Error
	return messages, err
}

func GetMessagesByBot(id string, startIdx int, len int, count bool) (messages []Message, total int64, err error) {
	if count {
		DB.Where("bot_id = ?", id).Select("id").Find(&messages).Count(&total)
	}
	err = DB.Where("bot_id = ?", id).Omit("updated_at", "deleted_at").Limit(len).Offset(startIdx).Order("created_at desc").Find(&messages).Error
	return messages, total, err
}

func LikeMessage(id string, like int) (int64, error) {
	result := DB.Model(Message{}).Where("id = ?", id).Update("like", like)
	return result.RowsAffected, result.Error
}
