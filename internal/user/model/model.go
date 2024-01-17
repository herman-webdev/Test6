package model

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type User struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Patronymic string    `json:"patronymic"`
	Age        int       `json:"age"`
	Gender     string    `json:"gender"`
	CountryId  string    `json:"country_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  null.Time `json:"updated_at"`
}
