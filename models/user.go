package models

type User struct {
	Username        string `json:"username" gorm:"primaryKey"`
	Password        string `json:"password" `
	Token           string `json:"token" gorm:"-"`
	Code            string `json:"code" gorm:"-" `
}

