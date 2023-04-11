package response

import (
	"encoding/json"
	"product-crud/models"
	"time"

	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

type GetProductResponse struct {
	ID                 uint             `json:"product_id"`
	ProductName        string           `json:"product_name"`
	ProductDescription string           `json:"product_description"`
	Photo              string           `json:"photo"`
	UploaderId         uint             `json:"uploader_id"`
	Uploader           *GetUserResponse `json:"uploader,omitempty"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	DeletedAt          gorm.DeletedAt   `json:"deleted_at"`
}

func (res GetProductResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res GetProductResponse) IsEmpty() bool {
	return cmp.Equal(res, GetProductResponse{})
}

func NewGetProductResponse(p models.Product) GetProductResponse {
	var productUploader *GetUserResponse
	if p.Uploader != nil {
		user := *p.Uploader
		productUploader = NewGetUserResponse(user).Pointer()
	}
	return GetProductResponse{
		ID:                 p.ID,
		ProductName:        p.ProductName,
		ProductDescription: p.ProductDescription,
		Photo:              p.Photo,
		UploaderId:         p.UploaderId,
		Uploader:           productUploader,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
		DeletedAt:          p.DeletedAt,
	}
}
