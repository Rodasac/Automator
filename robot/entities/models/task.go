package models

type Task struct {
	Id          string       `json:"id"`
	Title       string       `json:"name"`
	Description string       `json:"description"`
	Url         string       `json:"url"`
	Country     string       `json:"country"`
	WithProxy   bool         `json:"with_proxy"`
	Actions     []TaskAction `json:"actions"`
}
