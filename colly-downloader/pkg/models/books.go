package models

import "time"

type CatalogDoc struct {
	Id           string    `bson:"_id" json:"id"`
	ParentId     string    `bson:"parentId" json:"parentId"`
	Name         string    `bson:"name" json:"name"`
	Order        int32     `bson:"order" json:"order"`
	ArticleCount int32     `bson:"articleCount" json:"articleCount"`
	Description  string    `bson:"description" json:"description"`
	CreateDate   time.Time `bson:"createDate" json:"createDate"`
	LastUpdate   time.Time `bson:"lastUpdate" json:"lastUpdate"`
}

type Article struct {
	Id         string    `bson:"_id" json:"id"`
	Name       string    `bson:"name" json:"name"`
	CatalogId  string    `bson:"catalogId" json:"catalogId"`
	Content    string    `bson:"content,omitempty" json:"content"`
	CreateDate time.Time `bson:"createDate" json:"createDate"`
}
