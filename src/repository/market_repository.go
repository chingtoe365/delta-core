package repository

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RedisMarketRepository struct {
	client *redis.Client
}

func NewRedisMarketRepository() *RedisMarketRepository {
	// url := "redis://redis:6379/?db=0"
	// opts, err := redis.ParseURL(url)
	// if err != nil {
	// 	panic(err)
	// }
	return &RedisMarketRepository{
		client: redis.NewClient(&redis.Options{
			Addr: "redis:6379",
			DB:   0,
			// DialTimeout: 1,
			// PoolTimeout: 3,
		}),
	}
}

func (rc *RedisMarketRepository) FetchSeries(key string, start time.Time, end time.Time) string {
	log.Println(start.UnixMilli())
	log.Println(end.UnixMilli())
	log.Println(key)
	data, err := rc.client.GetRange(key, start.UnixMilli(), end.UnixMilli()).Result()
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	log.Println(data)
	return data
}
