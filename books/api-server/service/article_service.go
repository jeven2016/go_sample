package service

import (
	"api/common"
	"api/dto"
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
	result := artSrv.article.FindOne(context.TODO(), bson.M{"_id": id})
	err := result.Err()
	if err != nil {
		artSrv.log.Warn("error occurs when findById(id)", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	err = result.Decode(articleEntity)
	if err != nil {
		artSrv.log.Warn("error occurs when Decode(articleEntity)", zap.String("id", id), zap.Error(err))
	}
	return articleEntity, err
}

func (artSrv *ArticleService) List(catalogId string) ([]*dto.ArticlePageResponse, error) {
	var results []*entity.Article
	findOpt := options.Find()
	findOpt.SetLimit(10)
	findOpt.SetProjection(bson.M{"content": 0}) //不包含content内容

	cursor, err := artSrv.article.Find(context.TODO(), bson.M{"catalogId": catalogId}, findOpt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	if err != nil {
		artSrv.log.Warn("An error occurs while getting a list of articles", zap.Error(err))
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var article *entity.Article
		err := cursor.Decode(&article)
		if err != nil {
			artSrv.log.Warn("An error occurs while decoding a book article", zap.Error(err))
			return nil, err
		}
		results = append(results, article)
	}

	//查询总条数

	return nil, nil
}
