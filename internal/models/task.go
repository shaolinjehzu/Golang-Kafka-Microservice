package models

// Task for kafka update-balance topic.
type Task struct {
	Id    int64
	Type  string
	Value string
}
