package digpro

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	errs := []error{
		c.Supply("a"),
		c.Supply(1),
		c.Supply(true),
		c.Struct(new(Bar)),
	}
	for _, err := range errs {
		if err != nil {
			t.Errorf("c.Value() / c.Struct() error = %v, wantErr %v", err, false)
		}
	}
	err := c.Invoke(func(bar *Bar) {
		if bar.A != "a" {
			t.Errorf("bar.A got = %s, want = %s", bar.A, "a")
		}
		if bar.B != 1 {
			t.Errorf("bar.A got = %d, want = %d", bar.B, 1)
		}
		if bar.private != true {
			t.Errorf("bar.A got = %t, want = %t", bar.private, true)
		}
	})
	if err != nil {
		t.Errorf("c.Invoke() error = %v, wantErr %v", err, false)
	}
}

func TestContainerWrapper_Unwrap(t *testing.T) {
	c := New()
	if got := c.Unwrap(); !reflect.DeepEqual(got, &c.Container) {
		t.Errorf("ContainerWrapper.Unwrap() = %v, want %v", got, &c.Container)
		return
	}
}

func TestContainerWrapper_Visualize(t *testing.T) {
	c := New()
	if err := c.Visualize(bytes.NewBuffer(nil)); (err != nil) != false {
		t.Errorf("ContainerWrapper.Visualize() error = %v, wantErr %v", err, false)
		return
	}
}
