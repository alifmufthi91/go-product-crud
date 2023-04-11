package response

import (
	"encoding/json"
	"product-crud/models"
	"time"

	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

type GetUserResponse struct {
	ID        uint                 `json:"user_id"`
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Email     string               `json:"email"`
	Products  []GetProductResponse `json:"products"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	DeletedAt gorm.DeletedAt       `json:"deleted_at"`
}

func (res GetUserResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res GetUserResponse) IsEmpty() bool {
	return cmp.Equal(res, GetUserResponse{})
}

func (res GetUserResponse) Pointer() *GetUserResponse {
	return &res
}

func NewGetUserResponse(u models.User) GetUserResponse {
	productDatas := []GetProductResponse{}
	for _, product := range u.Products {
		productDatas = append(productDatas, NewGetProductResponse(product))
	}
	return GetUserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Products:  productDatas,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}
