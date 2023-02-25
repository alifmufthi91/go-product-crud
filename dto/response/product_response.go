package response

import (
	"encoding/json"
	"product-crud/models"

	"github.com/google/go-cmp/cmp"
)

type GetProductResponse struct {
	ID                 uint            `json:"product_id"`
	ProductName        string          `json:"product_name"`
	ProductDescription string          `json:"product_description"`
	Photo              string          `json:"photo"`
	UploaderId         uint            `json:"uploader_id"`
	Uploader           GetUserResponse `json:"uploader,omitempty"`
}

func (res GetProductResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res GetProductResponse) IsEmpty() bool {
	return cmp.Equal(res, GetProductResponse{})
}

func NewGetProductResponse(p models.Product) GetProductResponse {
	return GetProductResponse{
		ID:                 p.ID,
		ProductName:        p.ProductName,
		ProductDescription: p.ProductDescription,
		Photo:              p.Photo,
		UploaderId:         p.UploaderId,
		Uploader:           NewGetUserResponse(p.Uploader),
	}
}
