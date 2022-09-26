package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"api/app/books/dto"
	"api/app/books/entity"
	"api/pkg/clients"
	"api/pkg/global"
)

type ICatalog interface {
	insert()
}

type CatalogService struct {
	catalog *mongo.Collection
	log     *zap.Logger
	redis   *clients.RedisClient
}

func NewCatalogService(app *global.App) *CatalogService {
	return &CatalogService{
		catalog: app.MongoClient.Db.Collection("catalog"),
		log:     app.Log,
		redis:   app.RedisClient,
	}
}

func (c *CatalogService) List() (*dto.CatalogListResponse, error) {
	// return from cache if it exists
	redisResult := c.redis.Client.HSet(context.Background(), "catalogTree")
	if redisResult.

	var results = new([]*entity.BookCatalog)
	findOpts := options.Find()
	// findOpts.SetLimit(15)

	cursor, err := c.catalog.Find(context.Background(), bson.M{}, findOpts)
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			c.log.Warn("An error occurs while cursor closing", zap.Error(err))
		}
	}(cursor, context.TODO())

	if err != nil {
		c.log.Warn("An error occurs while getting a list of catalogs", zap.Error(err))
		return nil, err
	}

	err = cursor.All(context.Background(), results)
	if err != nil {
		c.log.Warn("An error occurs while decoding a book catalog", zap.Error(err))
	}
	return CatalogTree(*results), nil
}
