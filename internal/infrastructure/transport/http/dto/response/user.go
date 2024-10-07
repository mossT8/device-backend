package response

type User struct {
	AccountId       int64
	Email           string
	Cell            string
	FirstName       string
	LastName        string
	Verified        bool
	ReceivesUpdates bool
}
