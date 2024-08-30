package models

import "time"

type Categories struct {
	Id        int32
	Name      string
	Slug      string
	UpdatedAt time.Time
	CreatedAt time.Time
}
