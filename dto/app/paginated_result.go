package app

import (
	"encoding/json"

	"github.com/google/go-cmp/cmp"
)

type PaginatedResult[i any] struct {
	Items      []i `json:"items"`
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPage  int `json:"total_pages"`
}

func (res PaginatedResult[i]) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res PaginatedResult[i]) IsEmpty() bool {
	return cmp.Equal(res, PaginatedResult[i]{})
}
