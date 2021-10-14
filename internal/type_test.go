package internal

import (
	"testing"
)

func Test_initDigProvideOptionsType(t *testing.T) {
	want := "dig.provideOptions"
	got := initDigProvideOptionsType().String()
	if want != got {
		t.Errorf("initDigProvideOptionsType() want %s, got %s", want, got)
	}
}
