package model

type PaginationParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
