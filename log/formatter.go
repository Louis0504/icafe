package log

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	traceIDKey = "X-B3-Traceid" // equal with telemetry
)

var (
	processID = strconv.Itoa(os.Getpid())
)

var (
	severityMap = map[logrus.Level]byte{
		logrus.TraceLevel: 'D',
		logrus.DebugLevel: 'D',
		logrus.InfoLevel:  'I',
		logrus.WarnLevel:  'W',
		logrus.ErrorLevel: 'E',
		logrus.FatalLevel: 'F',
		logrus.PanicLevel: 'F',
	}
)

type formatter struct{}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Buffer.WriteByte('[')
	entry.Buffer.WriteByte(severityMap[entry.Level])
	entry.Buffer.WriteByte(' ')
	entry.Buffer.WriteString(entry.Time.Format("2006-01-02 15:04:05.000"))
	entry.Buffer.WriteByte(' ')

	var traceID string
	if entry.Context != nil {
		traceID, _ = entry.Context.Value(traceIDKey).(string)
	}
	if traceID != "" {
		entry.Buffer.WriteString(traceID)
	} else {
		entry.Buffer.WriteByte('-')
	}
	entry.Buffer.WriteByte(' ')
	if entry.Caller != nil {
		entry.Buffer.WriteString(filepath.Base(entry.Caller.File))
		entry.Buffer.WriteByte(':')
		entry.Buffer.WriteString(strconv.Itoa(entry.Caller.Line))
	} else {
		entry.Buffer.WriteByte('-')
	}

	entry.Buffer.WriteByte(' ')
	entry.Buffer.WriteString(hostName())
	entry.Buffer.WriteByte(':')
	entry.Buffer.WriteString(processID)
	entry.Buffer.WriteString("] ")
	if len(entry.Data) > 0 {
		if extraData, err := json.Marshal(entry.Data); err == nil {
			entry.Buffer.WriteString("[")
			entry.Buffer.Write(extraData)
			entry.Buffer.WriteString("] ")
		}
	}
	entry.Buffer.WriteString(entry.Message)
	entry.Buffer.WriteByte('\n')

	return entry.Buffer.Bytes(), nil
}

func hostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return regexp.MustCompile("[^0-9a-zA-Z]+").ReplaceAllString(hostname, "_")
}
