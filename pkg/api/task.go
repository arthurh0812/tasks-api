package api

import "time"

type Task struct {
	Title string `json:"title"`
	Text string `json:"text"`
	Date time.Time `json:"date"`
	Done bool `json:"bool"`
}
