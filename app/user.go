package app

import "encoding/json"

type User struct {
	ID        uint      `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Products  []Product `json:"products,omitempty"`
}

func (res User) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}
