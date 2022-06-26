package entity

import "time"

type Article struct {
	Id         string    `bson:"_id" json:"id"`
	Name       string    `bson:"name" json:"name"`
	CatalogId  string    `bson:"catalogId" json:"catalogId"`
	Content    string    `bson:"content,omitempty" json:"content"`
	CreateDate time.Time `bson:"createDate" json:"createDate"`
}
