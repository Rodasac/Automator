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
			name: "Navigate",
			a:    Navigate,
		},
		{
			name: "Click",
			a:    Click,
		},
		{
			name: "ScrollDown",
			a:    ScrollDown,
		},
		{
			name: "Capture",
			a:    Capture,
		},
		{
			name: "WaitSeconds",
			a:    WaitSeconds,
		},
		{
			name: "WriteInput",
			a:    WriteInput,
		},
		{
			name: "ClearInput",
			a:    ClearInput,
		},
		{
			name: "SelectOptions",
			a:    SelectOptions,
		},
		{
			name: "WriteTime",
			a:    WriteTime,
		},
		{
			name: "DownloadResource",
			a:    DownloadResource,
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
			name:    "Navigate",
			a:       Navigate,
			wantErr: false,
		},
		{
			name:    "Click",
			a:       Click,
			wantErr: false,
		},
		{
			name:    "ScrollDown",
			a:       ScrollDown,
			wantErr: false,
		},
		{
			name:    "Capture",
			a:       Capture,
			wantErr: false,
		},
		{
			name:    "WaitSeconds",
			a:       WaitSeconds,
			wantErr: false,
		},
		{
			name:    "WriteInput",
			a:       WriteInput,
			wantErr: false,
		},
		{
			name:    "ClearInput",
			a:       ClearInput,
			wantErr: false,
		},
		{
			name:    "SelectOptions",
			a:       SelectOptions,
			wantErr: false,
		},
		{
			name:    "WriteTime",
			a:       WriteTime,
			wantErr: false,
		},
		{
			name:    "DownloadResource",
			a:       DownloadResource,
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
			name:    "Navigate",
			value:   []byte("\"Navigate\""),
			wantErr: false,
		},
		{
			name:    "Click",
			value:   []byte("\"Click\""),
			wantErr: false,
		},
		{
			name:    "ScrollDown",
			value:   []byte("\"ScrollDown\""),
			wantErr: false,
		},
		{
			name:    "Capture",
			value:   []byte("\"Capture\""),
			wantErr: false,
		},
		{
			name:    "WaitSeconds",
			value:   []byte("\"WaitSeconds\""),
			wantErr: false,
		},
		{
			name:    "WriteInput",
			value:   []byte("\"WriteInput\""),
			wantErr: false,
		},
		{
			name:    "ClearInput",
			value:   []byte("\"ClearInput\""),
			wantErr: false,
		},
		{
			name:    "SelectOptions",
			value:   []byte("\"SelectOptions\""),
			wantErr: false,
		},
		{
			name:    "WriteTime",
			value:   []byte("\"WriteTime\""),
			wantErr: false,
		},
		{
			name:    "DownloadResource",
			value:   []byte("\"DownloadResource\""),
			wantErr: false,
		},
		{
			name:    "Invalid",
			value:   []byte("\"Invalid\""),
			wantErr: true,
		},
		{
			name:    "Empty",
			value:   []byte("\"\""),
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
