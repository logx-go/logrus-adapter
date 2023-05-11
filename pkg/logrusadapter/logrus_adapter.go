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

// LogrusAdapter implementation to wrap a format Logger
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
		fields:          l.fields,
		formatter:       l.formatter,
		logLevelMap:     l.logLevelMap,
		logLevelDefault: l.logLevelDefault,
	}
}

func (l *LogrusAdapter) format(v ...any) (messageF string, fieldsZ logrus.Fields) {
	fields := commons.FilterFieldsByName(l.fields, logx.FieldNameLogLevel)
	if len(v) < 1 {
		if l.formatter == nil {
			return "", fields
		}

		return l.formatter.Format("", fields)
	}

	msg := fmt.Sprintf(`%v`, v[0])

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

	if l.formatter == nil {
		return msg, fields
	}

	return l.formatter.Format(msg, fields)
}

func (l *LogrusAdapter) convertLogrusLevel(fields map[string]any) logrus.Level {
	lvl := commons.GetFieldAsIntOrElse(logx.FieldNameLogLevel, fields, l.logLevelDefault)

	if s, ok := l.logLevelMap[lvl]; ok {
		return s
	}

	return logrus.InfoLevel
}

func (l *LogrusAdapter) Fatal(v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	msgZ, fieldsZ := c.format(v...)
	c.logger.WithFields(fieldsZ).Fatal(msgZ)
}

func (l *LogrusAdapter) Panic(v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	msgZ, fieldsZ := c.format(v...)
	c.logger.WithFields(fieldsZ).Panic(msgZ)
}

func (l *LogrusAdapter) Print(v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	msgZ, fieldsZ := c.format(v...)
	c.logger.WithFields(fieldsZ).Log(c.convertLogrusLevel(c.fields), msgZ)
}

func (l *LogrusAdapter) Fatalf(format string, v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	msgZ, fieldsZ := c.format(fmt.Sprintf(format, v...))
	c.logger.WithFields(fieldsZ).Fatal(msgZ)
}

func (l *LogrusAdapter) Panicf(format string, v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	msgZ, fieldsZ := c.format(fmt.Sprintf(format, v...))
	c.logger.WithFields(fieldsZ).Panic(msgZ)
}

func (l *LogrusAdapter) Printf(format string, v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	msgZ, fieldsZ := c.format(fmt.Sprintf(format, v...))
	c.logger.WithFields(fieldsZ).Log(c.convertLogrusLevel(c.fields), msgZ)
}

func (l *LogrusAdapter) Debug(v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelDebug).Print(v...)
}

func (l *LogrusAdapter) Info(v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelInfo).Print(v...)
}

func (l *LogrusAdapter) Notice(v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelNotice).Print(v...)
}

func (l *LogrusAdapter) Warning(v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelWarning).Print(v...)
}

func (l *LogrusAdapter) Error(v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelError).Print(v...)
}

func (l *LogrusAdapter) Debugf(format string, v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelDebug).Printf(format, v...)
}

func (l *LogrusAdapter) Infof(format string, v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelInfo).Printf(format, v...)
}

func (l *LogrusAdapter) Noticef(format string, v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelNotice).Printf(format, v...)
}

func (l *LogrusAdapter) Warningf(format string, v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelWarning).Printf(format, v...)
}

func (l *LogrusAdapter) Errorf(format string, v ...any) {
	c := l.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelError).Printf(format, v...)
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
