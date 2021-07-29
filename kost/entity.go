package kost

import (
	"backendEkost/user"
	"time"
)

type Kost struct {
	ID                int
	UserID            int
	Name              string
	ShortDescription  string
	Description       string
	Perks             string
	Latitude          string
	Longitude         string
	Price             string
	CurrentSpaceCount int
	Slug              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	KostImages        []KostImage
	User              user.User
}

type KostImage struct {
	ID        int
	KostID    int
	FileName  string
	IsPrimary int
	CreatedAt time.Time
	UpdatedAt time.Time
}
