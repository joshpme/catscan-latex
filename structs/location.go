package structs

type Location struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func LocationIn(needle Location, haystack Location) bool {
	return needle.Start >= haystack.Start && needle.End <= haystack.End
}
