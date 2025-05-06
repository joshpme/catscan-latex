package finder

import (
	"catscan-latex/structs"
	"testing"
)

func TestFindDocument(t *testing.T) {
	tests := []struct {
		name     string
		contents string
		comments []structs.Comment
		expected structs.Document
	}{
		{
			name: "Commented \\begin{document}",
			contents: `%\begin{document}
\begin{document}`,
			comments: []structs.Comment{
				{Location: structs.Location{Start: 0, End: 17}}, // Commented lines
			},
			expected: structs.Document{
				Location: structs.Location{Start: 18, End: 33}, // Actual document location
			},
		},
		{
			name: "Commented \\begin{document}",
			contents: `\begin{document}
%\begin{document}`,
			comments: []structs.Comment{
				{Location: structs.Location{Start: 17, End: 33}}, // Commented lines
			},
			expected: structs.Document{
				Location: structs.Location{Start: 0, End: 33}, // Actual document location
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindDocument(tt.contents, tt.comments)
			if got.Location != tt.expected.Location {
				t.Errorf("FindDocument() = %v, want %v", got.Location, tt.expected.Location)
			}
		})
	}
}
