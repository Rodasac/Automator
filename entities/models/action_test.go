package models

import (
	"encoding/json"
	"testing"
)

func TestAction_String(t *testing.T) {
	tests := []struct {
		name string
		a    Action
	}{
		{
			name: "NavigateXpath",
			a:    NavigateXpath,
		},
		{
			name: "NavigateSelector",
			a:    NavigateSelector,
		},
		{
			name: "ClickXpath",
			a:    ClickXpath,
		},
		{
			name: "ClickSelector",
			a:    ClickSelector,
		},
		{
			name: "ScrollDown",
			a:    ScrollDown,
		},
		{
			name: "CaptureXpath",
			a:    CaptureXpath,
		},
		{
			name: "CaptureSelector",
			a:    CaptureSelector,
		},
		{
			name: "WaitSeconds",
			a:    WaitSeconds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.name {
				t.Errorf("Action.String() = %v, want %v", got, tt.name)
			}
			if got, _ := tt.a.FromString(tt.name); got != tt.a {
				t.Errorf("Action.FromString() = %v, want %v", got, tt.a)
			}
		})
	}
}

func TestAction_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		a       Action
		wantErr bool
	}{
		{
			name:    "NavigateXpath",
			a:       NavigateXpath,
			wantErr: false,
		},
		{
			name:    "NavigateSelector",
			a:       NavigateSelector,
			wantErr: false,
		},
		{
			name:    "ClickXpath",
			a:       ClickXpath,
			wantErr: false,
		},
		{
			name:    "ClickSelector",
			a:       ClickSelector,
			wantErr: false,
		},
		{
			name:    "ScrollDown",
			a:       ScrollDown,
			wantErr: false,
		},
		{
			name:    "CaptureXpath",
			a:       CaptureXpath,
			wantErr: false,
		},
		{
			name:    "CaptureSelector",
			a:       CaptureSelector,
			wantErr: false,
		},
		{
			name:    "WaitSeconds",
			a:       WaitSeconds,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.a.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Action.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && string(b) != "\""+tt.name+"\"" {
				t.Errorf("Action.MarshalJSON() = %v, want %v", string(b), tt.name)
			}
		})
	}
}

func TestAction_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		value   []byte
		wantErr bool
	}{
		{
			name:    "NavigateXpath",
			value:   []byte("\"NavigateXpath\""),
			wantErr: false,
		},
		{
			name:    "NavigateSelector",
			value:   []byte("\"NavigateSelector\""),
			wantErr: false,
		},
		{
			name:    "ClickXpath",
			value:   []byte("\"ClickXpath\""),
			wantErr: false,
		},
		{
			name:    "ClickSelector",
			value:   []byte("\"ClickSelector\""),
			wantErr: false,
		},
		{
			name:    "ScrollDown",
			value:   []byte("\"ScrollDown\""),
			wantErr: false,
		},
		{
			name:    "CaptureXpath",
			value:   []byte("\"CaptureXpath\""),
			wantErr: false,
		},
		{
			name:    "CaptureSelector",
			value:   []byte("\"CaptureSelector\""),
			wantErr: false,
		},
		{
			name:    "WaitSeconds",
			value:   []byte("\"WaitSeconds\""),
			wantErr: false,
		},
		{
			name:    "Invalid",
			value:   []byte("\"Invalid\""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a Action
			err := json.Unmarshal(tt.value, &a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Action.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && a.String() != tt.name {
				t.Errorf("Action.MarshalJSON() = %v, want %v", a.String(), tt.name)
			}
		})
	}
}
