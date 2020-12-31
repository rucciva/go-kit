package rzap

import (
	"testing"
)

type tWritter struct {
	t *testing.T
}

func (tw *tWritter) Write(b []byte) (int, error) {
	tw.t.Log(string(b[:len(b)-1]))
	return len(b) - 1, nil
}

func (tw *tWritter) Sync() error {
	return nil
}
