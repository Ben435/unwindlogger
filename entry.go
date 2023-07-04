package unwindlogger

import (
	"encoding/json"
	"fmt"
	"time"
)

type Entry struct {
	logger  *Logger
	message string
	time    time.Time
	fields  map[string]interface{}
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		logger:  logger,
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

func (e *Entry) Bytes() ([]byte, error) {
	data := make(map[string]interface{})

	for k, v := range e.fields {
		data[k] = v
	}
	data["message"] = e.message
	data["time"] = e.time.Format(time.RFC3339)

	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry: %w", err)
	}

	return b, nil
}

func (e *Entry) String() (string, error) {
	b, err := e.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *Entry) Debug(msg string) {
	e.message = msg

	b, err := e.Bytes()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%b", b)
}
