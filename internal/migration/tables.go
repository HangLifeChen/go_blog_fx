package migration

import (
	"go_blog/internal/model"

	"gorm.io/gorm"
)

func initTables(db *gorm.DB) (err error) {
	if err = db.Set("gorm:table_options", "COMMENT='user table'").AutoMigrate(&model.User{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='advertisement table'").AutoMigrate(&model.Advertisement{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='article category'").AutoMigrate(&model.ArticleCategory{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='article like '").AutoMigrate(&model.ArticleLike{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='article tag'").AutoMigrate(&model.ArticleTag{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='repository tag table'").AutoMigrate(&model.RepositoriesTag{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='user comment'").AutoMigrate(&model.Comment{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='Feedback'").AutoMigrate(&model.Feedback{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='browser footer link'").AutoMigrate(&model.FooterLink{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='friend link'").AutoMigrate(&model.FriendLink{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='image'").AutoMigrate(&model.Image{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='jwt blacklist'").AutoMigrate(&model.JwtBlacklist{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='login record'").AutoMigrate(&model.Login{}); err != nil {
		return
	}

	if err = db.Set("gorm:table_options", "COMMENT='admin user table'").AutoMigrate(&model.Admin{}); err != nil {
		return
	}
	if err = db.Set("gorm:table_options", "COMMENT='version table'").AutoMigrate(&model.Version{}); err != nil {
		return
	}

	if err = db.Set("gorm:table_options", "COMMENT='user migration table'").AutoMigrate(&model.UserMigration{}); err != nil {
		return
	}
	return
}
