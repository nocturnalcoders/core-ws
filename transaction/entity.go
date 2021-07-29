package transaction

import (
	"backendEkost/kost"
	"backendEkost/user"
	"time"
)

type Transaction struct {
	ID         int
	KostID     int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User       user.User
	Kost       kost.Kost
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
