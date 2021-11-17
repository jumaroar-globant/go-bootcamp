package repository

// User is the user type
type User struct {
	ID                    string
	PasswordHash          string
	Name                  string
	Age                   int
	AdditionalInformation string
	Parents               []string
}

// Parent is the parent type
type Parent struct {
	UserID string
	Name   string
}
