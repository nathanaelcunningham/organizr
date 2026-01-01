package models

import "time"

type Config struct {
	Key         string
	Value       string
	Description string
	UpdatedAt   time.Time
}
