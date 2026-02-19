package domain

import "time"

type (
	UserID string

	User struct {
		ID        UserID
		CreatedAt time.Time
	}
)
