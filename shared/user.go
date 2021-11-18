package shared

// User is the user type
type User struct {
	ID                    string   `json:"id,omitempty"`
	Password              string   `json:"password,omitempty"`
	Name                  string   `json:"name,omitempty"`
	Age                   int      `json:"age,omitempty"`
	AdditionalInformation string   `json:"additional_information,omitempty"`
	Parents               []string `json:"parents,omitempty"`
}

// Parent is the parent type
type Parent struct {
	UserID string
	Name   string
}
