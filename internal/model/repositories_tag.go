package model

type RepositoriesTag struct {
	Base
	Name string `json:"name" gorm:"size:100;not null;comment:repository tag name"`
	Pid  uint32 `json:"pid" gorm:"index;comment:parent id"`
}

const (
	RdsRepositoriesTag    = "repositories:tag"
	RdsRepositoriesExpire = 60 * 60 * 24 * 30 // 30 days
)
