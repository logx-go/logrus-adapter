package logrusadapter

import (
	"fmt"

	"github.com/logx-go/commons/pkg/commons"
	"github.com/logx-go/contract/pkg/logx"
	"github.com/sirupsen/logrus"
)

var _ logx.Logger = (*LogrusAdapter)(nil)
var _ logx.Adapter = (*LogrusAdapter)(nil)

// New returns a pointer to a new instance of LogrusAdapter
func New(logger *logrus.Logger) logx.Adapter {
	return &LogrusAdapter{
		logger:    logger,
		fields:    make(map[string]any),
		formatter: nil,
		logLevelMap: map[int]logrus.Level{
			logx.LogLevelDebug:   logrus.DebugLevel,
			logx.LogLevelInfo:    logrus.InfoLevel,
			logx.LogLevelNotice:  logrus.InfoLevel,
			logx.LogLevelWarning: logrus.WarnLevel,
			logx.LogLevelError:   logrus.ErrorLevel,
			logx.LogLevelFatal:   logrus.FatalLevel,
			logx.LogLevelPanic:   logrus.PanicLevel,
		},
		logLevelDefault: logx.LogLevelInfo,
	}
}

// LogrusAdapter implementation to wrap a formatMessage Logger
type LogrusAdapter struct {
	logger          *logrus.Logger
	formatter       logx.Formatter
	fields          map[string]any
	logLevelMap     map[int]logrus.Level
	logLevelDefault int
}

func (l *LogrusAdapter) clone() *LogrusAdapter {
	return &LogrusAdapter{
		logger:          l.logger,
		fields:          commons.CloneFieldMap(l.fields),
		formatter:       l.formatter,
		logLevelMap:     l.logLevelMap,
		logLevelDefault: l.logLevelDefault,
	}
}

func (l *LogrusAdapter) formatMessage(v ...any) string {
	if len(v) < 1 {
		if l.formatter == nil {
			return ""
		}

		return l.formatter.Format("", l.fields)
	}

	msg := fmt.Sprintf(`%v`, v[0])
	if l.formatter == nil {
		return msg
	}

	return l.formatter.Format(msg, l.fields)
}

func (l *LogrusAdapter) prepareFields(v ...any) map[string]any {
	fields := commons.CloneFieldMap(l.fields)

	// skip the first entry for it's the plain message
	for i := 1; i < len(v); i += 2 {
		fieldName := ""
		if i < len(v) {
			fieldName = fmt.Sprintf(`%v`, v[i])
		}

		var fieldValue any
		if i+1 < len(v) {
			fieldValue = v[i+1]
		}

		fields[fieldName] = fieldValue
	}

	return fields
}

func (l *LogrusAdapter) prepareEntry(v ...any) (*logrus.Entry, logrus.Level, string) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(2, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.fields = c.prepareFields(v...)

	lvl := c.resolveLevel(c.fields)
	msg := c.formatMessage(v...)
	logger := c.logger.WithFields(commons.FilterFieldsByName(c.fields, logx.FieldNameLogLevel))

	return logger, lvl, msg
}

func (l *LogrusAdapter) resolveLevel(fields map[string]any) logrus.Level {
	lvl := commons.GetFieldAsIntOrElse(logx.FieldNameLogLevel, fields, l.logLevelDefault)

	if s, ok := l.logLevelMap[lvl]; ok {
		return s
	}

	return logrus.InfoLevel
}

func (l *LogrusAdapter) Fatal(v ...any) {
	entry, _, message := l.prepareEntry(v...)
	entry.Fatal(message)
}

func (l *LogrusAdapter) Panic(v ...any) {
	entry, _, message := l.prepareEntry(v...)
	entry.Panic(message)
}

func (l *LogrusAdapter) Print(v ...any) {
	entry, level, message := l.prepareEntry(v...)
	entry.Log(level, message)
}

func (l *LogrusAdapter) Fatalf(format string, v ...any) {
	entry, _, _ := l.prepareEntry()
	entry.Fatal(fmt.Sprintf(format, v...))
}

func (l *LogrusAdapter) Panicf(format string, v ...any) {
	entry, _, _ := l.prepareEntry()
	entry.Panic(fmt.Sprintf(format, v...))
}

func (l *LogrusAdapter) Printf(format string, v ...any) {
	entry, level, _ := l.prepareEntry()
	entry.Log(level, fmt.Sprintf(format, v...))
}

func (l *LogrusAdapter) Debug(v ...any) {
	entry, level, message := l.prepareEntry(append(v, logx.FieldNameLogLevel, logx.LogLevelDebug)...)
	entry.Log(level, message)
}

func (l *LogrusAdapter) Info(v ...any) {
	entry, level, message := l.prepareEntry(append(v, logx.FieldNameLogLevel, logx.LogLevelInfo)...)
	entry.Log(level, message)
}

func (l *LogrusAdapter) Notice(v ...any) {
	entry, level, message := l.prepareEntry(append(v, logx.FieldNameLogLevel, logx.LogLevelNotice)...)
	entry.Log(level, message)
}

func (l *LogrusAdapter) Warning(v ...any) {
	entry, level, message := l.prepareEntry(append(v, logx.FieldNameLogLevel, logx.LogLevelWarning)...)
	entry.Log(level, message)
}

func (l *LogrusAdapter) Error(v ...any) {
	entry, level, message := l.prepareEntry(append(v, logx.FieldNameLogLevel, logx.LogLevelError)...)
	entry.Log(level, message)
}

func (l *LogrusAdapter) Debugf(format string, v ...any) {
	entry, level, _ := l.prepareEntry(nil, logx.FieldNameLogLevel, logx.LogLevelDebug)
	entry.Log(level, fmt.Sprintf(format, v...))
}

func (l *LogrusAdapter) Infof(format string, v ...any) {
	entry, level, _ := l.prepareEntry(nil, logx.FieldNameLogLevel, logx.LogLevelInfo)
	entry.Log(level, fmt.Sprintf(format, v...))
}

func (l *LogrusAdapter) Noticef(format string, v ...any) {
	entry, level, _ := l.prepareEntry(nil, logx.FieldNameLogLevel, logx.LogLevelNotice)
	entry.Log(level, fmt.Sprintf(format, v...))
}

func (l *LogrusAdapter) Warningf(format string, v ...any) {
	entry, level, _ := l.prepareEntry(nil, logx.FieldNameLogLevel, logx.LogLevelWarning)
	entry.Log(level, fmt.Sprintf(format, v...))
}

func (l *LogrusAdapter) Errorf(format string, v ...any) {
	entry, level, _ := l.prepareEntry(nil, logx.FieldNameLogLevel, logx.LogLevelError)
	entry.Log(level, fmt.Sprintf(format, v...))
}

func (l *LogrusAdapter) WithField(name string, value any) logx.Logger {
	c := l.clone()
	c.fields[name] = value

	return c
}

func (l *LogrusAdapter) WithFormatter(formatter logx.Formatter) logx.Adapter {
	c := l.clone()
	c.formatter = formatter

	return c
}
