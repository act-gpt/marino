package model

import (
	"time"
)

type Blocked struct {
	Id         int       `json:"id"`
	UserId     string    `json:"user_id"`
	BotId      string    `json:"bot_id"`
	Content    string    `json:"content" gorm:"type:text"`
	Reason     string    `json:"reason"`
	Suggestion string    `json:"suggestion"`
	CreatedAt  time.Time `json:"created_at"`
}

func GetAllBlocked(startIdx int, num int) ([]*Blocked, error) {
	var blocked []*Blocked
	err := DB.Order("created_at desc").Limit(num).Offset(startIdx).Find(&blocked).Error
	return blocked, err
}

func GetAllUserBlocked(userId int, startIdx int, num int) ([]*Blocked, error) {
	var blocked []*Blocked
	err := DB.Where("user_id = ?", userId).Order("created_at desc").Limit(num).Offset(startIdx).Find(&blocked).Error
	return blocked, err
}

func SearchBlocked(keyword string) (Blocked []*Blocked, err error) {
	err = DB.Where("reason = ?", keyword).Find(&Blocked).Error
	return Blocked, err
}

func (blocked *Blocked) Insert() error {
	err := DB.Create(blocked).Error
	return err
}
