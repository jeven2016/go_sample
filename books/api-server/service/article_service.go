package service

import (
	"api/common"
	"api/entity"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ArticleService struct {
	article *mongo.Collection
	log     *zap.Logger
}

func NewArticleService(app *common.App) *ArticleService {
	return &ArticleService{
		article: app.Db.Collection("article"),
		log:     app.Log,
	}
}

func (artSrv ArticleService) FindById(id string) (*entity.Article, error) {
	articleEntity := &entity.Article{}
	err := artSrv.article.FindOne(context.TODO(), bson.M{"Id": id}).Decode(articleEntity)
	if err != nil {
		artSrv.log.Warn("error occurs findById(id)", zap.String("id", id), zap.Error(err))
	}
	return articleEntity, err
}

func (artSrv *ArticleService) List() ([]*entity.Article, error) {
	var results []*entity.Article
	return results, nil
}
