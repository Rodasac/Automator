package models

import (
	"encoding/json"
	"fmt"
)

type Action uint8

const (
	Navigate Action = iota
	Click
	ScrollDown
	Capture
	WaitSeconds
	WriteInput
	SelectOptions
	WriteTime
	ClearInput
	DownloadResource
)

func (a *Action) String() string {
	return [...]string{
		"Navigate",
		"Click",
		"ScrollDown",
		"Capture",
		"WaitSeconds",
		"WriteInput",
		"SelectOptions",
		"WriteTime",
		"ClearInput",
		"DownloadResource",
	}[*a]
}

func (a *Action) FromString(s string) (Action, error) {
	switch s {
	case "Navigate":
		return Navigate, nil
	case "Click":
		return Click, nil
	case "ScrollDown":
		return ScrollDown, nil
	case "Capture":
		return Capture, nil
	case "WaitSeconds":
		return WaitSeconds, nil
	case "WriteInput":
		return WriteInput, nil
	case "ClearInput":
		return ClearInput, nil
	case "SelectOptions":
		return SelectOptions, nil
	case "WriteTime":
		return WriteTime, nil
	case "DownloadResource":
		return DownloadResource, nil
	default:
		return Navigate, fmt.Errorf("invalid action %s", s)
	}
}

func (a *Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	action, err := a.FromString(s)
	if err != nil {
		return err
	}

	*a = action
	return nil
}

type TaskAction struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Type  Action `json:"type"`
	Value string `json:"value"`
}
