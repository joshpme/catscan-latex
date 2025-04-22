package finder

import (
	"testing"
)

func Test_findLastDoi(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single DOI",
			input:    "This is a reference with DOI 10.1234/5678",
			expected: "10.1234/5678",
		},
		{
			name:     "Multiple DOIs - should return last one",
			input:    "D. P. Aguillard \\textit{et al.}, \"Measurement of the Positive Muon Anomalous Magnetic Moment to 0.20 ppm\", \\textit{Phys. Rev. Lett.}, vol. 131, no. 16, p. 161802, Oct. 2023. doi:10.1103/PhysRevLett.131.161802.",
			expected: "10.1103/PhysRevLett.131.161802.",
		},
		{
			name:     "No DOI",
			input:    "This is a reference without any DOI",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findLastDoi(tt.input); got != tt.expected {
				t.Errorf("findLastDoi() = %v, want %v", got, tt.expected)
			}
		})
	}
}
