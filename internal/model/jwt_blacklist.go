package model

// JwtBlacklist JWT 黑名单表
type JwtBlacklist struct {
	BaseDelete
	Jwt string `json:"jwt" gorm:"type:text"` // Jwt
}
