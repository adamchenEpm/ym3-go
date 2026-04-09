package redis

import (
	"context"
	"fmt"
	"github.com/adamchenEpm/ym3-go/internal/config"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

// RedisUtil Redis 操作工具类
type RedisUtil struct {
	client *redis.Client
}

var (
	instance *RedisUtil
	once     sync.Once
)

// GetInstance 获取 Redis 工具单例（自动初始化连接池）
func GetInstance() *RedisUtil {
	once.Do(func() {
		instance = &RedisUtil{}
		if err := instance.init(); err != nil {
			panic(fmt.Errorf("初始化 Redis 失败: %w", err))
		}
	})
	return instance
}

// init 内部初始化 Redis 客户端
func (r *RedisUtil) init() error {
	cfg := config.Get()
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  time.Duration(cfg.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Redis.WriteTimeout) * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return err
	}
	r.client = rdb
	return nil
}

// GetClient 返回原始 redis.Client 对象（供高级操作）
func (r *RedisUtil) GetClient() *redis.Client {
	return r.client
}

// Close 关闭 Redis 连接
func (r *RedisUtil) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// ================== 基本操作 ==================

// Set 设置键值对，带过期时间（可选）
func (r *RedisUtil) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(context.Background(), key, value, expiration).Err()
}

// Get 获取字符串值
func (r *RedisUtil) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

// Del 删除一个或多个键
func (r *RedisUtil) Del(keys ...string) error {
	return r.client.Del(context.Background(), keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisUtil) Exists(key string) (bool, error) {
	n, err := r.client.Exists(context.Background(), key).Result()
	return n > 0, err
}

// Expire 设置键的过期时间
func (r *RedisUtil) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(context.Background(), key, expiration).Err()
}

// TTL 获取键的剩余生存时间
func (r *RedisUtil) TTL(key string) (time.Duration, error) {
	return r.client.TTL(context.Background(), key).Result()
}

// ================== 哈希操作 ==================

// HSet 设置哈希字段值
func (r *RedisUtil) HSet(key string, values ...interface{}) error {
	return r.client.HSet(context.Background(), key, values...).Err()
}

// HGet 获取哈希字段值
func (r *RedisUtil) HGet(key, field string) (string, error) {
	return r.client.HGet(context.Background(), key, field).Result()
}

// HGetAll 获取所有哈希字段
func (r *RedisUtil) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(context.Background(), key).Result()
}

// HDel 删除哈希字段
func (r *RedisUtil) HDel(key string, fields ...string) error {
	return r.client.HDel(context.Background(), key, fields...).Err()
}

// ================== 列表操作 ==================

// LPush 从左侧推入元素
func (r *RedisUtil) LPush(key string, values ...interface{}) error {
	return r.client.LPush(context.Background(), key, values...).Err()
}

// RPush 从右侧推入元素
func (r *RedisUtil) RPush(key string, values ...interface{}) error {
	return r.client.RPush(context.Background(), key, values...).Err()
}

// LPop 从左侧弹出元素
func (r *RedisUtil) LPop(key string) (string, error) {
	return r.client.LPop(context.Background(), key).Result()
}

// RPop 从右侧弹出元素
func (r *RedisUtil) RPop(key string) (string, error) {
	return r.client.RPop(context.Background(), key).Result()
}

// LRange 获取列表片段
func (r *RedisUtil) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(context.Background(), key, start, stop).Result()
}

// ================== 集合操作 ==================

// SAdd 添加集合成员
func (r *RedisUtil) SAdd(key string, members ...interface{}) error {
	return r.client.SAdd(context.Background(), key, members...).Err()
}

// SMembers 获取所有集合成员
func (r *RedisUtil) SMembers(key string) ([]string, error) {
	return r.client.SMembers(context.Background(), key).Result()
}

// SIsMember 判断是否为集合成员
func (r *RedisUtil) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(context.Background(), key, member).Result()
}

// SRem 移除集合成员
func (r *RedisUtil) SRem(key string, members ...interface{}) error {
	return r.client.SRem(context.Background(), key, members...).Err()
}

// ================== 有序集合操作 ==================

// ZAdd 添加有序集合成员（带分数）
func (r *RedisUtil) ZAdd(key string, members ...redis.Z) error {
	return r.client.ZAdd(context.Background(), key, members...).Err()
}

// ZRange 按索引范围获取成员（升序）
func (r *RedisUtil) ZRange(key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(context.Background(), key, start, stop).Result()
}

// ZRangeByScore 按分数范围获取成员
func (r *RedisUtil) ZRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRangeByScore(context.Background(), key, opt).Result()
}

// ZRem 删除有序集合成员
func (r *RedisUtil) ZRem(key string, members ...interface{}) error {
	return r.client.ZRem(context.Background(), key, members...).Err()
}

// ================== 分布式锁（简单实现）==================

// Lock 获取分布式锁（基于 SETNX）
// 返回值：是否获取成功，以及解锁函数
func (r *RedisUtil) Lock(key string, expiration time.Duration) (bool, func(), error) {
	ctx := context.Background()
	ok, err := r.client.SetNX(ctx, key, "locked", expiration).Result()
	if err != nil {
		return false, nil, err
	}
	if ok {
		unlock := func() {
			r.client.Del(ctx, key)
		}
		return true, unlock, nil
	}
	return false, nil, nil
}
