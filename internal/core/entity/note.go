package entity

type Note struct {
	NoteID      string `json:"note_id"`
	Tittle      string `json:"tittle"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
