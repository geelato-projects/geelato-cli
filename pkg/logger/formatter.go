package logger

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type JSONFormatter struct {
	TimestampFormat string
	PrettyPrint    bool
}

func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)

	for k, v := range entry.Data {
		data[k] = v
	}

	data["level"] = entry.Level.String()
	data["time"] = entry.Time.Format(f.TimestampFormat)
	data["message"] = entry.Message

	if f.PrettyPrint {
		return json.MarshalIndent(data, "", "  ")
	}

	return json.Marshal(data)
}

func NewJSONFormatter(timestampFormat string, prettyPrint bool) *JSONFormatter {
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	return &JSONFormatter{
		TimestampFormat: timestampFormat,
		PrettyPrint:     prettyPrint,
	}
}

type TextFormatter struct {
	TimestampFormat string
	DisableColors   bool
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := entry.Level.String()
	message := entry.Message

	var output string
	if f.DisableColors {
		output = fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)
	} else {
		var color string
		switch entry.Level {
		case logrus.DebugLevel:
			color = "\033[37m"
		case logrus.InfoLevel:
			color = "\033[32m"
		case logrus.WarnLevel:
			color = "\033[33m"
		case logrus.ErrorLevel:
			color = "\033[31m"
		case logrus.FatalLevel, logrus.PanicLevel:
			color = "\033[35m"
		default:
			color = "\033[0m"
		}
		output = fmt.Sprintf("[%s] \033[%sm[%s]\033[0m %s\n", timestamp, color, level, message)
	}

	return []byte(output), nil
}
