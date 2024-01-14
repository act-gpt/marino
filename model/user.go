package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/act-gpt/marino/common"

	"gorm.io/gorm"
)

type User struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Username    string `json:"username" gorm:"unique;index" validate:"max=16"`
	Password    string `json:"password" gorm:"not null;" validate:"min=8,max=20"`
	DisplayName string `json:"display_name" gorm:"index" validate:"max=20"`
	Role        int    `json:"role" gorm:"type:int;default:1"`
	// enabled, disabled, pending, banned
	Status           int            `json:"status" gorm:"type:int;default:1"`
	Email            string         `json:"email" gorm:"unique;index" validate:"max=50"`
	Phone            string         `json:"phone" gorm:"index;type:varchar(20)"`
	WeChatId         string         `json:"wechat_id" gorm:"column:wechat_id;index"`
	VerificationCode string         `json:"verification_code" gorm:"-:all"`  // this field is only for Email verification, don't save it to database!
	AccessToken      string         `json:"access_token" gorm:"uniqueIndex"` // this token is for system management
	OrgId            string         `json:"org_id" gorm:"index"`
	Information      JSON           `json:"information"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	//Organizations    []*Organization `gorm:"many2many:users_oganizations;"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if strings.TrimSpace(user.Password) == "" {
		return fmt.Errorf("empty password")
	}
	user.Password, err = common.Password2Hash(user.Password)
	if err != nil {
		return err
	}
	if user.Id == "" {
		user.Id = common.GetUUID()
	}
	user.AccessToken = "AG." + common.GetUUID() + common.GetUUID()
	return
}

func GetAllUsers(startIdx int, num int) (users []*User, err error) {
	err = DB.Order("created_at desc").Limit(num).Offset(startIdx).Omit("password").Find(&users).Error
	return users, err
}

func SearchUsers(keyword string) (users []*User, err error) {
	err = DB.Omit("password", "updated_at").Where("id = ? or username LIKE ? or email LIKE ? or display_name LIKE ?", keyword, keyword+"%", keyword+"%", keyword+"%").Find(&users).Error
	return users, err
}

func SearchUsersInOrg(keyword string, org string) (users []*User, err error) {
	err = DB.Omit("password", "updated_at").Where("id = ? or username LIKE ? or email LIKE ? or display_name LIKE ?", keyword, keyword+"%", keyword+"%", keyword+"%").Find(&users, "org_id = ?", org).Error
	return users, err
}

func GetUserById(id string, selectAll bool) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("id 为空！")
	}
	user := User{Id: id}
	var err error = nil
	if selectAll {
		err = DB.First(&user, "id = ?", id).Error
	} else {
		err = DB.Omit("password").First(&user, "id = ?", id).Error
	}
	return &user, err
}

func DeleteUserById(id string) (err error) {
	if id == "" {
		return fmt.Errorf("id 为空！")
	}
	user := User{Id: id}
	return user.Delete()
}

func (user *User) Insert() error {
	err := DB.Create(user).Error
	return err
}

func (user *User) Update(updatePassword bool) error {
	var err error
	if updatePassword {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	err = DB.Model(user).Updates(user).Error
	return err
}

func (user *User) Delete() error {
	if user.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	err := DB.Delete(user).Error
	return err
}

// ValidateAndFill check password & user status
func (user *User) ValidateAndFill() (err error) {
	// When querying with struct, GORM will only query with non-zero fields,
	// that means if your field’s value is 0, '', false or other zero values,
	// it won’t be used to build query conditions
	password := user.Password
	if (user.Username == "" && user.Email == "") || password == "" {
		return fmt.Errorf("用户名或密码为空")
	}
	//DB.Where(User{Username: user.Username}).First(user)
	DB.Where("username = ? or email = ? ", user.Username, user.Email).First(user)
	okay := common.ValidatePasswordAndHash(password, user.Password)
	if !okay || user.Status != common.UserStatusEnabled {
		return fmt.Errorf("用户名或密码错误，或用户已被封禁")
	}
	return nil
}

func (user *User) FillUserById() error {
	if user.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	DB.Where(User{Id: user.Id}).First(user)
	return nil
}

func (user *User) FillUserByEmail() error {
	if user.Email == "" {
		return fmt.Errorf("email 为空！")
	}
	DB.Where(User{Email: user.Email}).First(user)
	return nil
}

func (user *User) FillUserByWeChatId() error {
	if user.WeChatId == "" {
		return fmt.Errorf("WeChat id 为空！")
	}
	DB.Where(User{WeChatId: user.WeChatId}).First(user)
	return nil
}

func (user *User) FillUserByUsername() error {
	if user.Username == "" {
		return fmt.Errorf("username 为空！")
	}
	DB.Where(User{Username: user.Username}).First(user)
	return nil
}

func IsEmailAlreadyTaken(email string) bool {
	return DB.Where("email = ?", email).Find(&User{}).RowsAffected == 1
}

func IsWeChatIdAlreadyTaken(wechatId string) bool {
	return DB.Where("wechat_id = ?", wechatId).Find(&User{}).RowsAffected == 1
}

func IsUsernameAlreadyTaken(username string) bool {
	return DB.Where("username = ?", username).Find(&User{}).RowsAffected == 1
}

func ResetUserPasswordByEmail(email string, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("邮箱地址或密码为空！")
	}
	hashedPassword, err := common.Password2Hash(password)
	if err != nil {
		return err
	}
	err = DB.Model(&User{}).Where("email = ?", email).Update("password", hashedPassword).Error
	return err
}

func IsUserEnabled(userId string) bool {
	if userId == "" {
		return false
	}
	var user User
	err := DB.Where("id = ?", userId).Select("status").Find(&user).Error
	if err != nil {
		return false
	}
	return user.Status == 1
}

func ValidateUserAccessToken(token string) (user *User) {
	if token == "" {
		return nil
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	user = &User{}
	if DB.Where("access_token = ?", token).First(user).RowsAffected == 1 {
		return user
	}
	return nil
}

func GetUserEmail(id string) (email string, err error) {
	err = DB.Model(&User{}).Where("id = ?", id).Select("email").Find(&email).Error
	return email, err
}

func GetRootUserEmail() (email string) {
	DB.Model(&User{}).Where("role = ?", "root").Select("email").Find(&email)
	return email
}

func ValidateUserToken(key string) (user *User, err error) {
	if key == "" {
		return nil, fmt.Errorf("未提供 token")
	}
	key = strings.Replace(key, "Bearer ", "", 1)
	user = &User{}
	err = DB.Where("`access_token` = ?", key).First(user).Error
	if err == nil {
		if user.Status != common.UserStatusEnabled {
			return nil, fmt.Errorf("该 user 状态不可用")
		}
		return user, nil
	}
	return nil, fmt.Errorf("无效的 token")
}
