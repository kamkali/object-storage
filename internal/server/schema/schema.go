package schema

const (
	ErrInternal   = "Internal server error"
	ErrNotFound   = "Not found"
	ErrBadRequest = "Bad request"
	ErrInvalidID  = "Invalid id - must be up to 32 alphanum"
	ErrTimedOut   = "Timed out"
)

type ServerError struct {
	Error string `json:"error"`
}

type PutObjectResponse struct {
	ID string
}
