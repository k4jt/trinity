package store

import (
	"strings"
	"time"
)

type User struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Family       string    `json:"family"`
	BirthDay     string    `json:"birth_day"`
	BirthDayFull time.Time `json:"birth_day_full"`
	Address      string    `json:"address"`
	Phone        []string  `json:"phone"`
}

func (u *User) Contains(q string) bool {

	if strings.Contains(strings.ToLower(u.Name), q) {
		return true
	}
	if strings.Contains(strings.ToLower(u.Family), q) {
		return true
	}
	if strings.Contains(strings.ToLower(u.BirthDay), q) {
		return true
	}
	if strings.Contains(strings.ToLower(u.Address), q) {
		return true
	}
	if strings.Contains(strings.ToLower(strings.Join(u.Phone, ", ")), q) {
		return true
	}

	return false
}
