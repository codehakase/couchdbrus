// Package couchdbrus defines a CouchDB Hook for the Logrus package
package couchdbrus

import (
	"fmt"
	"net/url"

	"github.com/ryanjyoder/couchdb"
	"github.com/sirupsen/logrus"
)

// Hook is a logrus Hook for dispatching messages to a couchdb bucket
type Hook struct {
	WithLevels   []logrus.Level
	DbClient     *couchdb.Client
	DB           couchdb.DatabaseService
	LogTypeTitle string
	Options      map[string]interface{}
}

// LogData defines the structure for persisted log items
type LogData struct {
	couchdb.Document
	Type     string      `json:"type"`
	LogEntry interface{} `json:"log_entry"`
}

// NewHook creates a new Hook with a user defined couchdb cient
func NewHook(client *couchdb.Client, dbname string) (*Hook, error) {
	_, err := url.Parse(client.BaseURL.String())
	if err != nil {
		return nil, fmt.Errorf("invalid BaseURL contained in client: %v", err)
	}
	db, err := createOrUseDb(client, dbname)
	if err != nil {
		return nil, fmt.Errorf("could not create db service: %v", err)
	}
	hook := &Hook{
		DbClient: client,
		DB:       db,
	}
	return hook, nil
}

// Fire a new event to couchdb
func (h *Hook) Fire(entry *logrus.Entry) error {
	entries := h.newEntry(entry)
	data := LogData{
		Type:     h.LogTypeTitle,
		LogEntry: entries.Message,
	}
	_, err := h.DB.Post(&data)
	if err != nil {
		return fmt.Errorf("could not post log entry: %v", err)
	}
	return nil
}

// newEntry adds a new log item
func (h *Hook) newEntry(entry *logrus.Entry) *logrus.Entry {
	data := map[string]interface{}{}
	for k, v := range h.Options {
		data[k] = v
	}
	for k, v := range entry.Data {
		data[k] = v
	}
	return &logrus.Entry{
		Logger:  entry.Logger,
		Time:    entry.Time,
		Level:   entry.Level,
		Data:    data,
		Message: entry.Message,
	}
}

// SetLevels for hook
func (h *Hook) SetLevels(levels []logrus.Level) {
	h.WithLevels = levels
}

// Levels sent to couchdb
func (h *Hook) Levels() []logrus.Level {
	if h.WithLevels != nil {
		return logrus.AllLevels
	}
	return h.WithLevels
}

func createOrUseDb(client *couchdb.Client, dbname string) (couchdb.DatabaseService, error) {
	_, err := client.Create(dbname)
	if err != nil {
		dberr, ok := err.(*couchdb.Error)
		if !ok {
			return nil, dberr
		}
	}
	return client.Use(dbname), nil
}
