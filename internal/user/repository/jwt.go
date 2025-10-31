package repository

import (
	"go_blog/internal/model"
	"go_blog/pkg/errorsf"

	"gorm.io/gorm"

	"go.uber.org/zap"
)

type JwtRepositoryI interface {
	CreateJwtBlacklist(jwt model.JwtBlacklist) (ef errorsf.ErrInter)
	GetJwtBlacklist() (blacklist []string, ef errorsf.ErrInter)
}

// JwtService 提供与JWT相关的服务
type JwtService struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewJwtRepository(db *gorm.DB) *JwtService {
	return &JwtService{
		db: db,
	}
}

func (jwtService *JwtService) CreateJwtBlacklist(jwtList model.JwtBlacklist) (ef errorsf.ErrInter) {
	// 将JWT记录插入到数据库中的黑名单表
	if err := jwtService.db.Create(&jwtList).Error; err != nil {
		ef = errorsf.DATABASE_INSERT_ERROR.Wrap(err)
		return
	}
	return
}

// LoadAll 从数据库加载所有的JWT黑名单并加入缓存
func (jwtService *JwtService) GetJwtBlacklist() (blacklist []string, ef errorsf.ErrInter) {
	// 从数据库中获取所有的黑名单JWT
	if err := jwtService.db.Model(&model.JwtBlacklist{}).Pluck("jwt", &blacklist).Error; err != nil {
		// 如果获取失败，记录错误日志
		jwtService.log.Error("Failed to load JWT blacklist from the database", zap.Error(err))
		ef = errorsf.GET_JWT_BLACKLIST_FAILED.Wrap(err)
		return
	}
	return
}
