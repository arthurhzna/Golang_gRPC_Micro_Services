package entity

import "time"

type UserRole struct {
	Id        string
	Name      string
	Code      string
	CreatedAt time.Time
	CreatedBy *string
	UpdateAt  time.Time
	UpdatedBy *string
	DeletedAt *time.Time
	DeletedBy *string
	IsDelated bool
}

type User struct {
	Id        string
	FullName  string
	Email     string
	Password  string
	RoleCode  string
	CreatedAt time.Time
	CreatedBy *string
	UpdateAt  time.Time
	UpdatedBy *string
	DeletedAt *time.Time
	DeletedBy *string
	IsDelated bool
}
