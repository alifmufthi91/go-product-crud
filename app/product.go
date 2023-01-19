package app

import "encoding/json"

type Product struct {
	ID                 uint   `json:"product_id"`
	ProductName        string `json:"product_name"`
	ProductDescription string `json:"product_description"`
	Photo              string `json:"photo"`
	UploaderId         uint   `json:"uploader_id"`
	Uploader           User   `json:"uploader,omitempty"`
}

func (res Product) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}
