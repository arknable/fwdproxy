package log

import (
	"fmt"
	"strings"

	"github.com/bclicn/color"
	"github.com/sirupsen/logrus"
)

// TextFormatter is logrus text formatter
type TextFormatter struct{}

// Format implements logrus.Formatter.Format
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var builder strings.Builder
	timeString := entry.Time.Format("02 January 2006 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	builder.WriteString(f.format(entry, level, f.colorized(entry, entry.Message)))
	builder.WriteString(f.format(entry, "time", timeString))

	for k, v := range entry.Data {
		builder.WriteString(f.format(entry, k, v))
	}

	if ((entry.Level == logrus.PanicLevel) ||
		(entry.Level == logrus.FatalLevel) ||
		(entry.Level == logrus.ErrorLevel)) &&
		(entry.Caller != nil) {
		builder.WriteString(f.format(entry, "source", fmt.Sprintf("%s:%v", entry.Caller.File, entry.Caller.Line)))
	}

	builder.WriteString("\n")
	return []byte(builder.String()), nil
}

// Format key and value to output message
func (f *TextFormatter) format(entry *logrus.Entry, key string, value interface{}) string {
	head := color.LightGray(fmt.Sprintf("%20s:", key))
	val := color.LightGray(fmt.Sprintf("%v", value))
	isDebugAndLower := (entry.Level == logrus.DebugLevel) || (entry.Level == logrus.TraceLevel)

	if (strings.ToUpper(entry.Level.String()) == strings.ToUpper(key)) || isDebugAndLower {
		head = f.colorized(entry, fmt.Sprintf("%20s:", key))
	}

	if isDebugAndLower {
		val = f.colorized(entry, fmt.Sprintf("%v", value))
	}

	return fmt.Sprintf("%s %s\n", head, val)
}

// Set color of the string
func (f *TextFormatter) colorized(entry *logrus.Entry, s string) string {
	result := color.BDarkGray(s)
	switch entry.Level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		result = color.BRed(s)
	case logrus.WarnLevel:
		result = color.BYellow(s)
	case logrus.InfoLevel:
		result = color.BBlue(s)
	case logrus.DebugLevel, logrus.TraceLevel:
		result = color.BDarkGray(s)
	}
	return result
}
