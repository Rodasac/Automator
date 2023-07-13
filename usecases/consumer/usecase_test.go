package consumer

import (
	"errors"
	"testing"
)

type mockTaskQueueConsumerAdapter struct {
	Err error
}

func (m *mockTaskQueueConsumerAdapter) ConsumeTasks() error {
	return m.Err
}

func TestStartConsumer(t *testing.T) {
	tests := []struct {
		name    string
		adapter TaskQueueConsumerAdapter
		wantErr bool
	}{
		{
			name:    "Test consumer without error",
			adapter: &mockTaskQueueConsumerAdapter{},
			wantErr: false,
		},
		{
			name:    "Test consumer with error",
			adapter: &mockTaskQueueConsumerAdapter{Err: errors.New("error")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tqc := NewTaskQueueConsumer(tt.adapter)
			if err := tqc.StartConsumer(); (err != nil) != tt.wantErr {
				t.Errorf("StartConsumer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
