package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"tools-admin/backend/common/config"
	"tools-admin/backend/pkg/log"
)

var Redis *RClient

type RClient struct {
	client *redis.Client
	ctx    context.Context
}

func init() {
	rConfig := config.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rConfig.Host, rConfig.Port),
		Password: rConfig.Password,
		DB:       rConfig.DB,
	})

	redisClient := &RClient{
		client: client,
		ctx:    context.Background(),
	}
	Redis = redisClient
	_, err := client.Ping(Redis.ctx).Result()
	if err != nil {
		log.Info("redis连接失败！")
		return
	}
	log.Info("redis连接成功！")

}

// Set 设置 key的值
func (rc *RClient) Set(key, value string) bool {
	result, err := rc.client.Set(rc.ctx, key, value, 0).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false
	}
	return result == "OK"
}

// SetEX 设置 key的值并指定过期时间
func (rc *RClient) SetEX(key, value string, ex time.Duration) bool {
	result, err := rc.client.Set(rc.ctx, key, value, ex).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false
	}
	return result == "OK"
}

// Get 获取 key的值
func (rc *RClient) Get(key string) (bool, string) {
	result, err := rc.client.Get(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false, ""
	}
	return true, result
}

// GetSet 设置新值获取旧值
func (rc *RClient) GetSet(key, value string) (bool, string) {
	oldValue, err := rc.client.GetSet(rc.ctx, key, value).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false, ""
	}
	return true, oldValue
}

