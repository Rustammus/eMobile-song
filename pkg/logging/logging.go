package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

type writerHook struct {
	Writer   []io.Writer
	LogLevel []logrus.Level
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevel
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return nil
}

func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		DisableQuote:    true,
		FullTimestamp:   false,
		TimestampFormat: "2006-01-02 15:04:05.000",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "msg"},
	})

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer:   []io.Writer{os.Stdout},
		LogLevel: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	logEntry = logrus.NewEntry(l)
}

var logEntry *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func (l *Logger) SetLevel(level logrus.Level) {
	l.Level = level
	l.Logger.SetLevel(level)
}

// GetLogger
// Return previously initialized logger
func GetLogger() Logger {
	return Logger{logEntry}
}

// GetLoggerTest
// Return logger with io.Discard output
func GetLoggerTest() Logger {
	l := logrus.New()
	l.Level = logrus.TraceLevel
	l.SetOutput(io.Discard)
	return Logger{logrus.NewEntry(l)}
}
