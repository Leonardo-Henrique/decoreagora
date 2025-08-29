package models

import "time"

type (
	AccessCode struct {
		ID       int       `json:"id"`
		Code     string    `json:"code"`
		IsUsed   bool      `json:"is_used"`
		UserID   int       `json:"user_id"`
		ExpireAt time.Time `json:"expire_at"`
	}
)
