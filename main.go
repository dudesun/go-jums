package jums

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type Jums struct {
	config Config
}

type Config struct {
	Key           string
	Secret        string
	EndPoint      string
	AccountKey    string
	AccountSecret string
}

// New 极光统一推送
func New(config Config) *Jums {
	return &Jums{
		config: config,
	}
}

// Message 普通消息推送
type Message struct {
	Key      string
	Secret   string
	EndPoint string
	Data   map[string]any
}

// Message 消息模式
func (u *Jums) Message() *Message {
	return &Message{
		Key:      u.config.Key,
		Secret:   u.config.Secret,
		EndPoint: u.config.EndPoint,
		Data:   map[string]any{},
	}
}

type Users struct {
	Key           string
	Secret        string
	EndPoint      string
	AccountKey    string
	AccountSecret string
	Data          map[string]any
	UserId        uint
}

// User 用户模式
func (u *Jums) User() *Users {
	return &Users{
		Key:           u.config.Key,
		Secret:        u.config.Secret,
		EndPoint:      u.config.EndPoint,
		AccountKey:    u.config.AccountKey,
		AccountSecret: u.config.AccountSecret,
		Data:          map[string]any{},
	}
}

// UserDel 批量删除用户
func (u *Jums) UserDel(userid ...uint) error {
	if len(userid) == 0 {
		return errors.New("删除用户为空")
	}
	return Request("v1/user/delete", u.config.EndPoint, u.config.AccountKey, u.config.AccountSecret, userid)
}

// Request 请求数据
func Request(url string, endpoint string, key string, secret string, data any) error {
	var err error
	_, body, errs := fiber.Post(endpoint+url).Debug().BasicAuth(key, secret).JSON(data).Bytes()
	if len(errs) > 0 {
		return errs[0]
	}
	var api struct {
		Code    int
		Message string
	}
	err = json.Unmarshal(body, &api)
	if err != nil {
		return errors.New("jums request failed")
	}

	if api.Code != 0 {
		return errors.New(api.Message)
	}
	return nil
}
