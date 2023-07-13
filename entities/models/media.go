package models

import "time"

type Media struct {
	Id            string                 `json:"id"`
	Attributes    map[string]interface{} `json:"attributes"`
	Height        float64                `json:"height"`
	Width         float64                `json:"width"`
	X             float64                `json:"x"`
	Y             float64                `json:"y"`
	Url           string                 `json:"url"`
	PHash         string                 `json:"phash"`
	Filename      string                 `json:"filename"`
	MediaUrl      string                 `json:"media_url"`
	ScreenshotUrl string                 `json:"screenshot_url"`
	TaskId        string                 `json:"task_id"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	DeletedAt     *time.Time             `json:"deleted_at"`
}
