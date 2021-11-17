package shared

// User is the user type
type User struct {
	ID                    string
	Password              string
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
