package store

import "time"

type User struct {
	ID       uint64    `json:"id"`
	Name     string    `json:"name"`
	Family   string    `json:"family"`
	BirthDay time.Time `json:"birth_day"`
	Address  string    `json:"address"`
	Phone    []string  `json:"phone"`
}
