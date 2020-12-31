package rzap

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rucciva/go-kit/log"
	"golang.org/x/xerrors"
)

func tInvokeLog(l log.Logger, preline string) {
	err := errors.New("first error string (unformattable)")
	err1 := xerrors.Errorf("second error string")
	err2 := xerrors.Errorf("third error string")
	data := map[string]interface{}{"data": map[string]interface{}{"test": true}}

	fmt.Println(preline)
	fmt.Println()
	l.WithFields("data", data, "error", err, "error1", err1, "error2", err2).Error("error")
	l.WithFields("data", data, "error", err, "error1", err1, "error2", err2).Error("error")
	fmt.Print("\n\n\n")

}

func TestLog(t *testing.T) {
	pl := NewLogger()
	t.Log("|||=== production logger ===|||")
	tInvokeLog(pl, "|||=== production logger ===|||")

	dl := NewLogger(WithConsoleEncoder(), WithErrorStackTrace(), WithCallerInfo())
	t.Log("|||=== development logger ===|||")
	tInvokeLog(dl, "|||=== development logger ===|||")

	tl := NewLogger(WithCallerInfo(), WithTWritter(t), WithConsoleEncoder())
	t.Log("|||=== test logger ===|||")
	tInvokeLog(tl, "|||=== test logger ===|||")

	tl1 := NewLogger(WithCallerInfo(), WithTWritter(t), WithConsoleEncoder(), WithErrorStackTrace())
	t.Log("|||=== test logger with stacktrace ===|||")
	tInvokeLog(tl1, "|||=== test logger stacktrace ===|||")

	t.Run("subtest", func(t *testing.T) {
		defer ChangeT(tl1, t)()
		t.Log("|||=== sub test logger with stacktrace ===|||")
		tInvokeLog(tl1, "|||=== sub test logger stacktrace ===|||")
	})

	t.Log("|||=== re-test logger with stacktrace ===|||")
	tInvokeLog(tl1, "|||=== re-test logger stacktrace ===|||")
}

func TestCaptureDefaultLog(t *testing.T) {
	type data struct {
		Test int `json:"test"`
	}

	type logged struct {
		Level string    `json:"level"`
		Time  time.Time `json:"time"`
		Msg   string    `json:"msg"`
		Data  data      `json:"data"`
		Error string    `json:"error"`
	}

	table := []struct {
		l   log.Level
		f   data
		m   string
		err error
	}{
		{log.DebugLevel, data{Test: 1}, "first", errors.New("first")},
		{log.InfoLevel, data{Test: 2}, "second", errors.New("second")},
		{log.WarnLevel, data{Test: 3}, "third", errors.New("third")},
		{log.ErrorLevel, data{Test: 4}, "fourth", errors.New("fourth")},
	}

	levelMap := map[string]log.Level{
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
	}

	lc := NewLogger(WithCapturer(false), WithLevel(log.DebugLevel))
	lcf := NewLogger(WithCapturer(true), WithLevel(log.DebugLevel))
	l := log.NewMultiplexer(lc, lcf,
		NewLogger(WithCallerInfo(), WithCallerSkip(1), WithConsoleEncoder(), WithErrorStackTrace(), WithTWritter(t)),
	)

	before := time.Now()
	for _, d := range table {
		switch d.l {
		case log.DebugLevel:
			l.WithFields("data", d.f, "error", d.err).Debug(d.m)
		case log.InfoLevel:
			l.WithFields("data", d.f, "error", d.err).Info(d.m)
		case log.WarnLevel:
			l.WithFields("data", d.f, "error", d.err).Warn(d.m)
		case log.ErrorLevel:
			l.WithFields("data", d.f, "error", d.err).Error(d.m)
		}
	}
	after := time.Now()
	if err := Close(l); err != nil {
		t.Fatalf("error when closing log: %+v", err)
	}
	cs := []*Capturer{
		GetCapturer(lc), GetCapturer(lcf),
	}
	for _, c := range cs {
		if c == nil {
			t.Fatal("should return capturer")
		}
		if len(c.Lines()) != len(table) {
			t.Fatalf("should capture all the logs, got %d line, want %d line", len(c.Lines()), len(table))
		}
		for i := len(table) - 1; i >= 0; i-- {
			d := table[i]

			s, ok := c.Pop()
			if !ok {
				t.Fatal("should return recorded log")
			}
			var v logged
			if err := json.Unmarshal([]byte(s), &v); err != nil {
				t.Fatalf("should format log as json: %v", err)
			}

			if levelMap[v.Level] != d.l {
				t.Errorf("should log correct level, got %s want %s", v.Level, d.l.String())
			}
			if v.Time.Before(before) || v.Time.After(after) {
				t.Errorf("logged time should be between %v and %v, got %v", before, after, v.Time)
			}
			if v.Msg != d.m {
				t.Errorf("should log correct message, got %s want %s", v.Msg, d.m)
			}
			if v.Data.Test != d.f.Test {
				t.Errorf("should log correct data, got %v want %v", v.Data, d.f)
			}
			if v.Error != d.err.Error() {
				t.Errorf("should log correct error, got %v want %v", v.Error, d.err.Error())
			}
		}

		if _, ok := c.Pop(); ok {
			t.Fatal("should not return recorded log anymore")
		}
	}
}
