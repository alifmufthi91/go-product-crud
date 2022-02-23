package app

type User struct {
	ID        uint      `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Products  []Product `json:"products,omitempty"`
}
