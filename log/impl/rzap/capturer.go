package rzap

import (
	"bufio"
	"io"
	"sync"

	"go.uber.org/zap/zapcore"
)

// Capturer is a WriteSyncer that can capture each line writen to it and optionally forward the line to another `WriteSyncer`
type Capturer struct {
	ws zapcore.WriteSyncer

	lines []string
	r     *io.PipeReader
	w     *io.PipeWriter
	wg    sync.WaitGroup
}

func newCapturer(ws zapcore.WriteSyncer) (c *Capturer) {
	c = &Capturer{ws: ws}
	c.r, c.w = io.Pipe()
	c.wg.Add(1)
	go c.parseLine()
	c.wg.Wait()
	return
}

func (c *Capturer) parseLine() {
	sc := bufio.NewScanner(c.r)
	c.wg.Done()
	for sc.Scan() {
		c.lines = append(c.lines, sc.Text())
	}
}

func (c *Capturer) Write(b []byte) (int, error) {
	i, err := c.w.Write(b)
	if c.ws != nil {
		i, err = c.ws.Write(b)
	}
	return i, err
}

func (c *Capturer) Sync() error {
	err := c.w.Close()
	if c.ws != nil {
		err = c.ws.Sync()
	}
	return err
}

// Lines return all captured lines
func (c *Capturer) Lines() []string {
	return c.lines
}

// Pop return the latest captured line and removed it from storage
func (c *Capturer) Pop() (string, bool) {
	if len(c.lines) == 0 {
		return "", false
	}
	s := c.lines[len(c.lines)-1]
	c.lines = c.lines[:len(c.lines)-1]
	return s, true
}

// Reset clear captured lines
func (c *Capturer) Reset() {
	c.lines = nil
}
