package finder

import (
	"github.com/dlclark/regexp2"
	"latex/structs"
)

var abstractBeginRegex = regexp2.MustCompile(`\\begin{abstract}`, 0)
var abstractEndRegex = regexp2.MustCompile(`\\end{abstract}`, 0)

func FindAbstractLocation(contents string, document structs.Document, comments []structs.Comment) *structs.Location {
	match, err := abstractBeginRegex.FindStringMatch(contents)
	start := 0
	foundStart := false
	for err == nil && match != nil {
		location := structs.Location{Start: match.Index, End: match.Index + match.Length}
		if structs.LocationIn(location, document.Location) && !locationInComments(location, comments) {
			start = location.Start
			foundStart = true
			break
		}
		match, err = abstractBeginRegex.FindNextMatch(match)
	}

	match, err = abstractEndRegex.FindStringMatch(contents)
	end := len(contents) - 1
	foundEnd := false
	for err == nil && match != nil {
		location := structs.Location{Start: match.Index, End: match.Index + match.Length}
		if structs.LocationIn(location, document.Location) && !locationInComments(location, comments) {
			end = location.Start
			foundEnd = true
			break
		}
		match, err = abstractEndRegex.FindNextMatch(match)
	}

	if !foundStart || !foundEnd {
		return nil
	}
	return &structs.Location{Start: start, End: end}
}
