package service

import (
	"go_blog/internal/jwt/repository"
	"go_blog/internal/model"
	"go_blog/pkg/cache"
	"go_blog/pkg/config"
	"go_blog/pkg/errorsf"
	"go_blog/pkg/utils"
	"time"

	"github.com/coocood/freecache"
	"github.com/gofrs/uuid"

	"go.uber.org/zap"
)

type JwtServiceI interface {
	SetRedisJWT(jwt string, uuid uuid.UUID) (ef errorsf.ErrInter)
	GetRedisJWT(uuid uuid.UUID) (jwtString string, ef errorsf.ErrInter)
	JoinInBlacklist(jwtList model.JwtBlacklist) (ef errorsf.ErrInter)
	LoadAll() (ef errorsf.ErrInter)
	IsInBlacklist(jwt string) bool
}

// JwtService 提供与JWT相关的服务
type JwtServiceImpl struct {
	cache      cache.CacheSqlManagerI
	log        *zap.Logger
	config     config.Config
	blackCache *freecache.Cache
	jwtRepo    repository.JwtRepositoryI
}

func NewJwtService(
	cache cache.CacheSqlManagerI,
	log *zap.Logger,
	config config.Config,
	blackCache *freecache.Cache,
	jwtRepo repository.JwtRepositoryI,
) JwtServiceI {
	return &JwtServiceImpl{
		cache:      cache,
		log:        log,
		config:     config,
		blackCache: blackCache,
		jwtRepo:    jwtRepo,
	}
}

// SetRedisJWT 将JWT存储到Redis中
func (o *JwtServiceImpl) SetRedisJWT(jwt string, uuid uuid.UUID) (ef errorsf.ErrInter) {
	// 解析配置中的JWT过期时间
	dr, err := utils.ParseDuration(o.config.Jwt.AccessTokenExpiryTime)
	if err != nil {
		ef = errorsf.DATABASE_SELECT_ERROR.Wrap(err)
	}
	expire := int64(time.Duration(dr.Seconds()))
	// 设置JWT在Redis中的过期时间
	err = o.cache.SetCache(uuid.String(), jwt, expire)
	if err != nil {
		ef = errorsf.DATABASE_SELECT_ERROR.Wrap(err)
	}
	return
}

// GetRedisJWT 从Redis中获取JWT
func (o *JwtServiceImpl) GetRedisJWT(uuid uuid.UUID) (jwtString string, ef errorsf.ErrInter) {
	// 从Redis获取指定uuid对应的JWT
	jwtString, err := o.cache.GetCache(uuid.String(), func() interface{} { return "" })
	if err != nil {
		ef = errorsf.DATABASE_SELECT_ERROR.Wrap(err)
	}
	return

}

// JoinInBlacklist 将JWT添加到黑名单
func (o *JwtServiceImpl) JoinInBlacklist(jwtList model.JwtBlacklist) (ef errorsf.ErrInter) {
	// 将JWT记录插入到数据库中的黑名单表
	o.jwtRepo.CreateJwtBlacklist(jwtList)
	// 将JWT添加到内存中的黑名单缓存
	o.blackCache.Set([]byte(jwtList.Jwt), []byte(""), 0)
	return nil
}

// IsInBlacklist 检查JWT是否在黑名单中
func (o *JwtServiceImpl) IsInBlacklist(jwt string) bool {
	// 从黑名单缓存中检查JWT是否存在
	_, err := o.blackCache.Get([]byte(jwt))
	ok := err == nil
	return ok
}

// LoadAll 从数据库加载所有的JWT黑名单并加入缓存
func (o *JwtServiceImpl) LoadAll() (ef errorsf.ErrInter) {
	var data []string
	// 从数据库中获取所有的黑名单JWT
	data, ef = o.jwtRepo.GetJwtBlacklist()
	// 将所有JWT添加到BlackCache缓存中
	for i := 0; i < len(data); i++ {
		o.blackCache.Set([]byte(data[i]), []byte(""), 0)
	}
	return ef
}
