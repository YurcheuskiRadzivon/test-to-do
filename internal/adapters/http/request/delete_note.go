package request

type DeleteNoteRequest struct {
	NoteID int `json:"id"` //validate:"required,min=1
}
