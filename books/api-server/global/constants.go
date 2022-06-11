package global

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Log /**全局Logger*/
var Log *zap.Logger

//Db mongo database
var Db *mongo.Database
