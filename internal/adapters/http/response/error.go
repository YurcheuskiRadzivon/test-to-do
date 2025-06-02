package response

const (
	ErrInvalidRequest  = "INVALID_REQUEST"
	ErrCodeInvalid     = "INVALID_CODE"
	ErrPassHashInvalid = "PASSWORD_HASH_INVALID"
	ErrSignInFailed    = "SIGN_IN_FAILED"
	ErrUnknown         = "UKNOWN_ERROR"
	ErrJWT             = "ENCRYPTION_ERROR"
)

type Error struct {
	Error string `json:"error" example:"message"`
}
