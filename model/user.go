package model

type UserRegisterInput struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Gender          string `json:"gender"`
	PhoneNumberCode string `json:"phone_number_code"`
	PhotoUrl        string `json:"photo_url"`
	AlgoAddress     string `json:"algo_address"`
}
