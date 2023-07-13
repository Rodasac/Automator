package models

import "encoding/json"

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

func (a *Action) FromString(s string) Action {
	switch s {
	case "NavigateXpath":
		return NavigateXpath
	case "NavigateSelector":
		return NavigateSelector
	case "ClickXpath":
		return ClickXpath
	case "ClickSelector":
		return ClickSelector
	case "ScrollDown":
		return ScrollDown
	case "CaptureXpath":
		return CaptureXpath
	case "CaptureSelector":
		return CaptureSelector
	default:
		return NavigateXpath
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

	*a = a.FromString(s)
	return nil
}

type TaskAction struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Type  Action `json:"type"`
	Value string `json:"value"`
}
