package entity

const (
	OwnerNote OwnerType = "NOTE"
)

type OwnerType string

type FileMeta struct {
	FileID      int
	ContentType string
	OwnerType   OwnerType
	OwnerID     int
	UserID      int
	URI         string
}
