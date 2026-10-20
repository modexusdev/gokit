package validate

type Error struct {
	Field   string
	Message string
}

func (e *Error) Error() string {
	return e.Field + " " + e.Message
}

func newError(
	field string,
	message string,
) error {
	return &Error{
		Field:   field,
		Message: message,
	}
}
