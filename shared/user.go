package shared

// User is the user type
type User struct {
	ID                    string   `bson:"user_id" json:"id,omitempty"`
	Password              string   `bson:"password" json:"password,omitempty"`
	Name                  string   `bson:"name" json:"name,omitempty"`
	Age                   int      `bson:"age" json:"age,omitempty"`
	AdditionalInformation string   `bson:"additional_information" json:"additional_information,omitempty"`
	Parents               []string `bson:"parents" json:"parents,omitempty"`
}

// Parent is the parent type
type Parent struct {
	UserID string
	Name   string
}
