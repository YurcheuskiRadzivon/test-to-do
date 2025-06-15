package request

import "github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"

type GetNoteReq struct {
	Noteinfo   entity.Note `json:"noteinfo"`
	FileIDList []int       `json:"file_id_list"`
}
