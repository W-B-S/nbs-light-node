package errors

const CommonError  = -10001


type NbsError struct {
	errno	int
	what 	string
	Err  	error
}

func New(text string) error {
	return &NbsError{errno:CommonError, what:text}
}

func New2(errno int, text string) error {
	return &NbsError{errno:errno, what:text}
}

func (err *NbsError) Error() string {
	return err.what
}

func (err *NbsError) ErrorNo() int {
	return err.errno
}