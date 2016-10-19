package store

import "time"

type User struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Family       string    `json:"family"`
	BirthDay     string    `json:"birth_day"`
	BirthDayFull time.Time `json:"birth_day_full"`
	Address      string    `json:"address"`
	Phone        []string  `json:"phone"`
}
