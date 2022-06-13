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

func (artSrv *ArticleService) List(catalogId string) ([]*entity.Article, error) {
	var results []*entity.Article
	findOpt := options.Find()
	findOpt.SetLimit(1)

	cursor, err := artSrv.article.Find(context.TODO(), bson.M{"catalogId": catalogId}, findOpt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	if err != nil {
		artSrv.log.Warn("An error occurs while getting a list of articles", zap.Error(err))
		return results, err
	}

	for cursor.Next(context.TODO()) {
		var article *entity.Article
		err := cursor.Decode(&article)
		if err != nil {
			artSrv.log.Warn("An error occurs while decoding a book article", zap.Error(err))
			return results, err
		}
		results = append(results, article)
	}
	return results, nil
}
