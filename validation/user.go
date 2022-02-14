package validation

type RegisterUser struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=10"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
