package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Media struct {
	bun.BaseModel `bun:"table:media,alias:media"`

	ID            string                 `bun:"id,pk"`
	Attributes    map[string]interface{} `bun:"attributes,type:jsonb,nullzero"`
	Height        float64                `bun:"height,notnull"`
	Width         float64                `bun:"width,notnull"`
	X             float64                `bun:"x,notnull"`
	Y             float64                `bun:"y,notnull"`
	Url           string                 `bun:"url,notnull"`
	PHash         string                 `bun:"phash,notnull"`
	Filename      string                 `bun:"filename,notnull"`
	MediaUrl      string                 `bun:"media_url,notnull"`
	ScreenshotUrl string                 `bun:"screenshot_url,notnull"`
	TaskId        string                 `bun:"task_id,notnull"`
	CreatedAt     time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt     bun.NullTime           `bun:"deleted_at"`
}
