package common

import "math"

// swagger:model PaginationMetadata
type PaginationMetadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

func CalculateMetadata(totalRecords, page, pageSize int) PaginationMetadata {
	if totalRecords == 0 {
		return PaginationMetadata{
			CurrentPage: page,
			PageSize:    pageSize,
			FirstPage:   1,
			LastPage:    1,
		}
	}

	return PaginationMetadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
