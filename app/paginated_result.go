package app

import (
	"encoding/json"

	"github.com/google/go-cmp/cmp"
)

type PaginatedResult struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	TotalItems int         `json:"total_items"`
	TotalPage  int         `json:"total_pages"`
}

func (res PaginatedResult) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res PaginatedResult) IsEmpty() bool {
	return cmp.Equal(res, PaginatedResult{})
}
