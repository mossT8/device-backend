package response

import "time"

type Account struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	ReceivesUpdates bool      `json:"receivesUpdates"`
	CreatedAt       time.Time `json:"createdAt"`
	ModifiedAt      time.Time `json:"modifiedAt"`
}
