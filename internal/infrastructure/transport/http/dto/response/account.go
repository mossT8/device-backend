package response

import "time"

type Account struct {
	Email           string
	Name            string
	ReceivesUpdates bool
	CreatedAt       time.Time
	ModifiedAt      time.Time
}
