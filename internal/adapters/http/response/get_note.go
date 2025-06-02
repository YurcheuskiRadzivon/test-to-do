package response

import "github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"

type GetNoteResponse struct {
	Note entity.Note `json:"note"`
}
