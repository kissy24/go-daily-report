package models

import "time"

type Report struct {
	ID      int
	Title   string
	Content string
	Date    time.Time
}
