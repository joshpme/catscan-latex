package finder

import (
	"catscan-latex/structs"
	"github.com/dlclark/regexp2"
	"strings"
)

func removeExcessWhitespace(contents string) string {
	var whitespaceRegex = regexp2.MustCompile(`\s+`, regexp2.Singleline)
	var withoutMoreThanOneSpace, _ = whitespaceRegex.Replace(contents, " ", 0, -1)
	return strings.Trim(withoutMoreThanOneSpace, " \n\t")
}

var bibItemRegex = regexp2.MustCompile(`(\\bibitem\{(.*?)})(.*?)(?=(\\bibitem|\\end\{thebibliography}))`, regexp2.Singleline)

func findBibItems(contents string) []structs.BibItem {
	var items []structs.BibItem
	match, err := bibItemRegex.FindStringMatch(contents)
	for err == nil && match != nil {
		items = append(items, structs.BibItem{
			Name:         match.Groups()[2].String(),
			Ref:          removeExcessWhitespace(match.Groups()[2].String()),
			OriginalText: match.Groups()[3].String(),
			Location: structs.Location{
				Start: match.Groups()[3].Index,
				End:   match.Groups()[3].Index + match.Groups()[2].Length,
			},
			LabelLocation: structs.Location{
				Start: match.Groups()[1].Index,
				End:   match.Groups()[1].Index + match.Groups()[1].Length,
			},
		})
		match, err = bibItemRegex.FindNextMatch(match)
	}
	return items
}

func filterBibItemsInComments(references []structs.BibItem, comments []structs.Comment) []structs.BibItem {
	var filtered []structs.BibItem
	for _, ref := range references {
		if !locationInComments(ref.LabelLocation, comments) {
			filtered = append(filtered, ref)
		}
	}
	return filtered
}

func filterBibItemInDocument(references []structs.BibItem, document structs.Document) []structs.BibItem {
	var filtered []structs.BibItem
	for _, ref := range references {
		if !structs.LocationIn(ref.Location, document.Location) {
			filtered = append(filtered, ref)
		}
	}
	return filtered
}

var doiRegex = regexp2.MustCompile(`10\.\d{4,9}/[-._\\;()/:a-zA-Z0-9]+`, regexp2.Singleline)

func findLastDoi(reference string) string {
	var lastDoi string
	match, err := doiRegex.FindStringMatch(reference)
	for err == nil && match != nil {
		lastDoi = match.Groups()[0].String()
		match, err = doiRegex.FindNextMatch(match)
	}
	return lastDoi
}

func findDois(references []structs.BibItem) []structs.BibItem {
	var dois []structs.BibItem
	for _, ref := range references {
		dois = append(dois, structs.BibItem{
			Name:         ref.Name,
			Ref:          ref.Ref,
			OriginalText: ref.OriginalText,
			Location:     ref.Location,
			Doi:          findLastDoi(ref.Ref),
		})
	}
	return dois
}

func FindValidBibItems(contents string, comments []structs.Comment, document structs.Document) []structs.BibItem {
	all := findBibItems(contents)
	noCommented := filterBibItemsInComments(all, comments)
	inDocument := filterBibItemInDocument(noCommented, document)
	withDois := findDois(inDocument)
	return withDois
}
