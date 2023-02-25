package response

import (
	"encoding/json"
	"product-crud/models"

	"github.com/google/go-cmp/cmp"
)

type GetUserResponse struct {
	ID        uint                 `json:"user_id"`
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Email     string               `json:"email"`
	Products  []GetProductResponse `json:"products,omitempty"`
}

func (res GetUserResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res GetUserResponse) IsEmpty() bool {
	return cmp.Equal(res, GetUserResponse{})
}

func NewGetUserResponse(u models.User) GetUserResponse {
	var productDatas []GetProductResponse
	for _, product := range u.Products {
		productDatas = append(productDatas, NewGetProductResponse(product))
	}
	return GetUserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Products:  productDatas,
	}
}
