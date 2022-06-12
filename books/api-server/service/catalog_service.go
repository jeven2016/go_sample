package service

import (
	"api/common"
	"api/entity"
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

func NewCatalogService(app *common.App) *CatalogService {
	return &CatalogService{
		catalog: app.Db.Collection("catalog"),
		log:     app.Log,
	}
}

func (c *CatalogService) List() ([]*entity.BookCatalog, error) {
	var results []*entity.BookCatalog
	findOpts := options.Find()
	findOpts.SetLimit(15)

	cursor, err := c.catalog.Find(context.Background(), bson.M{}, findOpts)
	defer cursor.Close(context.TODO())
	if err != nil {
		c.log.Warn("An error occurs while getting a list of catalogs", zap.Error(err))
		return results, err
	}

	for cursor.Next(context.TODO()) {
		var catalog *entity.BookCatalog
		err := cursor.Decode(&catalog)
		if err != nil {
			c.log.Warn("An error occurs while decoding a book catalog", zap.Error(err))
			return results, err
		}
		results = append(results, catalog)
	}
	return results, nil
}
