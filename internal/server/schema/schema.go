package schema

const (
	ErrInternal   = "Internal server error"
	ErrNotFound   = "Not found"
	ErrBadRequest = "Bad request"
	ErrInvalidID  = "ID must be alphanum up to 32 characters"
	ErrTimedOut   = "Timed out"
)

type ServerError struct {
	Error string `json:"error"`
}

type PutObjectResponse struct {
	ID string
}
