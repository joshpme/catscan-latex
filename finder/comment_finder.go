package finder

import (
	"catscan-latex/structs"
	"github.com/dlclark/regexp2"
)

var commentRegex = regexp2.MustCompile(`%.*?\n`, regexp2.Singleline)

func FindComments(contents string) []structs.Comment {
	match, err := commentRegex.FindStringMatch(contents)
	comments := make([]structs.Comment, 0)
	for err == nil && match != nil {
		comments = append(comments, structs.Comment{
			Location: structs.Location{
				Start: match.Index,
				End:   match.Index + match.Length,
			},
		})
		match, err = commentRegex.FindNextMatch(match)
	}
	return comments
}

func locationInComments(location structs.Location, comments []structs.Comment) bool {
	for _, comment := range comments {
		if structs.LocationIn(location, comment.Location) {
			return true
		}
	}
	return false
}
