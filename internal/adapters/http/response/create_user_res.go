package response

type CreateUserResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
