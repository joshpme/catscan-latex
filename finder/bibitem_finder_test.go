package finder

import (
	"log"
	"os"
	"testing"
)

func Test_findBibItems(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		comments        int
		allRefs         int
		uncommentedRefs int
	}{
		{
			name:  "No bibitems",
			input: "example.tex",

			comments:        1,
			allRefs:         2,
			uncommentedRefs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// decode json string to string
			content, err := os.ReadFile(tt.input)
			if err != nil {
				log.Fatalf("Failed to read file: %v", err)
			}
			fileContent := string(content)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}
			comments := FindComments(fileContent)
			if got := len(comments); got != tt.comments {
				t.Errorf("FindComments() = %v, want %v", got, tt.comments)
			}
			all := findBibItems(fileContent)
			if got := len(all); got != tt.allRefs {
				t.Errorf("findBibItems() = %v, want %v", got, tt.allRefs)
			}
			filtered := filterBibItemsInComments(all, comments)
			if got := len(filtered); got != tt.uncommentedRefs {
				t.Errorf("filterBibItemsInComments() = %v, want %v", got, tt.uncommentedRefs)
			}
		})
	}
}

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
