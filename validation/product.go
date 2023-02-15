package validation

type AddProduct struct {
	ProductName        string `json:"product_name" binding:"required,min=2,max=100"`
	ProductDescription string `json:"product_description" binding:"required,min=2,max=300"`
	Photo              string `json:"photo"`
}

type UpdateProduct struct {
	ProductName        string `json:"product_name" binding:"required,min=2,max=100"`
	ProductDescription string `json:"product_description" binding:"required,min=2,max=300"`
	Photo              string `json:"photo"`
}
