package checker

import (
	"catscan-latex/structs"
	"github.com/dlclark/regexp2"
	"regexp"
)

var containsEtAl = regexp2.MustCompile(`et al\.`, 0)
var commaProceedsEtAl = regexp2.MustCompile(`,\s*(\\emph\{|\\textit\{|\{\\it\s*|\{\\em\s*)?et al`, 0)
var containsDoi = regexp2.MustCompile(`doi:10.`, 0)
var containsSpace = regexp2.MustCompile(`doi:\s10`, 0)
var noPrefix = regexp2.MustCompile(`\\url{10\.`, 0)
var doiIsUrl = regexp2.MustCompile(`https?://doi.org`, 0)

var wrappedEtAl = regexp.MustCompile(`(\\emph\{|\\textit\{|\{\\it\s*|\{\\em\s*)et al\.,?\s*}`)
var wrappedDoi = regexp.MustCompile(`\\url\{doi:10\.`)

var apsStyleReference = regexp2.MustCompile(`\d+[, -]+?\d+ \(\d{4}\)`, 0)
var brokenStyleReference = regexp2.MustCompile(`: N\. p\., \d{4}. Web\.`, 0)

func detectReferenceStyleReference(bibItem structs.BibItem) (bool, *structs.Location) {
	// Check if the reference is in APS style
	match, err := apsStyleReference.FindStringMatch(bibItem.Ref)
	if err == nil && match != nil {
		return true, &structs.Location{
			Start: match.Index + bibItem.Location.Start,
			End:   match.Index + match.Length + bibItem.Location.Start,
		}
	}

	// check for broken style that appears in some references ": N. p., \d{4}. Web."
	match, err = brokenStyleReference.FindStringMatch(bibItem.Ref)
	if err == nil && match != nil {
		return true, &structs.Location{
			Start: match.Index + bibItem.Location.Start,
			End:   match.Index + match.Length + bibItem.Location.Start,
		}
	}

	return false, nil
}

func etAlNotWrapped(bibItem structs.BibItem) (bool, *structs.Location) {
	// Check if et al. not wrapped in a \it \emph or \textit macro
	// eg. L. Kiani et al.,
	match, err := containsEtAl.FindStringMatch(bibItem.OriginalText)
	if err == nil && match != nil {
		isWrapped := wrappedEtAl.FindString(bibItem.OriginalText)
		if isWrapped == "" {
			location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
			return true, &location
		}
	}
	return false, nil
}

func CheckBibItem(bibItem structs.BibItem) []structs.Issue {
	var issues []structs.Issue

	// et al. should not be proceeded by a comma
	// eg. L. Kiani et al.,
	match, err := commaProceedsEtAl.FindStringMatch(bibItem.OriginalText)
	if err == nil && match != nil {
		location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
		issues = append(issues, structs.Issue{Type: "ET_AL_WITH_COMMA", Location: location})
	}

	if found, location := etAlNotWrapped(bibItem); found {
		issues = append(issues, structs.Issue{Type: "ET_AL_NOT_WRAPPED", Location: *location})
	}

	// Check that the doi does not contain a space after the colon
	// e.g. doi: 10.1000/182
	match, err = containsSpace.FindStringMatch(bibItem.OriginalText)
	if err == nil && match != nil {
		location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
		issues = append(issues, structs.Issue{Type: "DOI_CONTAINS_SPACE", Location: location})
	}

	// check if reference is using incorrect reference style
	if found, location := detectReferenceStyleReference(bibItem); found {
		issues = append(issues, structs.Issue{Type: "INCORRECT_STYLE_REFERENCE", Location: *location})
	}

	// Is wrapped in a URL macro
	// e.g. doi:10.1000/182 (without \url{})
	match, err = containsDoi.FindStringMatch(bibItem.OriginalText)
	if err == nil && match != nil {
		isWrapped := wrappedDoi.FindString(bibItem.OriginalText)
		if isWrapped == "" {
			location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
			issues = append(issues, structs.Issue{Type: "DOI_NOT_WRAPPED", Location: location})
		}
	}

	// Check that doi has a doi: prefix
	// e.g. \url{10.1000/182}
	match, err = noPrefix.FindStringMatch(bibItem.OriginalText)
	if err == nil && match != nil {
		location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
		issues = append(issues, structs.Issue{Type: "NO_DOI_PREFIX", Location: location})
	}

	// Check that DOI is not a http link
	// e.g. \url{https://doi.org/10.1000/182}
	match, err = doiIsUrl.FindStringMatch(bibItem.OriginalText)
	if err == nil && match != nil {
		location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
		issues = append(issues, structs.Issue{Type: "DOI_IS_URL", Location: location})
	}

	return issues
}

func FindBibItemIssues(bibItems []structs.BibItem) []structs.Issue {
	var issues []structs.Issue

	for _, bibItem := range bibItems {
		issues = append(issues, CheckBibItem(bibItem)...)
	}

	return issues
}
