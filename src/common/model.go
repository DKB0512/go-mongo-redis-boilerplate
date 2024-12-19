package common

import (
	_redis "github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModelConstructor struct {
	Collection *mongo.Collection
	Redis      *_redis.Client
}
