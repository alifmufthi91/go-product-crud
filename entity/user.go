package entity

type User struct {
	UserId          string  `json:"user_id" db:"user_id"`
	Name            string  `json:"name" db:"name"`
	Email           string  `json:"email" db:"email"`
	PhoneNumberCode string  `json:"phone_number_code" db:"phone_number_code"`
	PhoneNumber     string  `json:"phone_number" db:"phone_number"`
	PhotoUrl        *string `json:"photo_url,omitempty" db:"photo_url"`
	Gender          string  `json:"gender" db:"gender"`
	AlgoAddress     string  `json:"algo_address" db:"algo_address"`
	Status          string  `json:"status" db:"status"`
	CreatedAt       int64   `json:"created_at" db:"created_at"`
	UpdatedAt       *int64  `json:"updated_at,omitempty" db:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
