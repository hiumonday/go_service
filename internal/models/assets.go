package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Folder struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name string    `gorm:"type:varchar(100);not null" json:"name"`

	OwnerID   uuid.UUID `gorm:"type:uuid;not null" json:"ownerId"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null" json:"teamId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Notes     []Note    `gorm:"foreignkey:FolderID"`
}

func (folder *Folder) BeforeCreate(tx *gorm.DB) (err error) {
	folder.ID = uuid.New()
	return
}

func (Folder) TableName() string {
	return "Folders"
}

type Note struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	FolderID  uuid.UUID `gorm:"type:uuid;not null" json:"folderId"`
	OwnerID   uuid.UUID `gorm:"type:uuid;not null" json:"ownerId"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null" json:"teamId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (note *Note) BeforeCreate(tx *gorm.DB) (err error) {
	note.ID = uuid.New()
	return
}

func (Note) TableName() string {
	return "Notes"
}

type Share struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	ResourceID   uuid.UUID `gorm:"type:uuid;not null" json:"resourceId"`
	ResourceType string    `gorm:"type:varchar(50);not null" json:"resourceType"` // "folder" or "note"
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	Permission   string    `gorm:"type:varchar(10);not null" json:"permission"` // "read" or "write"
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (share *Share) BeforeCreate(tx *gorm.DB) (err error) {
	share.ID = uuid.New()
	return
}

func (Share) TableName() string {
	return "Shares"
}
