package users

const (
	insert = `
		INSERT INTO users (username, password) VALUES (:username, :password)
	`
	findByUsername = `
		"SELECT id, username, password, role FROM users WHERE username = ?";
	`

	delete     = `DELETE FROM users WHERE id = :id`
	updateRole = `UPDATE users SET role = :role WHERE id = :id`
)
