package models

import "time"

type Listing struct {
	Id          string
	Title       string
	Description string
	Price       int
	City        string
	CreatedAt   time.Time
}
