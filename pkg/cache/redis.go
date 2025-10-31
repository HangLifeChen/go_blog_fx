package cache

import (
	"context"
	"errors"
	"fmt"
	"go_blog/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheSqlManagerI interface {
	SetCache(key string, value interface{}, configTime ...int64) error
	GetCache(key string, newCacheFunc func() interface{}, configTime ...int64) (string, error)
	DeleteCache(key ...string) error
	DelayDoubleDeleteCache(key ...string) error
	LockOperation(operationName string, lockSec int64) (isUnlock bool)
	SAdd(key string, value []interface{}, configTime ...int64) error
	SRem(key string, value []interface{}) error
	SMembersToMap(key string, newCacheFunc func() []interface{}, configTime ...int64) (map[string]struct{}, error)
	ZAdd(key string, value []redis.Z, configTime ...int64) error
	ZRangeAllToMap(key string, newCacheFunc func() []redis.Z, configTime ...int64) (map[string]float64, error)
	ZIncr(key string, value float64, member string, newCacheFunc ...func() []redis.Z) error
	LRange(key string, start, stop int64) ([]string, error)
	LPush(key string, value []interface{}) error
	Limit(key string, limitCnt int64, configTime ...int64) (isLimit bool)
}

type RedisSqlManagerImpl struct {
	rdb *redis.Client
}

func NewRedisSqlManager(rdb *redis.Client) CacheSqlManagerI {
	return &RedisSqlManagerImpl{
		rdb: rdb,
	}
}

type StandardCSMCache struct {
	D    string `json:"d"`    //user data
	L    int64  `json:"l"`    //last update time
	LW   int64  `json:"lw"`   //last work time
	W    bool   `json:"w"`    //is working update status
	Wait bool   `json:"wait"` //wait frist
}

const (
	DEFAULT_MAX_ALIVE_TIME        = int64(60 * 10) //cache will be alive for 10 minutes
	DEFAULT_FORCE_UPDATE_TIME     = int64(60 * 2)  //force update every 2 minutes
	DEFAULT_MAX_WORK_TIME         = int64(60 * 5)  //max work time for cache is 5 minutes
	DEFAULT_DELAY_DOUBLE_DEL_TIME = 1              //delay double delete
)

const (
	RdsLockOperation = "lock:operation:%s" //lock:operation:operation_name
)

func (csm *RedisSqlManagerImpl) SetCache(key string, value interface{}, configTime ...int64) error {
	maxAliveTime := DEFAULT_MAX_ALIVE_TIME
	if len(configTime) == 1 {
		maxAliveTime = configTime[0]
	}
	now := time.Now().Unix()
	newCacheData := StandardCSMCache{
		D:    utils.MarshalToString(value),
		L:    now + maxAliveTime,
		LW:   0,
		W:    false,
		Wait: false,
	}
	err := csm.rdb.Set(
		context.Background(),
		key,
		utils.MarshalToString(newCacheData),
		time.Duration(maxAliveTime)*time.Second,
	).Err()
	if err != nil {
		return err
	}
	return nil
}

