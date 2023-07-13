package models

import (
	"encoding/json"
	"fmt"
)

type Action uint8

const (
	NavigateXpath Action = iota
	NavigateSelector
	ClickXpath
	ClickSelector
	ScrollDown
	CaptureXpath
	CaptureSelector
)

func (a *Action) String() string {
	return [...]string{
		"NavigateXpath",
		"NavigateSelector",
		"ClickXpath",
		"ClickSelector",
		"ScrollDown",
		"CaptureXpath",
		"CaptureSelector",
	}[*a]
}

func (a *Action) FromString(s string) (Action, error) {
	switch s {
	case "NavigateXpath":
		return NavigateXpath, nil
	case "NavigateSelector":
		return NavigateSelector, nil
	case "ClickXpath":
		return ClickXpath, nil
	case "ClickSelector":
		return ClickSelector, nil
	case "ScrollDown":
		return ScrollDown, nil
	case "CaptureXpath":
		return CaptureXpath, nil
	case "CaptureSelector":
		return CaptureSelector, nil
	default:
		return NavigateXpath, fmt.Errorf("invalid action %s", s)
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
