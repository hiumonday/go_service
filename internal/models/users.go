package models

type User struct {
	ID       string ` json:"userid"`
	Username string ` json:"username"`
	Email    string ` json:"email"`
	Role     string ` json:"role"`
}
