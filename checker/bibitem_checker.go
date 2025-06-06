package checker

import (
	"catscan-latex/structs"
	"github.com/dlclark/regexp2"
	"regexp"
	"strings"
)

var noPrefix = regexp2.MustCompile(`\\url\s*{10\.`, 0)
var containsSpace = regexp2.MustCompile(`doi:\s10`, 0)

var doiIsUrl = regexp2.MustCompile(`https?://(dx\.)?doi.org`, 0)

var volumeIssue = regexp2.MustCompile(`Vol. \d+, Issue \d+,`, 0)
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

var containsEtAl = regexp2.MustCompile(`et al\.`, 0)
var validCommands = []string{
	`\\emph\s*{`,
	`\\textit\s*{`,
	`\\it\s*{`,
	`\\em\s*{`,
	`\\itshape\s*{`,
	`{\\it\s+`,
	`{\\em\s+`,
	`{\\itshape\s+`,
}
var validOptions = `(` + strings.Join(validCommands, "|") + ")"
var italicEtAl = regexp.MustCompile(validOptions + `\s*et al\.`)

func etAlNotItalic(bibItem structs.BibItem) (bool, *structs.Location) {
	// Check if et al. not wrapped in a \it \emph or \textit command
	// eg. L. Kiani et al.,
	match, err := containsEtAl.FindStringMatch(bibItem.OriginalText)
	if err == nil && match != nil {
		isWrapped := italicEtAl.FindString(bibItem.OriginalText)
		if isWrapped == "" {
			location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
			return true, &location
		}
	}
	return false, nil
}

var containsDoi = regexp2.MustCompile(`doi:\s?10.`, 0)
var wrappedDoi = regexp.MustCompile(`\\url\s*({|"|\||#|!|'})(doi:)?10\.`)

func detectContainsDoiNotWrappedInUrl(bibItem structs.BibItem) (bool, *structs.Location) {
	match, err := containsDoi.FindStringMatch(bibItem.Ref)
	if err == nil && match != nil {
		isWrapped := wrappedDoi.FindString(bibItem.Ref)
		if isWrapped == "" {
			location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
			return true, &location
		}
	}
	return false, nil
}

// return: issueCode, issueDetected
// func detectEtAlIssue(ref string) (string, bool) {

// 	findIfEtAlIsPresent := containsEtAl.FindString(ref)
// 	previousCharacterALetter := false
// 	for _, letter := range ref {

// 		if letter == '.' {
// 			authorInitialsFound := true
// 		}

// 		// prepare for next character
// 		previousCharacterALetter = false
// 		if unicode.IsLetter(letter) {
// 			previousCharacterALetter = true
// 		}
// 	}
// 	return "", false
// }

func CheckBibItem(bibItem structs.BibItem) []structs.Issue {
	var issues []structs.Issue

	// et al. should not be proceeded by a comma
	// eg. L. Kiani et al.,
	// if found, location := detectSingleAuthorCommaBeforeEtAl(bibItem); found {
	// 	issues = append(issues, structs.Issue{Name: bibItem.Name, Type: "ET_AL_WITH_COMMA", Location: *location})
	// }

	// if found, location := detectMultiAuthorNoCommaBeforeEtAl(bibItem); found {
	// 	issues = append(issues, structs.Issue{Name: bibItem.Name, Type: "ET_AL_WITHOUT_COMMA", Location: *location})
	// }

	if found, location := etAlNotItalic(bibItem); found {
		issues = append(issues, structs.Issue{Name: bibItem.Name, Type: "ET_AL_NOT_WRAPPED", Location: *location})
	}

	// Check that the doi does not contain a space after the colon
	// e.g. doi: 10.1000/182
	match, err := containsSpace.FindStringMatch(bibItem.Ref)
	if err == nil && match != nil {
		location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
		issues = append(issues, structs.Issue{Name: bibItem.Name, Type: "DOI_CONTAINS_SPACE", Location: location})
	}

	// check if reference is using incorrect reference style
	if found, location := detectReferenceStyleReference(bibItem); found {
		issues = append(issues, structs.Issue{Name: bibItem.Name, Type: "INCORRECT_STYLE_REFERENCE", Location: *location})
	}

	// Is wrapped in a URL command
	// e.g. doi:10.1000/182 (without \url{})
	// if found, location := detectContainsDoiNotWrappedInUrl(bibItem); found {
	// 	issues = append(issues, structs.Issue{Name: bibItem.Name, Type: "DOI_NOT_WRAPPED", Location: *location})
	// }

	// Check that doi has a doi: prefix
	// e.g. \url{10.1000/182}
	match, err = noPrefix.FindStringMatch(bibItem.Ref)
	if err == nil && match != nil {
		location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
		issues = append(issues, structs.Issue{Name: bibItem.Name, Type: "NO_DOI_PREFIX", Location: location})
	}

	// Check that DOI is not a http link
	// e.g. \url{https://doi.org/10.1000/182}
	match, err = doiIsUrl.FindStringMatch(bibItem.Ref)
	if err == nil && match != nil {
		location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
		issues = append(issues, structs.Issue{Name: bibItem.Name, Type: "DOI_IS_URL", Location: location})
	}

	match, err = volumeIssue.FindStringMatch(bibItem.Ref)
	if err == nil && match != nil {
		location := structs.Location{Start: match.Index + bibItem.Location.Start, End: match.Index + match.Length + bibItem.Location.Start}
		issues = append(issues, structs.Issue{Name: bibItem.Name, Type: "VOLUME_ISSUE", Location: location})
	}

	return issues
}
