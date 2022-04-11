package models

import "time"

// ErrorMessage for kafka dead-letter-queue topic.
type ErrorMessage struct {
	MessageID string    `json:"messageId"`
	Offset    int64     `json:"offset"`
	Partition int       `json:"partition"`
	Topic     string    `json:"topic"`
	Error     string    `json:"error"`
	Time      time.Time `json:"time"`
}

// SuccessMessage for kafka success-letter-queue topic.
type SuccessMessage struct {
	MessageID string    `json:"messageId"`
	Offset    int64     `json:"offset"`
	Partition int       `json:"partition"`
	Topic     string    `json:"topic"`
	Task      Task      `json:"task"`
	Time      time.Time `json:"time"`
}
