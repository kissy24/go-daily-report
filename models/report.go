package models

import "time"

type Report struct {
	ID      int
	Content string
	Date    time.Time
}
