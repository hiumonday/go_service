package dto

type CreateFolderRequest struct {
	Name string `json:"name" binding:"required"`

	TeamID string `json:"teamId" binding:"required"`
}

type UpdateFolderRequest struct {
	Name string `json:"name"`
}

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ShareRequest struct {
	UserID     string `json:"userId" binding:"required"`
	Permission string `json:"permission" binding:"required,oneof=read write"`
}
