package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID        uuid.UUID `gorm:"primary_key;column:id" json:"id"`
	TeamName  string    `gorm:"size:150;not null;unique;column:team_name" json:"team_name"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (team *Team) BeforeCreate(tx *gorm.DB) (err error) {
	// tạo một UUID mới
	if team.ID == uuid.Nil {
		team.ID = uuid.New()
	}
	return
}

func (Team) TableName() string {
	return "Teams"
}

type Roster struct {
	TeamID uuid.UUID `gorm:"not null;column:team_id" json:"team_id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;column:user_id" json:"user_id"`
	Role   string    `gorm:"default:false;column:role" json:"role"`

	// Foreign key relationships
	Team Team `gorm:"foreignKey:TeamID;references:ID"`
}

// TableName overrides the table name used by Roster to `Rosters`
func (Roster) TableName() string {
	return "Rosters"
}
