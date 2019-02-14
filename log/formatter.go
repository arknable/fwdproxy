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
	head := f.colorized(entry, fmt.Sprintf("%20s:", key))
	val := color.LightGray(fmt.Sprintf("%v", value))
	return fmt.Sprintf("%s %s\n", head, val)
}

// Set color of the string
func (f *TextFormatter) colorized(entry *logrus.Entry, s string) string {
	switch entry.Level {
	case logrus.PanicLevel:
	case logrus.FatalLevel:
		return color.GRed(s)
	case logrus.ErrorLevel:
		return color.BRed(s)
	case logrus.WarnLevel:
		return color.BYellow(s)
	case logrus.InfoLevel:
		return color.BLightGray(s)
	case logrus.DebugLevel:
		return color.BDarkGray(s)
	case logrus.TraceLevel:
		return color.BPurple(s)
	}
	return color.BDarkGray(s)
}
