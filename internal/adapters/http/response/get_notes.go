package response

import "github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"

type GetNotesResponse struct {
	Notes []entity.Note `json:"notes"`
}
