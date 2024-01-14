package model

import (
	"fmt"
	"time"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config/system"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Organization config
type OrganizationConfig struct {
	Type       string `json:"type"`
	Role       int    `json:"role"`
	Bots       int    `json:"bots"`
	Files      int    `json:"files"`
	Words      int    `json:"words"`
	Quota      int    `json:"quota"`
	Train      int    `json:"train"`
	Speech     int    `json:"speech"`
	Users      int    `json:"users"`
	Moderation bool   `json:"moderation"`
}

type Organization struct {
	Id          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;index"`
	Contact     string         `json:"contact"`
	Phone       string         `json:"phone"`
	Owner       string         `json:"owner"`
	AccessToken string         `json:"access_token" gorm:"uniqueIndex"` // this token is for bot management
	Admin       datatypes.JSON `json:"admin"`
	Information JSON           `json:"infomation"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	//Users       []*User        `gorm:"many2many:users_oganizations;"`
}

func (org *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if org.Id == "" {
		org.Id = common.GetUUID()
	}
	uuid := "OA." + common.GetUUID() + common.GetUUID()
	token, err := common.EncryptByAes([]byte(system.Config.Secret), []byte(uuid))
	org.AccessToken = token
	return err
}

func GetAllOrganization(startIdx int, num int) (orgs []*Organization, err error) {
	err = DB.Order("created_at desc").Omit("updated_at", "deleted_at", "access_token").Limit(num).Offset(startIdx).Find(&orgs).Error
	return orgs, err
}

func (org *Organization) Insert() error {
	err := DB.Create(org).Error
	return err
}

func (org *Organization) Update() error {
	err := DB.Model(org).Updates(org).Error
	return err
}

func (org *Organization) Delete() error {
	if org.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	err := DB.Model(org).Updates(org).Error
	return err
}

func GetOrganizationById(id string) (*Organization, error) {
	if id == "" {
		return nil, fmt.Errorf("id 为空！")
	}
	org := Organization{Id: id}
	var err error = nil
	err = DB.Omit("updated_at", "deleted_at", "access_token").First(&org).Error
	return &org, err
}

func GetOrganizations(user string) ([]*Organization, error) {
	var orgs []*Organization
	var err error
	if DB.Dialector.Name() == "postgres" {
		err = DB.Where(Organization{Owner: user}).Or("(admin)::jsonb ? '"+user+"'").Omit("updated_at", "deleted_at", "access_token").Find(&orgs).Error
	} else {
		err = DB.Where(Organization{Owner: user}).Or(datatypes.JSONArrayQuery("admin").Contains(user)).Omit("updated_at", "deleted_at", "access_token").Find(&orgs).Error
	}
	return orgs, err
}

func GetOrganization(id string, user string) (*Organization, error) {
	var org *Organization
	var err error
	if DB.Dialector.Name() == "postgres" {
		err = DB.Where(DB.Where(Organization{Id: id}).Where(DB.Where(Organization{Owner: user}).Or("(admin)::jsonb ? '"+user+"'"))).Omit("updated_at", "deleted_at", "access_token").First(&org).Error
	} else {
		err = DB.Where(DB.Where(Organization{Id: id}).Where(DB.Where(Organization{Owner: user}).Or(datatypes.JSONArrayQuery("admin").Contains(user)))).Omit("updated_at", "deleted_at", "access_token").First(&org).Error
	}

	return org, err
}
