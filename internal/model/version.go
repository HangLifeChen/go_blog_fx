package model

type Version struct {
	Base
	Version     string `json:"version" gorm:"size:20;comment:version number"`
	RefreshTime int64  `json:"refresh_time" gorm:"comment:refresh time in seconds"`
	ClearLocal  bool   `json:"clear_local" gorm:"comment:clear local cache"`
}

const (
	RdsVersion       = "version"
	RdsVersionExpire = 60 * 60 * 24 * 30 // 30 days
)
