package checker

import (
	"catscan-latex/structs"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func checkDOIExists(doi string) (bool, error) {
	doiLookupURL := "https://doi.org/" + doi

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Do not follow redirects
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Head(doiLookupURL)

	if err != nil {
		return false, fmt.Errorf("failed to check DOI: %w", err)
	}
	defer resp.Body.Close()

	// Consider DOI valid if status is 200 OK or a redirect (3xx)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true, nil
	}

	return false, nil
}

func tryTrimDOI(originalDOI string, cutset string) (string, error) {
	trimmedDOI := strings.TrimRight(originalDOI, cutset)
	if trimmedDOI != originalDOI {
		exists, err := checkDOIExists(trimmedDOI)
		if err != nil {
			return "", fmt.Errorf("error checking DOI with no %s: %w", cutset, err)
		}
		if exists {
			return trimmedDOI, nil
		}
	}
	return "", nil
}

func CheckDOIExists(bibItem structs.BibItem) *structs.Issue {
	currentDOI := bibItem.Doi
	if currentDOI == "" {
		return nil
	}
	doiExists, err := checkDOIExists(currentDOI)
	if err != nil {
		log.Printf("Error checking exists DOI %s: %v", bibItem.Doi, err)
		return nil
	}
	if !doiExists {
		newDOI, err := tryTrimDOI(currentDOI, ".")
		if err != nil {
			log.Printf("Error checked trimmed DOI %s: %v", currentDOI, err)
			return nil
		}
		if newDOI != "" {
			return &structs.Issue{
				Name:       bibItem.Name,
				Location:   bibItem.Location,
				Type:       "DOI_ENDS_IN_PERIOD",
				Suggestion: newDOI,
			}
		}
	}
	if !doiExists {
		newDOI, err := tryTrimDOI(currentDOI, ")")
		if err != nil {
			log.Printf("Error checked trimmed DOI %s: %v", currentDOI, err)
			return nil
		}
		if newDOI != "" {
			return &structs.Issue{
				Name:       bibItem.Name,
				Location:   bibItem.Location,
				Type:       "DOI_ENDS_IN_PARENTHESIS",
				Suggestion: newDOI,
			}
		}
	}

	if !doiExists {
		return &structs.Issue{
			Name:     bibItem.Name,
			Location: bibItem.Location,
			Type:     "DOI_NOT_FOUND",
		}
	}

	return nil
}
