package schema

const (
	ErrInternal   = "Internal server error"
	ErrBadRequest = "Bad request"
	ErrTimedOut   = "Timed out"
)

type ServerError struct {
	Description string `json:"description"`
}