func (csm *RedisSqlManagerImpl) GetCache(key string, newCacheFunc func() interface{}, configTime ...int64) (string, error) {
	forceUpdateTime := DEFAULT_FORCE_UPDATE_TIME
	maxWorkTime := DEFAULT_MAX_WORK_TIME
	maxAliveTime := DEFAULT_MAX_ALIVE_TIME
	for k, v := range configTime {
		switch k {
		case 0:
			forceUpdateTime = v
		case 1:
			maxWorkTime = v
		case 2:
			maxAliveTime = v
		}
	}
	now := time.Now().Unix()
	ctx := context.Background()
	val, err := csm.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	var cacheData StandardCSMCache
	utils.UnmarshalFromString(val, &cacheData)

	if cacheData.Wait {
		for {
			val, err := csm.rdb.Get(ctx, key).Result()
			if err != nil {
				return "", err
			}
			utils.UnmarshalFromString(val, &cacheData)
			if cacheData.D != "" {
				break
			}
			if !cacheData.Wait {
				break
			}

			time.Sleep(time.Millisecond * 3000)
			if time.Now().Unix()-now > 10000 {
				csm.rdb.Del(ctx, key)
				break
			}
		}
	}

	if cacheData.D == "" {
		var newCacheData StandardCSMCache
		newCacheData.Wait = true
		err := csm.rdb.Set(ctx, key, utils.MarshalToString(newCacheData), time.Duration(maxAliveTime)*time.Second).Err()
		if err != nil {
			return "", err
		}
		defer func() {
			if r := recover(); r != nil {
				csm.rdb.Del(ctx, key)
			}
		}()

		newData := newCacheFunc()
		newDataStr := utils.MarshalToString(newData)
		now := time.Now().Unix()
		newCacheData.D = newDataStr
		newCacheData.L = now
		newCacheData.LW = 0
		newCacheData.W = false
		newCacheData.Wait = false
		err = csm.rdb.Set(
			ctx,
			key,
			utils.MarshalToString(newCacheData),
			time.Duration(maxAliveTime)*time.Second,
		).Err()
		if err != nil {
			return "", err
		}

		return newDataStr, nil
	} else if ((now-cacheData.L >= forceUpdateTime) && !cacheData.W) || ((cacheData.LW != 0) && (now-cacheData.LW >= maxWorkTime)) {
		var newCacheData StandardCSMCache
		newCacheData.D = cacheData.D
		newCacheData.L = now
		newCacheData.LW = now
		newCacheData.W = true
		newCacheData.Wait = false

		err := csm.rdb.Set(ctx, key, utils.MarshalToString(newCacheData), time.Duration(maxAliveTime)*time.Second).Err()
		if err != nil {
			return "", err
		}

		defer func() {
			if r := recover(); r != nil {
				csm.rdb.Del(ctx, key)
			}
		}()

		newData := newCacheFunc()
		newDataStr := utils.MarshalToString(newData)
		now := time.Now().Unix()
		newCacheData.D = newDataStr
		newCacheData.L = now
		newCacheData.LW = 0
		newCacheData.W = false
		newCacheData.Wait = false

		err = csm.rdb.Set(ctx, key, utils.MarshalToString(newCacheData), time.Duration(maxAliveTime)*time.Second).Err()
		if err != nil {
			return "", err
		}
		return newDataStr, nil
	}
	return cacheData.D, nil
}

func (csm *RedisSqlManagerImpl) DeleteCache(key ...string) error {
	if len(key) == 0 {
		return nil
	}
	return csm.rdb.Del(context.Background(), key...).Err()
}

func (csm *RedisSqlManagerImpl) DelayDoubleDeleteCache(key ...string) error {
	err := csm.DeleteCache(key...)
	if err != nil {
		return err
	}
	time.AfterFunc(DEFAULT_DELAY_DOUBLE_DEL_TIME*time.Second, func() {
		csm.DeleteCache(key...)
	})
	return nil
}

func (csm *RedisSqlManagerImpl) LockOperation(operationName string, lockSec int64) (isUnlock bool) {
	key := fmt.Sprintf(RdsLockOperation, operationName)
	result, _ := csm.GetCache(key, func() interface{} {
		return 1
	})

	if result == "1" {
		csm.SetCache(key, 0, lockSec)
		return true
	}
	return false
}

func (csm *RedisSqlManagerImpl) SAdd(key string, value []interface{}, configTime ...int64) error {
	maxAliveTime := DEFAULT_MAX_ALIVE_TIME
	if len(configTime) == 1 {
		maxAliveTime = configTime[0]
	}
	ctx := context.Background()
	err := csm.rdb.SAdd(ctx, key, value...).Err()
	if err != nil {
		return err
	}
	csm.rdb.Expire(ctx, key, time.Duration(maxAliveTime)*time.Second)
	return nil
}

