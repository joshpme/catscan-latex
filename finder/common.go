package finder

import "catscan-latex/structs"

func Finder(in structs.Request) structs.Contents {
	filename := in.Filename
	contents := in.Content
	comments := FindComments(contents)
	document := FindDocument(contents, comments)
	bibItems := FindValidBibItems(contents, comments, document)
	return structs.Contents{
		Document: document,
		Comments: comments,
		BibItems: bibItems,
		Filename: filename,
		Content:  contents,
	}
}
