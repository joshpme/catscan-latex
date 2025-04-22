package main

import (
	"fmt"
	"io"
	"latex/checker"
	"latex/finder"
	"latex/structs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetContents(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	contents, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func Finder(in structs.Request) structs.Contents {
	filename := in.Filename
	contents := in.Content
	comments := finder.FindComments(contents)
	document := finder.FindDocument(contents, comments)
	bibItems := finder.FindValidBibItems(contents, comments, document)
	return structs.Contents{
		Document: document,
		Comments: comments,
		BibItems: bibItems,
		Filename: filename,
		Content:  contents,
	}
}

func findFiles(directory string) []string {
	files := make([]string, 0)
	entries, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("failed reading directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".tex" {
			files = append(files, filepath.Join(directory, entry.Name()))
		}
	}
	return files
}

func getResult(fileName string) (*structs.Contents, error) {
	contents, err := GetContents(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	result := Finder(structs.Request{Content: contents, Filename: fileName})
	return &result, nil
}

func getResults(fileNames []string) ([]structs.Contents, int) {
	results := make([]structs.Contents, 0)
	failures := 0
	for _, fileName := range fileNames {
		result, err := getResult(fileName)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", fileName, err)
			failures++
			continue
		}
		results = append(results, *result)
	}
	return results, failures
}

func issueToDescription(issue structs.Issue) string {
	switch issue.Type {
	case "INCORRECT_STYLE_REFERENCE":
		return "Reference not appear to be in the JACoW style, please adjust your reference style to the JACoW style reference, please see https://www.jacow.org/Authors/FormattingCitations"
	case "ET_AL_WITH_COMMA":
		return "et al. is preceded by a comma"
	case "ET_AL_NOT_WRAPPED":
		return "et al. is not wrapped in a macro to make it italic"
	case "DOI_CONTAINS_SPACE":
		return "DOI contains a space after the colon"
	case "DOI_NOT_WRAPPED":
		return "DOI not wrapped in \\url{} macro."
	case "NO_DOI_PREFIX":
		return "DOI does not contain \"doi:\" prefix. It should appear like this \\url{doi:10.18429/JACoW-IPAC2023-XXXX}"
	case "DOI_IS_URL":
		return "DOI is a URL, instead it should appear like this \\url{doi:10.18429/JACoW-IPAC2023-XXXX}"
	}
	return ""
}

func main() {

	// find all the files in the examples folder
	dirPath := "./examples"
	fileNames := findFiles(dirPath)
	if len(fileNames) == 0 {
		fmt.Println("No .tex files found in the directory.")
		return
	}

	// for each file, read the contents and run the main function
	results, failures := getResults(fileNames)
	if failures > 0 {
		fmt.Printf("Failed to process %d files.\n", failures)
	} else {
		fmt.Println("All files processed successfully.")
	}

	for _, result := range results {
		report := ""
		issueFound := false
		report += fmt.Sprintf("\n\n\nIn File: %s\n", result.Filename)
		issueCount := 0

		for _, bibItem := range result.BibItems {
			issues := checker.CheckBibItem(bibItem)
			if len(issues) > 0 {
				issueFound = true
				report += fmt.Sprintf("\nIssue found in reference:\n%s\n", strings.Trim(bibItem.Ref, " \t\n"))
				for _, issue := range issues {
					issueCount += 1
					descriptionOfIssue := issueToDescription(issue)
					if descriptionOfIssue != "" {
						issueFound = true
						report += fmt.Sprintf("[%d]: %s\n", issueCount, descriptionOfIssue)
					} else {
						report += fmt.Sprintf("[%d]: Unknown issue type: %s\n", issueCount, issue.Type)
					}
				}
			}

			doiResult, suggestion := checker.CheckDOIExists(bibItem)
			if doiResult == structs.HasIssue {
				issueFound = true
				if suggestion != nil {
					issueCount += 1
					report += fmt.Sprintf("\nIssue found in reference DOI for reference:\n%s\n", strings.Trim(bibItem.Ref, " \t\n"))
					report += fmt.Sprintf("[%d] %s\n", issueCount, suggestion.Description)
					report += fmt.Sprintf("Suggested DOI: %s\n", suggestion.Content)
				}
			}
		}
		if issueFound {
			fmt.Println(report)
		}
	}
}
