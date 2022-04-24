package main

import (
	"database/sql"
	"encoding/json"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type stateSinkCallback struct {
	db *sql.DB
}

type eventSinkCallback struct {
	db *sql.DB
}

type metricSinkCallback struct {
	db *sql.DB
}

type Metric struct {
	Time    string  `json:"time"`
	Device  string  `json:"device"`
	Channel string  `json:"channel"`
	Value   float32 `json:"value"`
}

func (cb *stateSinkCallback) Handle(data []byte) error {
	log.WithFields(log.Fields{
		"bus":  "state",
		"data": string(data),
	}).Debug("data received on message bus")
	return nil
}

func (cb *eventSinkCallback) Handle(data []byte) error {
	log.WithFields(log.Fields{
		"bus":  "event",
		"data": string(data),
	}).Debug("data received on message bus")
	return nil
}

func (cb *metricSinkCallback) Handle(data []byte) error {
	var metric Metric

	re := regexp.MustCompile(".*{")
	message := re.ReplaceAllString(string(data), "{")
	if err := json.Unmarshal([]byte(message), &metric); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"time":    metric.Time,
		"device":  metric.Device,
		"channel": metric.Channel,
		"value":   metric.Value,
	}).Debug("handler received metric")

	sql := "INSERT INTO metrics (time, device, channel, value) VALUES ($1, $2, $3, $4)"

	tx, err := cb.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(metric.Time, metric.Device, metric.Channel, metric.Value); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
