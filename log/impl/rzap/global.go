package rzap

import (
	"testing"

	"github.com/rucciva/go-kit/log"
)

func init() {
	log.Register(NewLogger())
}

// ChangeT change underlying `*testing.T` of `tWritter` if it's being used
// not thread safe
func ChangeT(l log.Logger, t *testing.T) (restorer func()) {
	restorer = func() {}

	r, ok := l.(*Logger)
	if !ok {
		return
	}
	tw, ok := r.cfg.ws.(*tWritter)
	if !ok {
		return
	}
	tOld := tw.t
	tw.t = t
	restorer = func() {
		tw.t = tOld
	}
	return
}

// GetCapturer return `Capturer` if it being used the Logger
func GetCapturer(l log.Logger) *Capturer {
	r, ok := l.(*Logger)
	if !ok {
		return nil
	}
	tw, _ := r.cfg.ws.(*Capturer)
	return tw
}

// Close flush the buffer
// if l is not Logger than nothing happens
func Close(l log.Logger) error {
	rz, ok := l.(*Logger)
	if !ok {
		return nil
	}
	return rz.Close()
}
