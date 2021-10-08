package digcopy

// ======================================================
// copy and modify from https://github.com/uber-go/dig/blob/v1.13.0/error.go
// ======================================================

import (
	"fmt"
	"io"
)

type causer interface {
	fmt.Formatter

	// Returns the next error in the chain.
	cause() error

	// Writes the message or context for this error in the chain.
	//
	// verb is either %v or %+v.
	writeMessage(w io.Writer, verb string)
}

func formatCauser(c causer, w fmt.State, v rune) {
	multiline := w.Flag('+') && v == 'v'
	verb := "%v"
	if multiline {
		verb = "%+v"
	}

	// "context: " or "context:\n"
	c.writeMessage(w, verb)
	io.WriteString(w, ":")
	if multiline {
		io.WriteString(w, "\n")
	} else {
		io.WriteString(w, " ")
	}

	fmt.Fprintf(w, verb, c.cause())
}

// ErrProvide is returned when a constructor could not be Provided into the
// container.
type ErrProvide struct {
	Func   *Func
	Reason error
}

var _ causer = ErrProvide{}

func (e ErrProvide) cause() error {
	return e.Reason
}

func (e ErrProvide) writeMessage(w io.Writer, verb string) {
	fmt.Fprintf(w, "cannot provide function "+verb, e.Func)
}

func (e ErrProvide) Error() string { return fmt.Sprint(e) }
func (e ErrProvide) Format(w fmt.State, c rune) {
	formatCauser(e, w, c)
}

// ErrConstructorFailed is returned when a user-provided constructor failed
// with a non-nil error.
type ErrConstructorFailed struct {
	Func   *Func
	Reason error
}

var _ causer = ErrConstructorFailed{}

func (e ErrConstructorFailed) cause() error {
	return e.Reason
}

func (e ErrConstructorFailed) writeMessage(w io.Writer, verb string) {
	fmt.Fprintf(w, "received non-nil error from function "+verb, e.Func)
}

func (e ErrConstructorFailed) Error() string { return fmt.Sprint(e) }
func (e ErrConstructorFailed) Format(w fmt.State, c rune) {
	formatCauser(e, w, c)
}

// ErrArgumentsFailed is returned when a function could not be run because one
// of its dependencies failed to build for any reason.
type ErrArgumentsFailed struct {
	Func   *Func
	Reason error
}

var _ causer = ErrArgumentsFailed{}

func (e ErrArgumentsFailed) cause() error {
	return e.Reason
}

func (e ErrArgumentsFailed) writeMessage(w io.Writer, verb string) {
	fmt.Fprintf(w, "could not build arguments for function "+verb, e.Func)
}

func (e ErrArgumentsFailed) Error() string { return fmt.Sprint(e) }
func (e ErrArgumentsFailed) Format(w fmt.State, c rune) {
	formatCauser(e, w, c)
}

// ErrMissingDependencies is returned when the dependencies of a function are
// not available in the container.
type ErrMissingDependencies struct {
	Func   *Func
	Reason error
}

var _ causer = ErrMissingDependencies{}

func (e ErrMissingDependencies) cause() error {
	return e.Reason
}

func (e ErrMissingDependencies) writeMessage(w io.Writer, verb string) {
	fmt.Fprintf(w, "missing dependencies for function "+verb, e.Func)
}

func (e ErrMissingDependencies) Error() string { return fmt.Sprint(e) }
func (e ErrMissingDependencies) Format(w fmt.State, c rune) {
	formatCauser(e, w, c)
}
