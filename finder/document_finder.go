package finder

import (
	"github.com/dlclark/regexp2"
	"latex/structs"
)

var documentBeginRegex = regexp2.MustCompile(`\\begin{document}`, 0)
var documentEndRegex = regexp2.MustCompile(`\\end{document}`, 0)

func FindDocument(contents string, comments []structs.Comment) structs.Document {
	match, err := documentBeginRegex.FindStringMatch(contents)
	start := 0
	for err == nil && match != nil {
		location := structs.Location{Start: match.Index, End: match.Index + match.Length}
		if !locationInComments(location, comments) {
			start = location.Start
			break
		}
		match, err = documentBeginRegex.FindNextMatch(match)
	}

	match, err = documentEndRegex.FindStringMatch(contents)
	end := len(contents) - 1
	for err == nil && match != nil {
		location := structs.Location{Start: match.Index, End: match.Index + match.Length}
		if !locationInComments(location, comments) {
			end = location.Start
			break
		}
		match, err = documentEndRegex.FindNextMatch(match)
	}

	return structs.Document{Location: structs.Location{Start: start, End: end}}
}
