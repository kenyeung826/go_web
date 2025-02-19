package error

type ApplicationError struct {
	s string
}

func (ae ApplicationError) Error() string {
	return ae.s
}

func NewApplicationError(text string) *ApplicationError {
	return &ApplicationError{
		s: text,
	}
}
