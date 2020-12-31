package log

type mBaseLogger struct {
	bl []BaseLogger
}

func (mbl mBaseLogger) Debug(msg string) {
	for _, bl := range mbl.bl {
		bl.Debug(msg)
	}
}
func (mbl mBaseLogger) Info(msg string) {
	for _, bl := range mbl.bl {
		bl.Info(msg)
	}
}
func (mbl mBaseLogger) Warn(msg string) {
	for _, bl := range mbl.bl {
		bl.Warn(msg)
	}
}
func (mbl mBaseLogger) Error(msg string) {
	for _, bl := range mbl.bl {
		bl.Error(msg)
	}
}
func (mbl mBaseLogger) Fatal(msg string) {
	for _, bl := range mbl.bl {
		bl.Fatal(msg)
	}
}
func (mbl mBaseLogger) Panic(msg string) {
	for _, bl := range mbl.bl {
		bl.Panic(msg)
	}
}

type mLogger struct {
	mBaseLogger
	l []Logger
}

func NewMultiplexer(l ...Logger) Logger {
	ml := mLogger{
		mBaseLogger: mBaseLogger{bl: make([]BaseLogger, 0, len(l))},
		l:           l,
	}
	for _, l := range l {
		ml.bl = append(ml.bl, l.WithFields())
	}
	return ml
}

func (ml mLogger) WithFields(KeyValuePairs ...interface{}) BaseLogger {
	mbl := mBaseLogger{bl: make([]BaseLogger, 0, len(ml.l))}
	for _, l := range ml.l {
		mbl.bl = append(mbl.bl, l.WithFields(KeyValuePairs...))
	}
	return mbl
}
