package repository

type User struct {
	ID                    string
	PasswordHash          string
	Name                  string
	Age                   int
	AdditionalInformation string
	Parents               []string
}

type Parent struct {
	UserID string
	Name   string
}
