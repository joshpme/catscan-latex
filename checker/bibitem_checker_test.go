package checker

import (
	"catscan-latex/structs"
	"testing"
)

func Test_etAlNotWrapped(t *testing.T) {
	type args struct {
		bibItem structs.BibItem
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no et al",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: "J. Smith, `Something else`",
					Ref:          "J. Smith, `Something else`",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "not wrapped",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " et al. ",
					Ref:          " et al. ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: true,
		},
		{
			name: "correct wrapped, but with space at end",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " \\textit{et al. } ",
					Ref:          " \\textit{et al. } ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "correct wrapped, but with space at start",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " \\textit{ et al.} ",
					Ref:          " \\textit{ et al.} ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "using \\textit command",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " \\textit{et al.} ",
					Ref:          " \\textit{et al.} ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "using \\emph command",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " \\emph{et al.,} ",
					Ref:          " \\emph{et al.,} ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "using \\it declaration inside group",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " {\\it et al.} ",
					Ref:          " {\\it et al.} ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "using \\em declaration inside group",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " {\\em et al.} ",
					Ref:          " {\\em et al.} ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "using \\em outside group",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " \\em{et al.} ",
					Ref:          " \\em{et al.} ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "using \\it outside group",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " \\it{et al.} ",
					Ref:          " \\it{et al.} ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "using \\em outside group",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " \\em{et al.} ",
					Ref:          " \\em{et al.} ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "using \\itshape declaration",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: " {\\itshape et al.} ",
					Ref:          " {\\itshape et al.} ",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, _ := etAlNotItalic(tt.args.bibItem)
			if found != tt.want {
				t.Errorf("etAlNotItalic() got = %v, want %v", found, tt.want)
			}
		})
	}
}

func Test_detectContainsDoiNotWrappedInUrl(t *testing.T) {
	type args struct {
		bibItem structs.BibItem
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 *structs.Location
	}{
		{
			name: "correctly wrapped",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: "\\url{doi:10.18429/JACoW-IPAC2023-TUPM055}",
					Ref:          "\\url{doi:10.18429/JACoW-IPAC2023-TUPM055}",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "wrapped using alternative delimiter",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: "\\url|doi:10.18429/JACoW-IPAC2023-TUPM055|",
					Ref:          "\\url|doi:10.18429/JACoW-IPAC2023-TUPM055|",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "correctly wrapped (with space)",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: "\\url {doi:10.18429/JACoW-IPAC2023-TUPM055}",
					Ref:          "\\url {doi:10.18429/JACoW-IPAC2023-TUPM055}",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "correctly wrapped (with space)",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: "\\url {10.18429/JACoW-IPAC2023-TUPM055}",
					Ref:          "\\url {10.18429/JACoW-IPAC2023-TUPM055}",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "not wrapped",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: "doi:10.18429/JACoW-IPAC2023-TUPM055",
					Ref:          "doi:10.18429/JACoW-IPAC2023-TUPM055",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := detectContainsDoiNotWrappedInUrl(tt.args.bibItem)
			if got != tt.want {
				t.Errorf("detectContainsDoiNotWrappedInUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}
