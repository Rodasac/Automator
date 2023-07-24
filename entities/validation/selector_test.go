package validation

import "testing"

func TestIsXpath(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "Xpath",
			s:    "//div[@id='test']",
			want: true,
		},
		{
			name: "Selector",
			s:    "div#test",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsXpath(tt.s); got != tt.want {
				t.Errorf("IsXpath() = %v, want %v", got, tt.want)
			}
		})
	}
}
