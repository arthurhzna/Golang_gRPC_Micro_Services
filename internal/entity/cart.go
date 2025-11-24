package entity

import "time"

type Cart struct {
	Id        string
	UserId    string
	ProductId string
	Quantity  int
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt *time.Time
	UpdatedBy *string
}
