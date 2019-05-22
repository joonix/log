package log

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Default logrus to FluentD severity map
var SeverityMap = map[string]string{
	"panic":   "fatal",
	"fatal":   "fatal",
	"warning": "warn",
	"debug":   "debug",
	"error":   "error",
	"trace":   "trace",
	"info":    "info",
}

// logrus to stackdriver severity map
func UseStackdriverSeverity() {
	SeverityMap = map[string]string{
		"panic":   "CRITICAL",
		"fatal":   "CRITICAL",
		"warning": "WARNING",
		"debug":   "DEBUG",
		"error":   "ERROR",
		"trace":   "DEBUG",
		"info":    "INFO",
	}
}

// FluentdFormatter is similar to logrus.JSONFormatter but with log level that are recongnized
// by kubernetes fluentd.
type FluentdFormatter struct {
	TimestampFormat string
	SeverityMap     map[string]string
}

// Format the log entry. Implements logrus.Formatter.
func (f *FluentdFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	prefixFieldClashes(data, entry.HasCaller())

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339Nano
	}

	data["timestamp"] = entry.Time.Format(timestampFormat)
	data["message"] = entry.Message

	if ms, ok := SeverityMap[entry.Level.String()]; ok {
		data["severity"] = ms
	} else {
		data["severity"] = SeverityMap["debug"]
	}

	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if funcVal != "" {
			data[logrus.FieldKeyFunc] = funcVal
		}
		if fileVal != "" {
			data[logrus.FieldKeyFile] = fileVal
		}
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

func prefixFieldClashes(data logrus.Fields, reportCaller bool) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}

	if m, ok := data["message"]; ok {
		data["fields.message"] = m
	}

	if l, ok := data["timestamp"]; ok {
		data["fields.timestamp"] = l
	}

	if l, ok := data["severity"]; ok {
		data["fields.severity"] = l
	}

	if reportCaller {
		if l, ok := data[logrus.FieldKeyFunc]; ok {
			data["fields."+logrus.FieldKeyFunc] = l
		}
		if l, ok := data[logrus.FieldKeyFile]; ok {
			data["fields."+logrus.FieldKeyFile] = l
		}
	}
}
