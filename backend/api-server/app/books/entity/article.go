package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"` // omitempty会自动生成
	// Name       string             `bson:"name" json:"name"`
	CatalogId primitive.ObjectID `bson:"catalogId" json:"catalogId"`
	// Content    string             `bson:"content,omitempty" json:"content"`
	CreateDate time.Time `bson:"createDate" json:"createDate"`
}
