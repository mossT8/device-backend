package response

import "time"

type User struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	Cell            string    `json:"cell"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Verified        bool      `json:"verified"`
	ReceivesUpdates bool      `json:"receivesUpdates"`
	CreatedAt       time.Time `json:"createdAt"`
	ModifiedAt      time.Time `json:"modifiedAt"`
}
