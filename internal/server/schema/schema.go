package schema

const (
	ErrInternal   = "Internal server error"
	ErrNotFound   = "Not found"
	ErrBadRequest = "Bad request"
	ErrTimedOut   = "Timed out"
)

type ServerError struct {
	Error string `json:"error"`
}
