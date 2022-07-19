package service

import (
	"api/app/books/dto"
	"api/app/books/entitie"
	"api/pkg/common"
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

func (c *CatalogService) List() (*dto.CatalogListResponse, error) {
	var results []*entitie.BookCatalog
	findOpts := options.Find()
	findOpts.SetLimit(15)

	cursor, err := c.catalog.Find(context.Background(), bson.M{}, findOpts)
	defer cursor.Close(context.TODO())
	if err != nil {
		c.log.Warn("An error occurs while getting a list of catalogs", zap.Error(err))
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var catalog *entitie.BookCatalog
		err := cursor.Decode(&catalog)
		if err != nil {
			c.log.Warn("An error occurs while decoding a book catalog", zap.Error(err))
			return nil, err
		}
		results = append(results, catalog)
	}

	var catalogList = &dto.CatalogListResponse{
		Count: int32(len(results)),
		List:  results,
	}
	return catalogList, nil
}
