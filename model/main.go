package model

import (
	"database/sql/driver"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config/system"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

//go:embed insert.sql
var sql string

type JSON map[string]interface{}

func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// 使用 field.Tag、field.TagSettings 获取字段的 tag
	// 查看 https://github.com/go-gorm/gorm/blob/master/schema/field.go 获取全部的选项
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

// Value Marshal
func (jsonField JSON) Value() (driver.Value, error) {
	return json.Marshal(jsonField)
}

// Scan Unmarshal
func (jsonField *JSON) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(data, &jsonField)
}

var DB *gorm.DB

func createRootAccountIfNeed() error {
	var user User
	var rootUser User
	if err := DB.First(&user).Error; err != nil {
		fmt.Println("\033[32;1;4mNo user exists, create a root user for you: username is admin, password is you@actgpt\033[0m")
		rootUser = User{
			Username:    "admin",
			Password:    "you@actgpt",
			Role:        common.RoleRootUser,
			Status:      common.UserStatusEnabled,
			DisplayName: "Admin",
			Information: map[string]interface{}{},
		}
		rootUser.Insert()
		DB.Exec(sql)
	}

	c := system.Config.Organization
	if c.Name != "" {
		var org Organization
		if err := DB.First(&org).Error; err != nil {
			Organization := Organization{
				Name:        c.Name,
				Contact:     c.Contact,
				Phone:       c.Phone,
				Owner:       rootUser.Id,
				Admin:       datatypes.JSON("[\"" + rootUser.Id + "\"]"),
				Information: map[string]interface{}{},
			}
			Organization.Insert()
			rootUser.OrgId = Organization.Id
			rootUser.Update(false)
		}
	}

	return nil
}

func CountTable(tableName string) (num int64) {
	DB.Table(tableName).Count(&num)
	return
}

func CheckConnection(source string) (bool, error) {
	_, err := gorm.Open(postgres.Open(source), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	return true, err
}

func InitDB(source string) (err error) {
	var db *gorm.DB
	if source == "" {
		return nil
	}
	var log logger.Interface
	if gin.Mode() == "debug" {
		log = logger.Default.LogMode(logger.Info)
	} else {
		log = logger.Default.LogMode(logger.Silent)
	}
	db, err = gorm.Open(postgres.Open(source), &gorm.Config{
		PrepareStmt: true,
		Logger:      log,
	})
	if err != nil {
		logx.Error(fmt.Sprintf("\033[31;1;4mDB access denied: %s\033[0m", err.Error()))
		return err
	}
	if err == nil {
		DB = db
		err := db.AutoMigrate(&User{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Organization{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Bot{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Folder{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Knowledge{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Template{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Message{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Blocked{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Config{})
		if err != nil {
			return err
		}
		InitSegments(system.Config)
		err = createRootAccountIfNeed()
		return err
	}
	return err
}

func CloseDB() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}
