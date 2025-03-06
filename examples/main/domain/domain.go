// Package domain contains business logic models.
package domain

// Account represents a user's account.
type Account struct {
	ID       int
	Username string `dbgo:"users.name"`
	Password string `dbgo:"users.password"`
	Name     string `dbgo:"accounts.name"`
}
