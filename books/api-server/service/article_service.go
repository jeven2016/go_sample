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
	"math"
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
		artSrv.log.Warn("Cannot findById(id)", zap.String("id", id), zap.Error(err))
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, common.NotFound
		}
		return nil, err
	}
	err = result.Decode(articleEntity)
	if err != nil {
		artSrv.log.Warn("error occurs when Decode(articleEntity)", zap.String("id", id), zap.Error(err))
	}
	return articleEntity, err
}

func (artSrv *ArticleService) List(catalogId string, pageRequest *dto.PageRequest) (*dto.ArticlePageResponse, error) {
	var results []*entity.Article
	findOpt := options.Find()
	findOpt.SetLimit(int64(pageRequest.PageSize))
	findOpt.SetSkip(int64((pageRequest.Page - 1) * pageRequest.PageSize))
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
	countOptions := &options.CountOptions{}
	count, err := artSrv.article.CountDocuments(context.TODO(), bson.M{"catalogId": catalogId}, countOptions)
	if err != nil {
		artSrv.log.Warn("An error occurs while counting articles", zap.Error(err))
		return nil, err
	}

	resp := &dto.ArticlePageResponse{
		Page:         pageRequest.Page,
		TotalPage:    int32(math.Ceil(float64(count) / float64(pageRequest.PageSize))),
		PageSize:     pageRequest.PageSize,
		TotalRecords: int32(count),
		Result: dto.Result{
			Payload: results,
		},
	}

	return resp, nil
}
