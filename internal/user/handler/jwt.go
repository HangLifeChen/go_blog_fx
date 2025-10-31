package handler

import (
	"go_blog/internal/jwt/service"
	"go_blog/internal/model"
	"go_blog/pkg/errorsf"
	"go_blog/pkg/utils"

	"github.com/gofrs/uuid"
)

func NewJwtHandler(
	response *utils.Response,
	jwtSrv service.JwtServiceI,
) *JwtHandler {
	return &JwtHandler{
		response: response,
		jwtSrv:   jwtSrv,
	}
}

type JwtHandler struct {
	response *utils.Response
	jwtSrv   service.JwtServiceI
}

// SetRedisJWT 将JWT存储到Redis中
func (o *JwtHandler) SetRedisJWT(jwt string, uuid uuid.UUID) errorsf.ErrInter {
	return o.jwtSrv.SetRedisJWT(jwt, uuid)
}

// GetRedisJWT 从Redis中获取JWT
func (o *JwtHandler) GetRedisJWT(uuid uuid.UUID) (string, errorsf.ErrInter) {
	// 从Redis获取指定uuid对应的JWT
	return o.jwtSrv.GetRedisJWT(uuid)
}

// JoinInBlacklist 将JWT添加到黑名单
func (o *JwtHandler) JoinInBlacklist(jwtList model.JwtBlacklist) errorsf.ErrInter {
	return o.jwtSrv.JoinInBlacklist(jwtList)
}

// IsInBlacklist 检查JWT是否在黑名单中
func (o *JwtHandler) IsInBlacklist(jwt string) bool {
	return o.jwtSrv.IsInBlacklist(jwt)
}

// LoadAll 从数据库加载所有的JWT黑名单并加入缓存
func (o *JwtHandler) LoadAll() {
	o.jwtSrv.LoadAll()
}
