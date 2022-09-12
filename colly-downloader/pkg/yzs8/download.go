package yzs8

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"colly-downloader/pkg/models"
)

// ensureCatalog, return existed, error
func ensureCatalog(collection *mongo.Collection, catalog *models.CatalogDoc) (string, error) {
	result := collection.FindOne(context.TODO(), bson.M{"name": catalog.Name})
	err := result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			insertResult, err := collection.InsertOne(context.TODO(), catalog)

			return insertResult.InsertedID.(string), err
		}
		return "", err
	}
	catalog = new(models.CatalogDoc)
	err = result.Decode(catalog)
	if err != nil {
		return catalog.Id, nil
	}
	return "", nil
}

func GoDown() {
	log := initLog()
	db, err := CreateMongoClient(log)
	if err != nil {
		return
	}

	catalogCol := db.Collection("catalog")

	individualCatalog := &models.CatalogDoc{
		ParentId:     "",
		Name:         "yzs8",
		Order:        1,
		ArticleCount: 0,
		Description:  "yzs8",
		CreateDate:   time.Now(),
		LastUpdate:   time.Now(),
	}

	catalogId, err := ensureCatalog(catalogCol, individualCatalog)
	if err != nil {
		log.Error("failed to check the catalog", zap.Error(err))
		return
	}
	log.Info("catalog id is", zap.String("id", catalogId))

}

// https://yazhouse8.com/article.php?cate=1  首页
// 文章：<div class="articleList"><p><span>12、</span>
// <a class="img-center" target="_blank" href="article/27078.html">有雪</a></p></div>
// 下一页 ： https://yazhouse8.com/article.php?page=2&cate=1
// check： <ul class="pager"><li><a href="article.php?cate=1">首页</a>&nbsp;&nbsp;
//      </li><li><a href="article.php?page=2&amp;cate=1">下一页</a></li></ul>
