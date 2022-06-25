package dto

import "api/entity"

type CatalogListResponse struct {
	Count int32
	List  []*entity.BookCatalog
}

type PageRequest struct {
	Page int32 `form:"page"  binding:"gte=1,lte=1000000"`

	//因为是使用Query参数查询方式，格式上兼容form，所以需要添加form的方式
	PageSize int32  `form:"pageSize"   binding:"gte=10,lte=100"`
	Search   string `form:"search" json:"search"`
}
