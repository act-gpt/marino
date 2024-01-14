package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config/system"

	"gorm.io/gorm"
)

type Config struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type"`
	Setting   JSON      `json:"setting"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (config *Config) BeforeCreate(tx *gorm.DB) (err error) {
	if config.Id == "" {
		config.Id = common.GetUUID()
	}
	return err
}

func (config *Config) Insert() error {
	if DB == nil {
		return nil
	}
	err := DB.Create(config).Error
	return err
}

func (config *Config) Update() error {
	// This can update zero values
	err := DB.Model(config).Updates(config).Error
	return err
}

func LoadConfig(tp string) (Config, error) {
	var conf Config
	if DB == nil {
		return conf, nil
	}
	err := DB.Where("type = ?", tp).First(&conf).Error
	return conf, err
}

func InsertOrSaveSystemConfig(config system.SystemConfig) (system.SystemConfig, error) {
	if config.Initialled.Db {
		return SaveSystemConfig(config)
	} else {
		conf := Config{
			Type:    "stystem",
			Setting: common.Struct2JSON(config),
		}
		err := conf.Insert()
		return config, err
	}
}

func SaveSystemConfig(config system.SystemConfig) (system.SystemConfig, error) {
	conf, err := LoadConfig("stystem")
	if err != nil {
		return config, err
	}
	conf.Setting = common.Struct2JSON(config)
	err = conf.Update()
	return config, err
}

func LoadSystemConfig() (system.SystemConfig, error) {

	var conf Config
	var config = system.WithDefault()
	if DB == nil {
		return config, errors.New("db not init")
	}
	err := DB.Where("type = ?", "stystem").First(&conf).Error
	if err != nil {
		return config, err
	}
	m, err := json.Marshal(conf.Setting)
	if err != nil {
		return config, err
	}
	var c system.SystemConfig
	if err := json.Unmarshal(m, &c); err != nil {
		return config, err
	}
	return c, nil
}
