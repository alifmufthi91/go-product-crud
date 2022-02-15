package validation

type AddProduct struct {
	ProductName        string `json:"product_name" validate:"required,min=2,max=100"`
	ProductDescription string `json:"product_description" validate:"required,min=2,max=300"`
	Photo              string `json:"photo"`
}
