package model

import (
	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config/system"

	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BotSetting struct {
	// Language
	Lang string `json:"lang"`
	// Name
	Name string `json:"name"`
	//Type int    `json:"type"`
	// Engine model
	Model  string `json:"model"`
	Avatar string `json:"avatar"`
	// Custome prompt
	Prompt string `json:"prompt"`
	// Welcome message
	Welcome     string  `json:"welcome"`
	Description string  `json:"description"`
	Corpus      string  `json:"corpus"`
	Link        string  `json:"link"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
	// verctor score
	Score float64 `json:"score"`
	// contexts in LLM prompt
	Contexts int `json:"contexts"`
	// histories in LLM prompt
	Histories int `json:"histories"`
}

type Bot struct {
	Id       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" validate:"max=32"`
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

func (bot *Bot) BeforeCreate(tx *gorm.DB) (err error) {
	if bot.Id == "" {
		bot.Id = common.GetUUID()
	}
	uuid := common.GetUUID() + common.GetUUID()
	token, err := common.EncryptByAes([]byte(system.Config.Secret), []byte(uuid))
	bot.AccessToken = "BA." + token
	return err
}

func (bot *Bot) Insert() error {
	err := DB.Create(bot).Error
	return err
}

func (bot *Bot) Update() error {
	// This can update zero values
	err := DB.Model(bot).Select("name", "enabled", "approved", "status", "owner", "config", "setting", "admin").Updates(bot).Error
	return err
}

func (bot *Bot) Delete() error {
	if bot.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	err := DB.Delete(bot).Error
	return err
}

func BotAccessToken(token string) (*Bot, error) {
	var bot *Bot
	err := DB.Where("enabled = true").First(&bot, "access_token = ?", token).Error
	return bot, err
}

func GetBotById(id string) (*Bot, error) {
	if id == "" {
		return nil, fmt.Errorf("id 为空！")
	}
	bot := Bot{Id: id}
	var err error = nil
	err = DB.Omit("updated_at", "deleted_at").Where("enabled = true").First(&bot, "id = ?", id).Error
	return &bot, err
}

func EnabledBotById(id string, enabled bool) error {
	return DB.Exec("UPDATE bots SET enabled=$1 WHERE id = $2", enabled, id).Error
}

func GetAllBotByOrgId(orgId string, startIdx int, num int) ([]*Bot, error) {
	var bots []*Bot
	err := DB.Omit("updated_at", "deleted_at").Where("org_id = ?", orgId).Order("created_at desc").Limit(num).Offset(startIdx).Find(&bots).Error
	return bots, err
}

func CountAllBotByOrgId(orgId string) (count int64, err error) {
	var bot *Bot
	err = DB.Model(&bot).Where("org_id = ?", orgId).Count(&count).Error
	return count, err
}

func GetBotByAdminAndId(id string, admin string) (*Bot, error) {
	if id == "" {
		return nil, fmt.Errorf("id 为空！")
	}
	var bot *Bot
	var err error
	if DB.Dialector.Name() == "postgres" {
		err = DB.Where("(admin)::jsonb ? '"+admin+"'").Where("id = ?", id).Omit("updated_at", "deleted_at").First(&bot).Error
	} else {
		err = DB.Where("id = ?", id).Where(datatypes.JSONArrayQuery("admin").Contains(admin)).Omit("updated_at", "deleted_at").First(&bot).Error
	}
	return bot, err
}

func GetAllBotByOwner(owner string) ([]*Bot, error) {
	var bots []*Bot
	err := DB.Where("owner = ?", owner).Omit("updated_at", "deleted_at").Order("created_at desc").Find(&bots).Error
	return bots, err
}

func GetAllBotByAdmin(admin string) ([]*Bot, error) {
	var bots []*Bot
	var err error
	if DB.Dialector.Name() == "postgres" {
		err = DB.Where("(admin)::jsonb ? '"+admin+"'").Omit("updated_at", "deleted_at").Order("created_at desc").Find(&bots).Error
	} else {
		err = DB.Where(datatypes.JSONArrayQuery("admin").Contains(admin)).Omit("updated_at", "deleted_at").Order("created_at desc").Find(&bots).Error
	}
	return bots, err
}

func GetSetting(id string, admin string) (BotSetting, error) {
	var bot Bot
	var err error

	if admin != "" {
		if DB.Dialector.Name() == "postgres" {
			err = DB.Where("(admin)::jsonb ? '"+admin+"'").First(&bot, "id = ?", id).Error
		} else {
			err = DB.Where(datatypes.JSONArrayQuery("admin").Contains(admin)).First(&bot, "id = ?", id).Error
		}
	} else {
		err = DB.First(&bot, "id = ?", id).Error
	}

	var setting BotSetting
	text, e := json.Marshal(bot.Setting)
	if e != nil {
		return setting, e
	}
	json.Unmarshal([]byte(text), &setting)
	return setting, err
}

func UpdateSetting(id string, admin string, settings BotSetting) (BotSetting, error) {
	var bot Bot
	var err error
	var setting BotSetting

	if admin != "" {
		if DB.Dialector.Name() == "postgres" {
			err = DB.Where("(admin)::jsonb ? '"+admin+"'").First(&bot, "id = ?", id).Error
		} else {
			err = DB.Where(datatypes.JSONArrayQuery("admin").Contains(admin)).First(&bot, "id = ?", id).Error
		}
	} else {
		err = DB.First(&bot, "id = ?", id).Error
	}
	if err != nil {
		return setting, err
	}
	data := common.Struct2JSON(settings)
	bot.Name = data["name"].(string)
	bot.Setting = data
	bot.Update()
	return setting, err
}
