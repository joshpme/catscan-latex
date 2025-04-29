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
			name: "et al. not wrapped",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: "J. S. Berg et al. ``Lattice Design for the Hadron Storage Ring of the Electron-Ion Collider'', presented at IPAC'23, Venice, Italy, May 2023, paper MOPL156, this conference.",
					Ref:          "J. S. Berg et al. ``Lattice Design for the Hadron Storage Ring of the Electron-Ion Collider'', presented at IPAC'23, Venice, Italy, May 2023, paper MOPL156, this conference.",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: true,
		},
		{
			name: "et al. not wrapped",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: "J. S. Berg \\emph{et al.,} ``Lattice Design for the Hadron Storage Ring of the Electron-Ion Collider'', presented at IPAC'23, Venice, Italy, May 2023, paper MOPL156, this conference.",
					Ref:          "J. S. Berg \\emph{et al.,} ``Lattice Design for the Hadron Storage Ring of the Electron-Ion Collider'', presented at IPAC'23, Venice, Italy, May 2023, paper MOPL156, this conference.",
					Location: structs.Location{
						Start: 0,
						End:   0,
					},
				},
			},
			want: false,
		},
		{
			name: "et al. not wrapped",
			args: args{
				bibItem: structs.BibItem{
					OriginalText: "J. S. Berg {\\it et al.} ``Lattice Design for the Hadron Storage Ring of the Electron-Ion Collider'', presented at IPAC'23, Venice, Italy, May 2023, paper MOPL156, this conference.",
					Ref:          "J. S. Berg {\\it et al.} ``Lattice Design for the Hadron Storage Ring of the Electron-Ion Collider'', presented at IPAC'23, Venice, Italy, May 2023, paper MOPL156, this conference.",
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
			found, _ := etAlNotWrapped(tt.args.bibItem)
			if found != tt.want {
				t.Errorf("etAlNotWrapped() got = %v, want %v", found, tt.want)
			}
		})
	}
}
