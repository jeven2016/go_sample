package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookCatalog struct {
	Id           primitive.ObjectID `bson:"_id" json:"id"`
	ParentId     primitive.ObjectID `bson:"parentId" json:"parentId"`
	Name         string             `bson:"name" json:"name"`
	Order        int32              `bson:"order" json:"order"`
	ArticleCount int32              `bson:"articleCount" json:"articleCount"`
	Description  string             `bson:"description" json:"description"`
}
