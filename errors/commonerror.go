package errors

const CommonError  = -10001

type NbsError interface {
	ErrorNo()	int
	Error() 	string
}

func New(text string) error {
	return &nbsError{errno:CommonError, what:text}
}

func New2(errno int, text string) error {
	return &nbsError{errno:errno, what:text}
}

type nbsError struct {
	errno	int
	what 	string
}

func (err *nbsError) Error() string {
	return err.what
}

func (err *nbsError) ErrorNo() int {
	return err.errno
}