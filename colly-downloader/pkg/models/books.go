package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatalogDoc struct {
	// 添加omitempty，当为空时，mongo driver会自动生成
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ParentId     primitive.ObjectID `bson:"parentId,omitempty" json:"parentId"`
	Name         string             `bson:"name" json:"name"`
	Order        int32              `bson:"order" json:"order"`
	ArticleCount int32              `bson:"articleCount" json:"articleCount"`
	Description  string             `bson:"description" json:"description"`
	CreateDate   time.Time          `bson:"createDate" json:"createDate"`
	LastUpdate   time.Time          `bson:"lastUpdate" json:"lastUpdate"`
}

type Article struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	CatalogId  primitive.ObjectID `bson:"catalogId,omitempty" json:"catalogId"`
	Content    string             `bson:"content,omitempty" json:"content"`
	CreateDate time.Time          `bson:"createDate" json:"createDate"`
}

type ArticlePage struct {
	Name string // 有可能是繁体中文
	Url  string
}
