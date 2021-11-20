package digcopy

// ======================================================
// copy and modify from https://github.com/uber-go/dig/blob/v1.13.0/error.go
// ======================================================

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"go.uber.org/dig"
)

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

type Key struct {
	T reflect.Type

	// Only one of name or group will be set.
	Name  string
	Group string
}

func (po *Key) String() string {
	if po.Name == "" && po.Group == "" {
		return po.T.String()
	} else {
		s := []string{}
		if po.Name != "" {
			s = append(s, fmt.Sprintf("name=\"%s\"", po.Name))
		}
		if po.Group != "" {
			s = append(s, fmt.Sprintf("group=\"%s\"", po.Group))
		}
		return fmt.Sprintf("%s[%s]", po.T.String(), strings.Join(s, ","))
	}
}

// ErrParamSingleFailed is returned when a paramSingle could not be built.
type ErrParamSingleFailed struct {
	Key    Key
	Reason error
}

var _ causer = ErrParamSingleFailed{}

func (e ErrParamSingleFailed) cause() error {
	return e.Reason
}

func (e ErrParamSingleFailed) writeMessage(w io.Writer, _ string) {
	fmt.Fprintf(w, "failed to build %s", e.Key.String())
}

func (e ErrParamSingleFailed) Error() string { return fmt.Sprint(e) }
func (e ErrParamSingleFailed) Format(w fmt.State, c rune) {
	formatCauser(e, w, c)
}

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

func RootCause(err error) error {
	for {
		if e, ok := err.(causer); ok {
			err = e.cause()
		} else {
			return dig.RootCause(err)
		}
	}
}