func (csm *RedisSqlManagerImpl) SRem(key string, value []interface{}) error {
	ctx := context.Background()
	err := csm.rdb.SRem(ctx, key, value...).Err()
	if err != nil {
		return err
	}
	return nil
}
func (csm *RedisSqlManagerImpl) SMembersToMap(key string, newCacheFunc func() []interface{}, configTime ...int64) (map[string]struct{}, error) {
	maxAliveTime := DEFAULT_MAX_ALIVE_TIME
	if len(configTime) == 1 {
		maxAliveTime = configTime[0]
	}
	ctx := context.Background()
	resMap, err := csm.rdb.SMembersMap(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if len(resMap) > 0 {
		return resMap, nil
	}
	data := newCacheFunc()
	if len(data) == 0 {
		return nil, nil
	}
	err = csm.SAdd(key, data, maxAliveTime)
	if err != nil {
		return nil, err
	}
	res := make(map[string]struct{})
	for _, v := range data {
		res[v.(string)] = struct{}{}
	}
	return res, nil
}

func (csm *RedisSqlManagerImpl) ZAdd(key string, value []redis.Z, configTime ...int64) error {
	maxAliveTime := DEFAULT_MAX_ALIVE_TIME
	if len(configTime) == 1 {
		maxAliveTime = configTime[0]
	}
	ctx := context.Background()
	err := csm.rdb.ZAdd(ctx, key, value...).Err()
	if err != nil {
		return err
	}
	csm.rdb.Expire(ctx, key, time.Duration(maxAliveTime)*time.Second)
	return nil
}

func (csm *RedisSqlManagerImpl) ZRangeAllToMap(key string, newCacheFunc func() []redis.Z, configTime ...int64) (map[string]float64, error) {
	maxAliveTime := DEFAULT_MAX_ALIVE_TIME
	if len(configTime) == 1 {
		maxAliveTime = configTime[0]
	}
	ctx := context.Background()
	resArr, err := csm.rdb.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if len(resArr) > 0 {
		tmp := make(map[string]float64)
		for _, res := range resArr {
			tmp[res.Member.(string)] = res.Score
		}
		return tmp, nil
	}
	data := newCacheFunc()
	if len(data) == 0 {
		return nil, nil
	}
	err = csm.ZAdd(key, data, maxAliveTime)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (csm *RedisSqlManagerImpl) ZIncr(key string, value float64, member string, newCacheFunc ...func() []redis.Z) error {
	ctx := context.Background()
	if len(newCacheFunc) != 0 {
		_, err := csm.ZRangeAllToMap(key, newCacheFunc[0])
		if err != nil {
			return err
		}
	}
	err := csm.rdb.ZIncrBy(ctx, key, value, member).Err()
	if err != nil {
		return err
	}
	return nil
}

func (csm *RedisSqlManagerImpl) LRange(key string, start, stop int64) ([]string, error) {
	ctx := context.Background()
	resArr, err := csm.rdb.LRange(ctx, key, start, stop).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if len(resArr) > 0 {
		return resArr, nil
	}
	return resArr, nil
}

func (csm *RedisSqlManagerImpl) LPush(key string, value []interface{}) error {
	ctx := context.Background()
	err := csm.rdb.LPush(ctx, key, value...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (csm *RedisSqlManagerImpl) Limit(key string, limitCnt int64, configTime ...int64) (isLimit bool) {
	ctx := context.Background()
	count, err := csm.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false
	}
	maxAliveTime := DEFAULT_MAX_ALIVE_TIME
	if len(configTime) == 1 {
		maxAliveTime = configTime[0]
	}
	if count == 1 && maxAliveTime > 0 {
		csm.rdb.Expire(ctx, key, time.Duration(maxAliveTime)*time.Second)
	}
	if count > limitCnt {
		return true
	}
	return false
}
