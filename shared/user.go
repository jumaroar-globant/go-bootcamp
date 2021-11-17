package shared

type User struct {
	ID                    string   `json:"id,omitempty"`
	Password              string   `json:"password,omitempty"`
	Name                  string   `json:"name"`
	Age                   int      `json:"age,omitempty"`
	AdditionalInformation string   `json:"additional_information,omitempty"`
	Parents               []string `json:"parents,omitempty"`
}

type Parent struct {
	UserID string
	Name   string
}
