package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gin_bbs/app/models"
	"gin_bbs/pkg/ginutils/utils"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	userCache = cache.New(30*time.Minute, 1*time.Hour)
)

// User 用户模型
type User struct {
	models.BaseModel
	Name         string `gorm:"column:name;type:varchar(255);not null"`
	Email        string `gorm:"column:email;type:varchar(255);unique;not null"`
	Avatar       string `gorm:"column:avatar;type:varchar(255);not null"`
	Introduction string `gorm:"column:introduction;type:varchar(255);not null"`
	Password     string `gorm:"column:password;type:varchar(255);not null"`
	// 是否为管理员
	IsAdmin uint `gorm:"column:is_admin;type:tinyint(1)"`
	// 用户激活
	ActivationToken string    `gorm:"column:activation_token;type:varchar(255)"`
	Activated       uint      `gorm:"column:activated;type:tinyint(1);not null"`
	EmailVerifiedAt time.Time `gorm:"column:email_verified_at"` // 激活时间
	// 用于实现记住我功能，存入 cookie 中，下次带上时，即可直接登录
	RememberToken string `gorm:"column:remember_token;type:varchar(100)"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// BeforeSave - hook
func (u *User) BeforeCreate() (err error) {
	if u.Password == "" {
		return errors.New("User Model 创建失败")
	}

	if isEncrypted := passwordEncrypted(u.Password); !isEncrypted {
		if err = u.Encrypt(); err != nil {
			return errors.New("User Model 创建失败")
		}
	}

	// 生成用户 remember_token
	if u.RememberToken == "" {
		u.RememberToken = string(utils.RandomCreateBytes(10))
	}

	// 生成用户激活 token
	if u.ActivationToken == "" {
		u.ActivationToken = string(utils.RandomCreateBytes(30))
	}

	// 生成用户头像
	if u.Avatar == "" {
		hash := md5.Sum([]byte(u.Email))
		u.Avatar = "http://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
	}

	return err
}

// BeforeUpdate - hook
func (u *User) BeforeUpdate() (err error) {
	if isEncrypted := passwordEncrypted(u.Password); !isEncrypted {
		if err = u.Encrypt(); err != nil {
			return errors.New("User Model 更新失败")
		}
	}

	return
}

// ------------ private
func passwordEncrypted(pwd string) (status bool) {
	return len(pwd) == 60 // 长度等于 60 说明加密过了
}

func setToCache(user *User) {
	key := strconv.Itoa(int(user.ID))
	userCache.Set(key, user, cache.DefaultExpiration)
}

func getFromCache(id int) (*User, bool) {
	cachedUser, ok := userCache.Get(strconv.Itoa(id))
	if !ok {
		return nil, false
	}

	u, ok := cachedUser.(*User)
	if !ok {
		return nil, false
	}

	return u, true
}

func delCache(id int) {
	userCache.Delete(strconv.Itoa(id))
}
