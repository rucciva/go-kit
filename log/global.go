package log

var (
	logger Logger = defaultLogger{defaultBaseLogger: defaultBaseLogger{}}
)

// Register register a global Log Implementator
// return error when Log Implementator already registered before
func Register(log Logger) {
	if log == nil {
		return
	}
	logger = log
}

// GetGlobal return the global logger implementator
func GetGlobal() Logger {
	return logger
}

func Debug(msg string) {
	logger.Debug(msg)
}
func Info(msg string) {
	logger.Info(msg)
}
func Warn(msg string) {
	logger.Warn(msg)
}
func Error(msg string) {
	logger.Error(msg)
}
func Fatal(msg string) {
	logger.Fatal(msg)
}
func Panic(msg string) {
	logger.Panic(msg)
}
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}
func WithFields(KeyValuePairs ...interface{}) BaseLogger {
	return logger.WithFields(KeyValuePairs)
}