// Incr key值每次加一 并返回新值
func (rc *RClient) Incr(key string) int64 {
	val, err := rc.client.Incr(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return val
}

// IncrBy key值每次加指定数值 并返回新值
func (rc *RClient) IncrBy(key string, incr int64) int64 {
	val, err := rc.client.IncrBy(rc.ctx, key, incr).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return val
}

// IncrByFloat key值每次加指定浮点型数值 并返回新值
func (rc *RClient) IncrByFloat(key string, incrFloat float64) float64 {
	val, err := rc.client.IncrByFloat(rc.ctx, key, incrFloat).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return val
}

// Decr key值每次递减 1 并返回新值
func (rc *RClient) Decr(key string) int64 {
	val, err := rc.client.Decr(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return val
}

// DecrBy key值每次递减指定数值 并返回新值
func (rc *RClient) DecrBy(key string, incr int64) int64 {
	val, err := rc.client.DecrBy(rc.ctx, key, incr).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return val
}

// Del 删除 key
func (rc *RClient) Del(key string) bool {
	result, err := rc.client.Del(rc.ctx, key).Result()
	if err != nil {
		return false
	}
	return result == 1
}

// Expire 设置 key的过期时间
func (rc *RClient) Expire(key string, ex time.Duration) bool {
	result, err := rc.client.Expire(rc.ctx, key, ex).Result()
	if err != nil {
		return false
	}
	return result
}

/*------------------------------------ list 操作 ------------------------------------*/

// LPush 从列表左边插入数据，并返回列表长度
func (rc *RClient) LPush(key string, date ...interface{}) int64 {
	result, err := rc.client.LPush(rc.ctx, key, date).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return result
}

// RPush 从列表右边插入数据，并返回列表长度
func (rc *RClient) RPush(key string, date ...interface{}) int64 {
	result, err := rc.client.RPush(rc.ctx, key, date).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return result
}

// LPop 从列表左边删除第一个数据，并返回删除的数据
func (rc *RClient) LPop(key string) (bool, string) {
	val, err := rc.client.LPop(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false, ""
	}
	return true, val
}

// RPop 从列表右边删除第一个数据，并返回删除的数据
func (rc *RClient) RPop(key string) (bool, string) {
	val, err := rc.client.RPop(rc.ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	return true, val
}

// LIndex 根据索引坐标，查询列表中的数据
func (rc *RClient) LIndex(key string, index int64) (bool, string) {
	val, err := rc.client.LIndex(rc.ctx, key, index).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false, ""
	}
	return true, val
}

// LLen 返回列表长度
func (rc *RClient) LLen(key string) int64 {
	val, err := rc.client.LLen(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return val
}

// LRange 返回列表的一个范围内的数据，也可以返回全部数据
func (rc *RClient) LRange(key string, start, stop int64) []string {
	vales, err := rc.client.LRange(rc.ctx, key, start, stop).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return vales
}

// LRem 从列表左边开始，删除元素data， 如果出现重复元素，仅删除 count次
func (rc *RClient) LRem(key string, count int64, data interface{}) bool {
	_, err := rc.client.LRem(rc.ctx, key, count, data).Result()
	if err != nil {
		fmt.Println(err)
	}
	return true
}

// LInsert 在列表中 pivot 元素的后面插入 data
func (rc *RClient) LInsert(key string, pivot int64, data interface{}) bool {
	err := rc.client.LInsert(rc.ctx, key, "after", pivot, data).Err()
	if err != nil {
		log.Error("redis err : %v", err)
		return false
	}
	return true
}

/*------------------------------------ set 操作 ------------------------------------*/

// SAdd 添加元素到集合中
func (rc *RClient) SAdd(key string, data ...interface{}) bool {
	err := rc.client.SAdd(rc.ctx, key, data).Err()
	if err != nil {
		log.Error("redis err : %v", err)
		return false
	}
	return true
}

// SCard 获取集合元素个数
func (rc *RClient) SCard(key string) int64 {
	size, err := rc.client.SCard(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return size
}

// SIsMember 判断元素是否在集合中
func (rc *RClient) SIsMember(key string, data interface{}) bool {
	ok, err := rc.client.SIsMember(rc.ctx, key, data).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return ok
}

// SMembers 获取集合所有元素
func (rc *RClient) SMembers(key string) []string {
	es, err := rc.client.SMembers(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return es
}

// SRem 删除 key集合中的 data元素
func (rc *RClient) SRem(key string, data ...interface{}) bool {
	_, err := rc.client.SRem(rc.ctx, key, data).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false
	}
	return true
}

// SPopN 随机返回集合中的 count个元素，并且删除这些元素
func (rc *RClient) SPopN(key string, count int64) []string {
	vales, err := rc.client.SPopN(rc.ctx, key, count).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return vales
}

/*------------------------------------ hash 操作 ------------------------------------*/

// HSet 根据 key和 field字段设置，field字段的值
func (rc *RClient) HSet(key, field, value string) bool {
	err := rc.client.HSet(rc.ctx, key, field, value).Err()
	if err != nil {
		return false
	}
	return true
}

// HGet 根据 key和 field字段，查询field字段的值
func (rc *RClient) HGet(key, field string) string {
	val, err := rc.client.HGet(rc.ctx, key, field).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return val
}

// HMGet 根据key和多个字段名，批量查询多个 hash字段值
func (rc *RClient) HMGet(key string, fields ...string) []interface{} {
	vales, err := rc.client.HMGet(rc.ctx, key, fields...).Result()
	if err != nil {
		panic(err)
	}
	return vales
}

// HGetAll 根据 key查询所有字段和值
func (rc *RClient) HGetAll(key string) map[string]string {
	data, err := rc.client.HGetAll(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return data
}

// HKeys 根据 key返回所有字段名
func (rc *RClient) HKeys(key string) []string {
	fields, err := rc.client.HKeys(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return fields
}

// HLen 根据 key，查询hash的字段数量
func (rc *RClient) HLen(key string) int64 {
	size, err := rc.client.HLen(rc.ctx, key).Result()
	if err != nil {
		log.Error("redis err : %v", err)
	}
	return size
}

// HMSet 根据 key和多个字段名和字段值，批量设置 hash字段值
func (rc *RClient) HMSet(key string, data map[string]interface{}) bool {
	result, err := rc.client.HMSet(rc.ctx, key, data).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false
	}
	return result
}

// HSetNX 如果 field字段不存在，则设置 hash字段值
func (rc *RClient) HSetNX(key, field string, value interface{}) bool {
	result, err := rc.client.HSetNX(rc.ctx, key, field, value).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false
	}
	return result
}

// HDel 根据 key和字段名，删除 hash字段，支持批量删除
func (rc *RClient) HDel(key string, fields ...string) bool {
	_, err := rc.client.HDel(rc.ctx, key, fields...).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false
	}
	return true
}

// HExists 检测 hash字段名是否存在
func (rc *RClient) HExists(key, field string) bool {
	result, err := rc.client.HExists(rc.ctx, key, field).Result()
	if err != nil {
		log.Error("redis err : %v", err)
		return false
	}
	return result
}
