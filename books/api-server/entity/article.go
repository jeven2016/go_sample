package entity

import "time"

type Article struct {
	Id         string        `bson:"_id"`
	Name       string        `bson:"name"`
	CatalogId  string        `bson:"catalogId"`
	content    string        `bson:"content"`
	createDate time.Duration `bson:"createDate"`
}
