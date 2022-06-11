package service

import (
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
}

func (c *CatalogService) New() *CatalogService {
	return &CatalogService{
		catalog: global.Db.Collection("catalog"),
	}
}

func (c *CatalogService) list() {
	findOpts := options.Find()
	 findOpts.SetLimit(15)

	cursor, err := c.catalog.Find(context.Background(), bson.D{}, findOpts)
	if err != nil {
		global.Log.Warn("An error occurs while getting a list of catalogs", zap.Error(err))
	}
	for
}
