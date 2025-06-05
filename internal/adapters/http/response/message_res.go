package response

const (
	MessageSuccsessfully = "SUCCESFULLY"
)

type MessageResponse struct {
	Message string `json:"message"`
}
