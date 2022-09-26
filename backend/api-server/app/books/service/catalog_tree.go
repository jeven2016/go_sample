package service

import (
	"github.com/duke-git/lancet/v2/slice"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"api/app/books/dto"
	"api/app/books/entity"
)

var blankObjectId = [12]byte{}

// CatalogTree Generate a complete tree by catalog entity array
func CatalogTree(catalogs []*entity.BookCatalog) *dto.CatalogListResponse {
	if len(catalogs) == 0 {
		return &dto.CatalogListResponse{}
	}

	var rootResponse = &dto.CatalogListResponse{List: []*dto.CatalogResponse{}}

	// 使用mongo的objectId做key
	existingMap := make(map[primitive.ObjectID]*dto.CatalogResponse, len(catalogs))

	var rootNodes = getRootNodes(catalogs)
	addRootTreeNodes(rootResponse, rootNodes, existingMap)
	addChildren(catalogs, rootNodes, existingMap)

	return rootResponse
}

func addChildren(catalogs []*entity.BookCatalog, nodes []*entity.BookCatalog, existingMap map[primitive.ObjectID]*dto.CatalogResponse) {
	var nodeCount = 0
	var total = len(catalogs)

	if nodeCount < total {
		for _, c := range catalogs {
			if _, exists := existingMap[c.Id]; exists {
				// it has been processed, just .ignore
				nodeCount++
				continue
			}
			if parentNode, parentExists := existingMap[c.ParentId]; parentExists {
				parentNode.Children = append(parentNode.Children, toCatalogResponse(c))
				sort(parentNode.Children)
				nodeCount++
			}
		}
	}
}

func addRootTreeNodes(root *dto.CatalogListResponse, rootNodes []*entity.BookCatalog,
	existingMap map[primitive.ObjectID]*dto.CatalogResponse) {

	for _, node := range rootNodes {
		var treeNode = toCatalogResponse(node)
		existingMap[node.Id] = treeNode
		root.List = append(root.List, treeNode)
	}
	sort(root.List)
}

func sort(array []*dto.CatalogResponse) {
	err := slice.SortByField(array, "Order", "asc")
	if err != nil {
		panic(err)
	}
}

func getRootNodes(catalogs []*entity.BookCatalog) []*entity.BookCatalog {
	return slice.Filter(catalogs, func(index int, item *entity.BookCatalog) bool {
		if item.ParentId == blankObjectId {
			return true
		}
		return false
	})
}

func toCatalogResponse(node *entity.BookCatalog) *dto.CatalogResponse {
	var treeNode = &dto.CatalogResponse{
		Id:           node.Id.Hex(),
		ParentId:     node.ParentId.Hex(),
		Name:         node.Name,
		Order:        node.Order,
		ArticleCount: node.ArticleCount,
		Description:  node.Description,
		Children:     []*dto.CatalogResponse{},
	}
	return treeNode
}
