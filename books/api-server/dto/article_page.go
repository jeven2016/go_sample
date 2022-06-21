package dto

type ArticlePageResponse struct {
	Page         int32 `json:"page"`
	TotalPage    int32 `json:"totalPage"`
	PageSize     int32 `json:"pageSize"`
	TotalRecords int32 `json:"totalRecords"`
	Result
}
