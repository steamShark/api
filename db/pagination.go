package db

// Generic list wrapper
type PaginatedListResult[T any] struct {
	Items  []T
	Count  int64
	Limit  int
	Offset int
}
