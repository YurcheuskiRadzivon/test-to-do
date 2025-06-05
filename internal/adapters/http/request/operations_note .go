package request

type OperationNoteRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` //validate:"omitempty,oneof=SUCCESS IN_PROGRESS NOT_START"
}
