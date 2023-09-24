package models

type AppError struct {
	Message string
	Err   error
}

func (e *AppError) Error() string {
    return e.Message
}
