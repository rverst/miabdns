package models

type Error struct {
	Status   int
	Err      error
	ToClient bool
}

func NewError(status int, err error, toClient bool) *Error {
	return &Error{
		Status:   status,
		Err:      err,
		ToClient: toClient,
	}
}
