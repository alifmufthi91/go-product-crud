package entity

type User struct {
	UserId          string `json:"user_id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	PhoneNumberCode string `json:"phone_number_code"`
	PhoneNumber     string `json:"phone_number"`
	PhotoUrl        string `json:"photo_url,omitempty"`
	Gender          string `json:"gender"`
	AlgoAddress     string `json:"algo_address"`
	Status          string `json:"status"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at,omitempty"`
}
