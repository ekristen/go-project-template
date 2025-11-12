package common

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ConsoleFormatter formats logs in a human-readable console format similar to zerolog
type ConsoleFormatter struct {
	TimestampFormat string
	NoColor         bool
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorWhite  = "\033[97m"

	colorBoldRed    = "\033[1;31m"
	colorBoldGreen  = "\033[1;32m"
	colorBoldYellow = "\033[1;33m"
	colorBoldCyan   = "\033[1;36m"
)

// Format renders a single log entry
func (f *ConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	// Timestamp
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = "3:04:05PM"
	}

	timestamp := entry.Time.Format(timestampFormat)
	if !f.NoColor {
		b.WriteString(colorGray)
	}
	b.WriteString(timestamp)
	if !f.NoColor {
		b.WriteString(colorReset)
	}
	b.WriteString(" ")

	// Log Level
	level := strings.ToUpper(entry.Level.String())
	levelColor := f.getLevelColor(entry.Level)

	if !f.NoColor {
		b.WriteString(levelColor)
	}
	b.WriteString(fmt.Sprintf("%-5s", level))
	if !f.NoColor {
		b.WriteString(colorReset)
	}
	b.WriteString(" ")

	// Caller information (if present)
	if entry.HasCaller() {
		caller := fmt.Sprintf("%s:%d",
			filepath.Base(entry.Caller.File),
			entry.Caller.Line)

		if !f.NoColor {
			b.WriteString(colorGray)
		}
		b.WriteString(caller)
		if !f.NoColor {
			b.WriteString(colorReset)
		}
		b.WriteString(" > ")
	}

	// Message
	if !f.NoColor {
		b.WriteString(colorWhite)
	}
	b.WriteString(entry.Message)
	if !f.NoColor {
		b.WriteString(colorReset)
	}

	// Fields
	if len(entry.Data) > 0 {
		f.writeFields(&b, entry)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// getLevelColor returns the color for a log level
func (f *ConsoleFormatter) getLevelColor(level logrus.Level) string {
	if f.NoColor {
		return ""
	}

	switch level {
	case logrus.TraceLevel:
		return colorPurple
	case logrus.DebugLevel:
		return colorCyan
	case logrus.InfoLevel:
		return colorGreen
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel:
		return colorRed
	case logrus.FatalLevel, logrus.PanicLevel:
		return colorBoldRed
	default:
		return colorWhite
	}
}

// writeFields writes the log fields in a key=value format
func (f *ConsoleFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	// Sort keys for consistent output
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		b.WriteByte(' ')

		// Key
		if !f.NoColor {
			b.WriteString(colorCyan)
		}
		b.WriteString(key)
		if !f.NoColor {
			b.WriteString(colorReset)
		}
		b.WriteByte('=')

		// Value
		value := entry.Data[key]
		f.writeValue(b, value)
	}
}

// writeValue writes a field value with appropriate formatting
func (f *ConsoleFormatter) writeValue(b *bytes.Buffer, value interface{}) {
	if !f.NoColor {
		b.WriteString(colorWhite)
	}

	switch v := value.(type) {
	case string:
		// Quote strings that contain spaces
		if strings.Contains(v, " ") {
			fmt.Fprintf(b, "%q", v)
		} else {
			b.WriteString(v)
		}
	case error:
		if !f.NoColor {
			b.WriteString(colorRed)
		}
		fmt.Fprintf(b, "%q", v.Error())
	case time.Time:
		b.WriteString(v.Format(time.RFC3339))
	case time.Duration:
		b.WriteString(v.String())
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, bool:
		fmt.Fprintf(b, "%v", v)
	default:
		fmt.Fprintf(b, "%v", v)
	}

	if !f.NoColor {
		b.WriteString(colorReset)
	}
}
