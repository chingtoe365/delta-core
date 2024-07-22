package repository

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisMarketRepository struct {
	client *redis.Client
}

type TsPoint struct {
	Time  string  `json:"timestring"`
	Value float64 `json:"value"`
}

func NewRedisMarketRepository() *RedisMarketRepository {
	return &RedisMarketRepository{
		client: redis.NewClient(&redis.Options{
			Addr: "redis:6379",
			DB:   0,
		}),
	}
}

func (rc *RedisMarketRepository) FetchSeries(ctx context.Context, key string, start time.Time, end time.Time) []TsPoint {
	log.Println(start.UnixMilli())
	log.Println(end.UnixMilli())
	log.Println(key)
	var result []TsPoint
	data, err := rc.client.TSRange(ctx, key, int(start.UnixMilli()), int(end.UnixMilli())).Result()
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	// log.Println(data)
	// return data
	for x := range len(data) {
		result = append(result, TsPoint{
			Time:  time.UnixMilli(data[x].Timestamp).Format(time.RFC3339),
			Value: data[x].Value,
		})

	}
	return result
}
