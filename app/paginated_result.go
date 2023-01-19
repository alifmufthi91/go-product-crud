package app

import "encoding/json"

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
