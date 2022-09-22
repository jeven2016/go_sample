package yzs8

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"colly-downloader/pkg/models"
)

// converting

// ConvertCatalog 将ID、 parentId从string转为objectId
func ConvertCatalog() {
	// mongo
	client, err := CreateMongoClient(log)
	handleError(err)

	database := client.Database("books")
	catalogCol := database.Collection("catalog")

	cur, err := catalogCol.Find(context.Background(), bson.M{})

	var results = new([]models.CatalogDocMap)
	err = cur.All(context.Background(), results)
	handleError(err)

	for _, c := range *results {
		if reflect.TypeOf(c.Id) == reflect.TypeOf(string("")) {
			// 转换ID
			stringId := c.Id.(string)
			objectId, err := primitive.ObjectIDFromHex(stringId)
			handleError(err)

			// mongodb无法更新_id, 只能新建一个同时删除旧的
			if err == nil {
				// 插入新的
				c.Id = objectId
				one, err := catalogCol.InsertOne(context.Background(), c)
				handleError(err)

				if err == nil {
					newId := one.InsertedID.(primitive.ObjectID).Hex()
					println("new Id=", newId)

					// 删除旧文档
					_, err := catalogCol.DeleteOne(context.Background(), bson.M{"_id": stringId})
					if err != nil {
						panic(err)
					}
				}

			}

		}

		// 转换parentId
		if reflect.TypeOf(c.ParentId) == reflect.TypeOf(string("")) {

		}
	}

	print(results)
}
