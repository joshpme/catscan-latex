package main

import (
	"catscan-latex/checker"
	"catscan-latex/finder"
	"catscan-latex/structs"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/rs/cors"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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

func getResult(fileName string, contents string) (*structs.Contents, error) {
	result := Finder(structs.Request{Content: contents, Filename: fileName})
	return &result, nil
}

func issueToDescription(issue structs.Issue) string {
	switch issue.Type {
	case "INCORRECT_STYLE_REFERENCE":
		return "Reference does not appear to be in the JACoW style, please adjust your reference style to be consistent with the JACoW style reference, please see https://www.jacow.org/Authors/FormattingCitations"
	case "ET_AL_WITH_COMMA":
		return "et al. is preceded by a comma, which is incorrect. Please remove the comma before the et al."
	case "ET_AL_NOT_WRAPPED":
		return "et al. is not wrapped in a macro to make it italic. Please use \\emph{et al.} instead of et al."
	case "DOI_CONTAINS_SPACE":
		return "DOI contains a space after the colon. Please remove the space."
	case "DOI_NOT_WRAPPED":
		return "DOI not wrapped in \\url{} macro. Please use \\url{doi:10.18429/JACoW-IPAC2023-XXXX} instead of doi:10.18429/JACoW-IPAC2023-XXXX"
	case "NO_DOI_PREFIX":
		return "DOI does not contain \"doi:\" prefix. It should appear like this \\url{doi:10.18429/JACoW-IPAC2023-XXXX}"
	case "DOI_IS_URL":
		return "DOI is written as a web URL (including https://doi.org/) which is incorrect. Remove the https://doi.org/, and write it as per this example. \\url{doi:10.18429/JACoW-IPAC2023-XXXX}"
	case "VOLUME_ISSUE":
		return "JACoW references use vol. X and no. X. You have used not Vol. X, Issue X, which is incorrect. Please correct your reference style. You can generate correctly formatted references at https://refs.jacow.org/ or you can refer to the JACoW reference style guide at https://www.jacow.org/Authors/FormattingCitations"
	}
	return ""
}

type Request struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type Response struct {
	StatusCode    int               `json:"statusCode,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Body          string            `json:"body,omitempty"`
	IsAbbreviated bool              `json:"isabbreviated"`
	IssuesFound   int               `json:"issuesFound"`
	Unabbreviated string            `json:"unabbreviated"`
}

func geminiSummarize(content string) (string, error) {
	apiKey := os.Getenv("GEMINI_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_KEY environment variable not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("models/gemini-2.0-flash") // Or "gemini-pro"

	prompt := "You are an editor correcting bibitem references in a latex paper for a scientific conference.\n\n"
	prompt += "The text to be summerized is:\n"
	prompt += content + "\n\n"
	prompt += "Provide a summary of the issues. Do not include any introductory or concluding text.\n"

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("error generating content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}

	part := resp.Candidates[0].Content.Parts[0]
	switch part := part.(type) {
	case genai.Text:
		return string(part), nil
	}
	return "", fmt.Errorf("unexpected part type: %T", part)
}

type Report struct {
	issueFound    bool
	issueCount    int
	output        string
	unabbreviated string
}

func getReport(result structs.Contents) Report {
	report := Report{
		issueFound:    false,
		issueCount:    0,
		output:        "No issues found",
		unabbreviated: "",
	}
	for _, bibItem := range result.BibItems {
		issues := checker.CheckBibItem(bibItem)
		if len(issues) > 0 {
			report.issueFound = true
			report.unabbreviated += fmt.Sprintf("\nIssue found in reference %s:\n%s\n", strings.Trim(bibItem.Name, " \t\r\n"), strings.Trim(bibItem.Ref, " \t\n"))
			for _, issue := range issues {
				descriptionOfIssue := issueToDescription(issue)
				if descriptionOfIssue != "" {
					report.issueFound = true
					report.issueCount += 1
					report.unabbreviated += fmt.Sprintf(" %s\n", descriptionOfIssue)
				}
			}
		}

		doiResult, suggestion := checker.CheckDOIExists(bibItem)
		if doiResult == structs.HasIssue {
			report.issueFound = true
			if suggestion != nil {
				report.unabbreviated += fmt.Sprintf("\nIssue found in reference DOI for reference %s:\n%s\n", strings.Trim(bibItem.Name, " \t\r\n"), strings.Trim(bibItem.Ref, " \t\n"))
				report.unabbreviated += fmt.Sprintf("%s\n", suggestion.Description)
				report.issueCount += 1
				if suggestion.Content != "" {
					report.unabbreviated += fmt.Sprintf("Suggested DOI: %s\n", suggestion.Content)
				}
			}
		}
	}
	return report
}

func Main(in Request) (*Response, error) {
	fileName := in.Filename
	contents := in.Content
	isAbbreviated := false
	// for each file, read the contents and run the main function
	result, err := getResult(fileName, contents)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	report := getReport(*result)

	if report.issueFound {
		report.output = report.unabbreviated
		if report.issueCount > 3 {
			isAbbreviated = true
			geminiOutput, err := geminiSummarize(report.unabbreviated)
			if err == nil {
				report.output = geminiOutput
			}
		}
	}

	return &Response{
		StatusCode:    200,
		Body:          report.output,
		IsAbbreviated: isAbbreviated,
		IssuesFound:   report.issueCount,
		Unabbreviated: report.unabbreviated,
	}, nil
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Log the incoming request
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON body
	var req Request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Call the Main function
	resp, err := Main(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func getContents(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	files := findFiles("examples")
	reports := make([]Report, 0)
	var result *structs.Contents
	for _, fileName := range files {
		contents, err := getContents(fileName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading file: %s", fileName), http.StatusInternalServerError)
			return
		}
		result, err = getResult(fileName, contents)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting result: %s", fileName), http.StatusInternalServerError)
			break
		}
		report := getReport(*result)
		reports = append(reports, report)
	}
	resp := Response{
		StatusCode: 200,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", baseHandler)
	mux.HandleFunc("/dry-run", exampleHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Or specifically list your frontend domains
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}).Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	bindAddr := fmt.Sprintf(":%s", port)

	log.Printf("Starting server on %s", bindAddr)
	if err := http.ListenAndServe(bindAddr, corsHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
