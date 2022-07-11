package types

import (
	"reflect"
	"testing"
)

func TestSurvey_BuildIndex(t *testing.T) {
	dataNodes := []string{"N01.01.", "N01.01.", "N01.01.01.", "N01.01.02.", "N01.02.01.02."}
	orgNodes := OrgStructure{nodes: []string{"N01.", "N01.01.", "N01.01.01.", "N01.01.02.", "N01.02.", "N01.02.01.01.", "N01.02.01.02."}}
	got := buildIndex(orgNodes, dataNodes).data
	want := map[string]loc{
		"N01.":          loc{0, 4, -1, -1},
		"N01.01":        loc{0, 3, 0, 1},
		"N01.01.01.":    loc{2, 3, 2, 2},
		"N01.01.02.":    loc{3, 3, 3, 3},
		"N01.02.":       loc{4, 4, -1, -1},
		"N01.02.01.01.": loc{-1, -1, -1, -1},
		"N01.02.01.02.": loc{4, 4, 4, 4},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %v, expected %v", got, want)
	}
}
