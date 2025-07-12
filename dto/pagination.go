package dto

type Pagination struct {
	Page       int                   `json:"page" example:"1"`
	Size       int                   `json:"size" example:"10"`
	TotalCount int                   `json:"total_count" example:"100"`
	TotalPages int                   `json:"total_pages" example:"10"`
	List       []*ResponseGetUserAll `json:"list"`
}
