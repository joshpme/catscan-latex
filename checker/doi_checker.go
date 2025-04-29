package checker

import (
	"catscan-latex/structs"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

func searchRefSearch(reference string) (string, error) {
	urlEncodedReference := url.QueryEscape(reference)
	jacowRefSearchURL := fmt.Sprintf("https://refs.jacow.org/internal?query=%s", urlEncodedReference)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(jacowRefSearchURL)
	if err != nil {
		return "", fmt.Errorf("failed to check DOI: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to check DOI: %s", resp.Status)
	}

	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	type RefSearchResult struct {
		Query []struct {
			DOI string `json:"doi"`
		} `json:"query"`
	}

	var result RefSearchResult
	err = json.Unmarshal(bodyContent, &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(result.Query) > 0 && result.Query[0].DOI != "" {
		return result.Query[0].DOI, nil
	}
	return "", nil
}

func searchDOI(reference string) (string, error) {
	if strings.Contains(reference, "10.18429") || strings.Contains(reference, "JACoW") {
		return searchRefSearch(reference)
	}
	return "", nil
}

func CheckDOIExists(bibItem structs.BibItem) (structs.CheckResult, *structs.Suggestion) {
	currentDOI := bibItem.Doi
	if currentDOI == "" {
		return structs.NoIssue, nil
	}
	doiExists, err := checkDOIExists(currentDOI)
	if err != nil {
		log.Printf("Error checking exists DOI %s: %v", bibItem.Doi, err)
		return structs.NoSure, nil
	}
	if !doiExists {
		newDOI, err := tryTrimDOI(currentDOI, ".")
		if err != nil {
			log.Printf("Error checked trimmed DOI %s: %v", currentDOI, err)
			return structs.NoSure, nil
		}
		if newDOI != "" {
			return structs.HasIssue, &structs.Suggestion{
				Content:     newDOI,
				Description: "This DOI ends with a period, which is not correct for this specific DOI. Please remove the period.",
			}
		}
	}
	if !doiExists {
		newDOI, err := tryTrimDOI(currentDOI, ")")
		if err != nil {
			log.Printf("Error checked trimmed DOI %s: %v", currentDOI, err)
			return structs.NoSure, nil
		}
		if newDOI != "" {
			return structs.HasIssue, &structs.Suggestion{
				Content:     newDOI,
				Description: "Do not use parathesis e.g () around your DOI. Instead wrap it with a \\url{} command.",
			}
		}
	}

	if !doiExists {
		return structs.HasIssue, &structs.Suggestion{
			Description: "The DOI does not appear to be valid. Please check the DOI.",
		}
	}

	return structs.NoIssue, nil
}
