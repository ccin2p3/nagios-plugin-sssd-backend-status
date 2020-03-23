package errors

type CmdError struct {
	Err error
	Rc  int
}

func (e CmdError) Error() string {
	return e.Err.Error()
}

func NewCmdError(err error, rc int) CmdError {
	return CmdError{
		Err: err,
		Rc:  rc,
	}
}
