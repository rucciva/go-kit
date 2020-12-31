package log

type BaseLogger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)
}

type Logger interface {
	BaseLogger
	WithFields(KeyValuePairs ...interface{}) BaseLogger
}
