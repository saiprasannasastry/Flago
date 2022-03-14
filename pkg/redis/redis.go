package redis

import (
	"github.com/apex/log"
	"github.com/go-redis/redis"
	"time"
)

const (
	DISABLE_ALL_KEY = "disable:all"
	ENABLE_ALL_KEY  = "enable:all"
)

type RedisOpts struct {
	Password string
}

type Pool struct {
	RedisClient *redis.Client
}
type PoolInterface interface {
	EnableAllCustomers(feature string) error
	DisableAllCustomers(feature string) error
	AddToSetOfcustomers(customerName string, customerId string, feature string) error
	AddToRef(refType string, feature string) error
	GetAllCustomers(refType string) ([]string, error)
	GetFlagForCustomer(customerDetails, feature string) (bool, error)
}

func NewPool(redisClient *redis.Client) PoolInterface {
	return Pool{RedisClient: redisClient}
}

func GetRedisClient(redisAddr string, maxConnections int, timeout time.Duration, opts ...RedisOpts) (*redis.Client, error) {
	redisPassword := ""
	for _, opt := range opts {
		redisPassword = opt.Password
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     redisPassword,
		MinIdleConns: maxConnections,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("failed to connect to redis %v", err)
		return nil, err
	}
	return rdb, nil
}
