package model

import "gorm.io/plugin/soft_delete"

// 管理员表
type Admin struct {
	BaseNoDelete
	DeletedAt     soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:idx_username_deleted_at"`
	Username      string                `json:"username" gorm:"uniqueIndex:idx_username_deleted_at;size:20"`
	Password      EncryptedString       `json:"password" gorm:"size:64"`
	Nickname      string                `json:"nickname" gorm:"size:20"`
	Avatar        string                `json:"avatar" gorm:"size:255"`
	LastLoginTime int64                 `json:"last_login_time"`
	LastLoginIp   string                `json:"last_login_ip" gorm:"size:40"`
}
