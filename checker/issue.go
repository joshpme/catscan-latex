package checker

import "catscan-latex/structs"

func GetIssues(result structs.Contents) []structs.Issue {
	issues := make([]structs.Issue, 0)
	for _, bibItem := range result.BibItems {
		bibItemIssues := CheckBibItem(bibItem)
		issues = append(issues, bibItemIssues...)
		issue := CheckDOIExists(bibItem)
		if issue != nil {
			issues = append(issues, *issue)
		}
	}
	return issues
}
