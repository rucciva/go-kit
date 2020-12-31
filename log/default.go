package log

import "fmt"

type defaultBaseLogger struct {
	fields []interface{}
}

func (d defaultBaseLogger) Debug(msg string) {
	a := append(make([]interface{}, 0, len(d.fields)+2), DebugLevel, msg)
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Info(msg string) {
	a := append(make([]interface{}, 0, len(d.fields)+2), InfoLevel, msg)
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Warn(msg string) {
	a := append(make([]interface{}, 0, len(d.fields)+2), WarnLevel, msg)
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Error(msg string) {
	a := append(make([]interface{}, 0, len(d.fields)+2), ErrorLevel, msg)
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Fatal(msg string) {
	a := append(make([]interface{}, 0, len(d.fields)+2), FatalLevel, msg)
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Panic(msg string) {
	a := append(make([]interface{}, 0, len(d.fields)+2), PanicLevel, msg)
	fmt.Println(append(a, d.fields...))
}

func (d defaultBaseLogger) Debugf(format string, args ...interface{}) {
	a := append(make([]interface{}, 0, len(d.fields)+2), DebugLevel, fmt.Sprintf(format, args...))
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Infof(format string, args ...interface{}) {
	a := append(make([]interface{}, 0, len(d.fields)+2), InfoLevel, fmt.Sprintf(format, args...))
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Warnf(format string, args ...interface{}) {
	a := append(make([]interface{}, 0, len(d.fields)+2), WarnLevel, fmt.Sprintf(format, args...))
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Errorf(format string, args ...interface{}) {
	a := append(make([]interface{}, 0, len(d.fields)+2), ErrorLevel, fmt.Sprintf(format, args...))
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Fatalf(format string, args ...interface{}) {
	a := append(make([]interface{}, 0, len(d.fields)+2), FatalLevel, fmt.Sprintf(format, args...))
	fmt.Println(append(a, d.fields...))
}
func (d defaultBaseLogger) Panicf(format string, args ...interface{}) {
	a := append(make([]interface{}, 0, len(d.fields)+2), PanicLevel, fmt.Sprintf(format, args...))
	fmt.Println(append(a, d.fields...))
}

type defaultLogger struct {
	defaultBaseLogger
}

func (defaultLogger) WithFields(KeyValuePairs ...interface{}) BaseLogger {
	return defaultBaseLogger{fields: KeyValuePairs}
}
