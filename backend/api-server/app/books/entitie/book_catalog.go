package entitie

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type BookCatalog struct {
	Id       string `bson:"_id" json:"id"`
	ParentId string `bson:"parentId" json:"parentId"`
	// Name         string `bson:"name" json:"name"`
	Order        int32 `bson:"order" json:"order"`
	ArticleCount int32 `bson:"articleCount" json:"articleCount"`
	// Description  string `bson:"description" json:"description"`
	catalog *mongo.Collection
}
