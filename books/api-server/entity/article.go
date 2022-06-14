package entity

import "time"

type Article struct {
	Id string `bson:"_id" json:"id"`
	//Name       string        `bson:"name"`
	CatalogId string `bson:"catalogId" json:"catalogId"`
	//content    string        `bson:"content,omitempty"`
	CreateDate time.Time `bson:"createDate" json:"createDate"`
}
