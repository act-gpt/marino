package model

import (
	"fmt"
	"time"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config/system"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Setting struct {
	Id       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"" validate:"max=32"`
	OrgId    string `json:"org_id" gorm:"index"`
	Enabled  bool   `json:"enabled" gorm:"type:bool;default:true"`
	Approved int    `json:"approved" gorm:"type:int;default:1"`
	// 0: failed, 1: success, 2: queue
	Status      int            `json:"status" gorm:"type:int;default:1"`
	Owner       string         `json:"owner" gorm:"index"`
	AccessToken string         `json:"access_token,omitempty" gorm:"uniqueIndex"` // this token is for bot management
	UpdateUser  string         `json:"update_user"`
	Setting     JSON           `json:"setting"`
	Admin       datatypes.JSON `json:"admin,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (bot *Setting) BeforeCreate(tx *gorm.DB) (err error) {
	if bot.Id == "" {
		bot.Id = common.GetUUID()
	}
	uuid := "BA." + common.GetUUID() + common.GetUUID()
	token, err := common.EncryptByAes([]byte(system.Config.Secret), []byte(uuid))
	bot.AccessToken = token
	return err
}

func (bot *Setting) Insert() error {
	err := DB.Create(bot).Error
	return err
}

func (bot *Setting) Update() error {
	// This can update zero values
	err := DB.Model(bot).Select("name", "enabled", "access_token", "approved", "status", "owner", "config", "setting", "admin").Updates(bot).Error
	return err
}

func (bot *Setting) Delete() error {
	if bot.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	err := DB.Delete(bot).Error
	return err
}
