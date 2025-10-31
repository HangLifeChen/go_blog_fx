package model

// ArticleLike 文章收藏表
type ArticleLike struct {
	BaseDelete
	ArticleID string `json:"article_id"` // 文章 ID
	UserID    uint   `json:"user_id"`    // 用户 ID
	User      User   `json:"-" gorm:"foreignKey:UserID"`
}
