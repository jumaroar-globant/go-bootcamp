package repository

const (
	// PasswordHashQuery is a SQL query to obtain a password hash
	PasswordHashQuery string = "SELECT password_hash FROM users WHERE name=?"
	// InsertUserStatement is a SQL statement to insert a user
	InsertUserStatement string = "INSERT INTO users (id, name, password_hash, age, additional_information) VALUES(?, ?, ?, ?, ?)"
	// InsertParentStatement is an SQL statement to insert a parent
	InsertParentStatement string = "INSERT INTO user_parents (user_id, name) VALUES(?, ?)"
	// UpdateUserStatement is an SQL statement to update a user
	UpdateUserStatement string = "UPDATE users SET name=?, age=?, additional_information=?  WHERE id = ?"
	// DeleteUserParentsStatement is an SQL statement to delete a user parents
	DeleteUserParentsStatement string = "DELETE FROM user_parents WHERE user_id=?"
	// UserDataQuery is a SQL query to obtain a user data
	UserDataQuery string = "SELECT id, name, age, additional_information FROM users WHERE id=?"
	// UserParentsQuery is a SQL query to obtain a user parents
	UserParentsQuery string = "SELECT name FROM user_parents WHERE user_id=?"
	//DeleteUserStatement is a SQL statement to delete a user
	DeleteUserStatement string = "DELETE FROM users WHERE id=?"

	usernameKey = "users.users_UN"
)
