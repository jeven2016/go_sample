package entity

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type BookCatalog struct {
	ParentId     string `bson:"parentId"`
	Name         string `bson:"name"`
	Order        int32  `bson:"order"`
	ArticleCount int32  `bson:"articleCount"`
	Description  string `bson:"description"`
	catalog      *mongo.Collection
}
