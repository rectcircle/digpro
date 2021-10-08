package digpro

type err struct {
	msg string
}

func wrapError(inner error) error {
	if inner == nil {
		return nil
	}
	return &err{
		msg: inner.Error(),
	}
}

func (e *err) Error() string {
	return e.msg
}
