package service

import (
	"api/entity"
	"api/global"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type ICatalog interface {
	insert()
}

type CatalogService struct {
	catalog *mongo.Collection
	log     *zap.Logger
}

func New(log *zap.Logger) *CatalogService {
	return &CatalogService{
		catalog: global.Db.Collection("catalog"),
		log:     log,
	}
}

func (c *CatalogService) List() []*entity.BookCatalog {
	var results []*entity.BookCatalog
	findOpts := options.Find()
	findOpts.SetLimit(15)

	cursor, err := c.catalog.Find(context.Background(), bson.D{}, findOpts)
	defer cursor.Close(context.TODO())
	if err != nil {
		c.log.Warn("An error occurs while getting a list of catalogs", zap.Error(err))
		return results
	}

	for cursor.Next(context.TODO()) {
		var catalog *entity.BookCatalog
		err := cursor.Decode(catalog)
		if err != nil {
			c.log.Warn("An error occurs while decoding a book catalog", zap.Error(err))
		}
		results = append(results, catalog)
	}
	return results
}
