package model

type UserMigration struct {
	BaseDelete
	OldUid string `json:"old_uid" gorm:"size:128;comment:old uid"`
	NewUid string `json:"new_uid" gorm:"size:128;comment:new uid"`
	Email  string `json:"email" gorm:"size:128;comment:email"`
}
