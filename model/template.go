package model

import (
	"fmt"
	"time"

	"github.com/act-gpt/marino/common"

	"gorm.io/gorm"
)

type Template struct {
	Id          string         `json:"id" gorm:"primaryKey"`
	Avatar      string         `json:"avatar"`
	Name        string         `json:"name"  validate:"max=32"`
	Description string         `json:"description"`
	Kind        string         `json:"kind"`
	Language    string         `json:"language"`
	Setting     JSON           `json:"setting"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (template *Template) BeforeCreate(tx *gorm.DB) (err error) {
	if template.Id == "" {
		template.Id = common.GetUUID()
	}
	return
}

func (template *Template) Insert() error {
	err := DB.Create(template).Error
	return err
}

func (template *Template) Update() error {
	// This can update zero values
	err := DB.Model(template).Updates(template).Error
	return err
}

func (template *Template) Delete() error {
	if template.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	err := DB.Delete(template).Error
	return err
}

func GetTemplates(lang string) ([]*Template, error) {
	var templates []*Template
	err := DB.Order("created_at asc").Omit("updated_at", "deleted_at").Find(&templates, "language = ?", lang).Error
	return templates, err
}

func GetTemplate(id string) (*Template, error) {
	var template *Template
	err := DB.Omit("updated_at", "deleted_at").First(&template, "id = ?", id).Error
	return template, err
}
