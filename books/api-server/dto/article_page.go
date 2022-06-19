package dto

type ArticlePageResponse struct {
	Page         int32
	TotalPage    int32
	PageSize     int32
	TotalRecords int32
	Result
}
