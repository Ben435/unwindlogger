package unwindlogger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Entry struct {
	logger  *Logger
	ctx     context.Context
	level   Level
	message string
	time    time.Time
	fields  map[string]interface{}
}

func NewEntry(logger *Logger, ctx context.Context) *Entry {
	return &Entry{
		logger:  logger,
		ctx:     ctx,
		level:   0,
		message: "",
		time:    time.Now(),
		fields:  make(map[string]interface{}),
	}
}

func (e *Entry) WithField(k string, v interface{}) *Entry {
	if v != nil {
		e.fields[k] = v
	}
	return e
}

func (e *Entry) WithError(err error) *Entry {
	return e.WithField("error", err.Error())
}

func (e *Entry) Bytes() ([]byte, error) {
	data := make(map[string]interface{})

	for k, v := range e.fields {
		data[k] = v
	}
	data["level"] = e.level.Format()
	data["msg"] = e.message
	data["time"] = e.time.Format(time.RFC3339)

	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry: %w", err)
	}

	return append(b, '\n'), nil
}

func (e *Entry) String() (string, error) {
	b, err := e.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *Entry) Debug(msg string) {
	e.Log(DEBUG, msg)
}

func (e *Entry) Info(msg string) {
	e.Log(INFO, msg)
}

func (e *Entry) Warn(msg string) {
	e.Log(WARN, msg)
}

func (e *Entry) Error(msg string) {
	e.Log(ERROR, msg)
}

func (e *Entry) Log(level Level, msg string) {
	e.message = msg
	e.level = level

	if !e.logger.shouldImmediatelyLog(e.level) {
		e.logger.deferEntry(e)
		return
	}

	e.write()
}

func (e *Entry) write() {
	b, err := e.Bytes()
	if err != nil {
		panic(err)
	}

	if _, err := e.logger.out.Write(b); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to write to log, %v\n", err)
	}
}
