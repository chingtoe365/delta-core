package repository

import (
	"context"
	"delta-core/domain"
	"delta-core/mongo"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MarketRepository struct {
	redis           *redis.Client
	mongoDb         mongo.Database
	mongoCollection string
}

type TsPoint struct {
	Time  string  `json:"timestring"`
	Value float64 `json:"value"`
}

func NewMarketRepository(db mongo.Database, collection string) *MarketRepository {
	return &MarketRepository{
		redis: redis.NewClient(&redis.Options{
			Addr: "redis:6379",
			DB:   0,
		}),
		mongoDb:         db,
		mongoCollection: collection,
	}
}

// Redis transactions

func (mr *MarketRepository) FetchRawSeries(ctx context.Context, key string, start time.Time, end time.Time) []redis.TSTimestampValue {
	data, err := mr.redis.TSRange(ctx, key, int(start.UnixMilli()), int(end.UnixMilli())).Result()
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	return data
}

func (mr *MarketRepository) FetchSeries(ctx context.Context, key string, start time.Time, end time.Time) []TsPoint {
	log.Println(start.UnixMilli())
	log.Println(end.UnixMilli())
	log.Println(key)
	var result []TsPoint
	// data, err := mr.redis.TSRange(ctx, key, int(start.UnixMilli()), int(end.UnixMilli())).Result()
	// if err != nil {
	// 	log.Fatalln(err)
	// 	panic(err)
	// }
	data := mr.FetchRawSeries(ctx, key, start, end)
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

//Mongo transactions

func (mr *MarketRepository) CreateSignaler(ctx context.Context, marketSignalDto domain.MarketSignalDto) error {
	collection := mr.mongoDb.Collection(mr.mongoCollection)
	_, err := collection.InsertOne(ctx, marketSignalDto)
	return err
}

func (mr *MarketRepository) FetchSignalerById(c context.Context, signalerId string) (domain.MarketSignalDto, error) {
	collection := mr.mongoDb.Collection(mr.mongoCollection)
	var signaler domain.MarketSignalDto

	idHex, err := primitive.ObjectIDFromHex(signalerId)
	if err != nil {
		return domain.MarketSignalDto{}, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&signaler)

	return signaler, err
}

func (mr *MarketRepository) FetchByUserID(c context.Context, userID string) ([]domain.MarketSignalDto, error) {
	collection := mr.mongoDb.Collection(mr.mongoCollection)
	var signaler []domain.MarketSignalDto

	idHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(c, bson.M{"userID": idHex})
	if err != nil {
		return nil, err
	}

	err = cursor.All(c, &signaler)
	if signaler == nil {
		return []domain.MarketSignalDto{}, err
	}

	return signaler, err
}

func (mr *MarketRepository) Delete(c context.Context, signal *domain.MarketSignalDto) error {
	collection := mr.mongoDb.Collection(mr.mongoCollection)

	_, err := collection.DeleteOne(c, signal)
	return err
}
