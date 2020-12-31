package rzap

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type plainError struct {
	e error
}

func (pe plainError) Error() string {
	return pe.e.Error()
}

type errorVerboseToStacktraceEncoder struct {
	zapcore.Encoder
}

func (he errorVerboseToStacktraceEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	filteredFields := make([]zapcore.Field, 0)

	var errsF []zapcore.Field
	for _, f := range fields {
		err, ok := f.Interface.(error)
		if ok {
			_, ok = f.Interface.(fmt.Formatter)
		}
		if ok {
			errsF = append(errsF, f)
			filteredFields = append(filteredFields, zap.NamedError(f.Key, plainError{err}))
			continue
		}

		filteredFields = append(filteredFields, f)
	}

	var sb strings.Builder
	for i, err := range errsF {
		sb.WriteString(" === stacktrace for field `")
		sb.WriteString(err.Key)
		sb.WriteString("` ===\n")
		sb.WriteString(fmt.Sprintf(" %+v", err.Interface.(error)))
		if i < len(errsF)-1 {
			sb.WriteString("\n")
		}
	}
	ent.Stack = sb.String()

	return he.Encoder.EncodeEntry(ent, filteredFields)
}
