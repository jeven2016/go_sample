package global

import (
	"api/pkg/clients"
	"go.uber.org/zap"
)

type App struct {
	// Log /**全局Logger*/
	Log *zap.Logger

	//Db mongo database
	//Db *mongo.Database

	RedisClient *clients.RedisClient

	MongoClient *clients.MongoClient

	AuthClient *clients.AuthClient
}
