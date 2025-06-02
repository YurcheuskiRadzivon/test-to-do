package entity

type Note struct {
	NoteID      int    `json:"note_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
