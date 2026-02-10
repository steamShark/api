package model_pagination

/*
This is  the object requested in the params
*/
type Pagination struct {
	Page     int
	PageSize int
}

/*
This is the object to be returned
*/
type PaginationMeta struct {
	Page     int
	PageSize int
	Total    *int64 `json:"total,omitempty"`
}

// Generic list wrapper
type PaginatedListResult[T any] struct {
	Items    []T
	Total    int64
	Page     int
	PageSize int
}
