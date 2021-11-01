package entity

type User struct {
	UserId      string `json:"user_id" db:"user_id" gorm:"primaryKey"`
	AlgoAddress string `json:"algo_address" db:"algo_address"`
	Status      string `json:"status" db:"status"`
	CreatedAt   int64  `json:"created_at" db:"created_at"`
	UpdatedAt   *int64 `json:"updated_at,omitempty" db:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
