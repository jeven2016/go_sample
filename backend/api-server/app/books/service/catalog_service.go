package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"api/app/books/book_common"
	"api/app/books/dto"
	"api/app/books/entity"
	"api/pkg/clients"
	"api/pkg/global"
)

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
	// retrieve from cache if it exists
	data, err := c.retrieveCache()
	if err == nil || err != redis.Nil {
		return data, err
	}

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
		return nil, err
	}
	var finalResult = CatalogTree(*results)

	// store to cache
	jsonData, err := convertor.ToJson(finalResult)
	if err == nil {
		_, err = c.redis.Client.HSetNX(context.Background(), book_common.KEY_CATALOG_TREE,
			book_common.FIELD_JSON_DATA, jsonData).Result()
		if err == nil {
			c.redis.Client.Expire(context.Background(), book_common.KEY_CATALOG_TREE, 30*time.Minute)
		}
	}
	c.log.Warn("An error occurs while retrieving CatalogListResponse", zap.Error(err))
	return finalResult, nil
}

func (c *CatalogService) retrieveCache() (*dto.CatalogListResponse, error) {
	redisResult := c.redis.Client.HGet(context.Background(), book_common.KEY_CATALOG_TREE,
		book_common.FIELD_JSON_DATA)
	jsonData, err := redisResult.Result()

	if err == nil {
		var catResp dto.CatalogListResponse
		err := json.Unmarshal([]byte(jsonData), &catResp)
		return &catResp, err
	}
	return nil, err
}
