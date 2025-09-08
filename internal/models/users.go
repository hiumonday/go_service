package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID ` json:"userid"`
	Username string    ` json:"username"`
	Email    string    ` json:"email"`
	Role     string    ` json:"role"`
}
