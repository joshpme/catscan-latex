package main

import (
	"bytes"
	"catscan-latex/checker"
	"catscan-latex/finder"
	"catscan-latex/structs"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type detailEntry struct {
	FileName string
	Issues   []structs.Issue
}

func main() {
	files := findFiles("examples")
	details := make([]detailEntry, 0)
	for _, fileName := range files {
		fmt.Printf("Checking %s\n", fileName)
		contents, err := getContents(fileName)
		if err != nil {
			log.Fatalf("Error reading file '%v': %v", fileName, err)
		}

		result := finder.Finder(structs.Request{Content: contents, Filename: fileName})

		entry := detailEntry{
			FileName: fileName,
			Issues:   checker.GetIssues(result),
		}

		sort.Slice(entry.Issues, func(i, j int) bool {
			return entry.Issues[i].Name < entry.Issues[j].Name
		})

		details = append(details, entry)
	}

	sort.Slice(details, func(i, j int) bool {
		return details[i].FileName < details[j].FileName
	})

	outFileName := "stats/details.csv"
	err := writeDetailsFile(details, outFileName)
	if err != nil {
		log.Fatalf("Error writing details to file '%v': %v", outFileName, err)
	}

	log.Printf("%d details written to '%v'", len(details), outFileName)

	summaryFileName := "stats/summary.csv"
	err = writeSummaryToFile(details, summaryFileName)
	if err != nil {
		log.Fatalf("Error writing summary to file '%v': %v", summaryFileName, err)
	}
}

func writeDetailsFile(checksums []detailEntry, fileName string) error {
	var buf bytes.Buffer

	for _, entry := range checksums {
		for _, issue := range entry.Issues {
			buf.WriteString(fmt.Sprintf("%s,%v\n", entry.FileName, issue.Type))
		}
	}

	return os.WriteFile(fileName, buf.Bytes(), 0644)
}

func writeSummaryToFile(checksums []detailEntry, fileName string) error {
	m := make(map[string]int)
	for _, entry := range checksums {
		for _, issue := range entry.Issues {
			m[issue.Type]++
		}
	}

	issueTypes := make([]string, 0, len(m))
	for issueType := range m {
		issueTypes = append(issueTypes, issueType)
	}

	// Sort the issue types alphabetically
	sort.Strings(issueTypes)

	var buf bytes.Buffer
	for _, issueType := range issueTypes {
		count := m[issueType]
		buf.WriteString(fmt.Sprintf("%s,%d\n", issueType, count))
	}

	return os.WriteFile(fileName, buf.Bytes(), 0644)
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

func getContents(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
