package models

import (
	"math/rand/v2"
	"time"
)

type (
	AccessCode struct {
		ID       int       `json:"id"`
		Code     string    `json:"code"`
		IsUsed   *bool     `json:"is_used"`
		UserID   int       `json:"user_id"`
		ExpireAt time.Time `json:"expire_at"`
	}
)

func (a *AccessCode) Generate(userID int) {
	a.UserID = userID
	notUsed := false
	a.IsUsed = &notUsed
	a.Code = a.generateAccessCode()
	a.ExpireAt = time.Now().Add(time.Minute * 15)
}

func (a *AccessCode) generateAccessCode() string {
	var possible_chars = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var elected string
	for range 4 {
		pos := rand.IntN(len(possible_chars))
		elected += possible_chars[pos]
	}
	return elected
}
