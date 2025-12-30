package task

import "time"

type Task struct {
	ID          uint64
	UserID      uint64
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
