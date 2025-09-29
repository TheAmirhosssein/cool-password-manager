package errors

func NewServerError() error {
	return NewError("Internal Server Error", 500_00)
}
