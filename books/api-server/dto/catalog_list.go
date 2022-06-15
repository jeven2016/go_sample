package dto

import "api/entity"

type CatalogListResponse struct {
	Count int32
	List  []*entity.BookCatalog
}

type PageRequest struct {
	//Page     int32 `json:"page" validate:"required,gte=1,lte=1000000"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"pageSize" validate:"required,gte=10,lte=100"`
}
