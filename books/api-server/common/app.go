package common

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type App struct {
	// Log /**全局Logger*/
	Log *zap.Logger

	//Db mongo database
	Db *mongo.Database
}

func (app *App) NewApp(log *zap.Logger, db *mongo.Database) {
	app.Log = log
	app.Db = db
}
