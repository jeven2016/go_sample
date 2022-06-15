package entity

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type BookCatalog struct {
	Id       string `bson:"_id" bson:"id"`
	ParentId string `bson:"parentId" json:"parentId"`
	//Name         string `bson:"name"`
	Order        int32 `bson:"order" json:"order"`
	ArticleCount int32 `bson:"articleCount" json:"articleCount"`
	//Description  string `bson:"description"`
	catalog *mongo.Collection
}
