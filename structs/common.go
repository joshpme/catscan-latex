package structs

type Document struct {
	Location Location `json:"location"`
}

type Comment struct {
	Location Location `json:"location"`
}

type Citation struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
}

type BibItem struct {
	Name          string   `json:"-"`
	OriginalText  string   `json:"-"`
	Doi           string   `json:"doi"`
	Ref           string   `json:"ref"`
	Location      Location `json:"location"`
	LabelLocation Location `json:"labelLocation"`
}

type Issue struct {
	Type     string   `json:"type"`
	Location Location `json:"location"`
}

type CheckResult int

const (
	NoIssue CheckResult = iota
	HasIssue
	NoSure
)

type Suggestion struct {
	Description string `json:"description"`
	Content     string `json:"content"`
}

type Request struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type Contents struct {
	Filename string    `json:"filename"`
	Content  string    `json:"content"`
	BibItems []BibItem `json:"bibItems"`
	Document Document  `json:"-"`
	Comments []Comment `json:"-"`
}
