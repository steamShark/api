package models

import (
	model_filters "steamshark-api/internal/models/filters"
	model_pagination "steamshark-api/internal/models/pagination"
	model_response "steamshark-api/internal/models/response"
	model_statistics "steamshark-api/internal/models/statistics"
	model_website "steamshark-api/internal/models/websites"
)

type (
	/* WEBISTE related */
	Website    = model_website.Website
	Occurrence = model_website.Occurrence

	/* Filters */
	ListWebsitesFilter          = model_filters.ListWebsitesFilter
	ListWebsitesExtensionFilter = model_filters.ListWebsitesExtensionFilter

	/* PAGINATION */
	Pagination                 = model_pagination.Pagination
	PaginatedListResult[T any] = model_pagination.PaginatedListResult[T]
	PaginationMeta             = model_pagination.PaginationMeta

	/* STATISTICS */
	Statistics = model_statistics.Statistics

	/* RESPONSE models */
	APIResponse[T any] = model_response.ApiResponse[T]
)
