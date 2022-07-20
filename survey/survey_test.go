package survey

import (
	"reflect"
	"testing"
)

func TestSurvey_BuildIndex(t *testing.T) {
	dataNodes := []string{
		"N01.01.",
		"N01.01.",
		"N01.01.01.",
		"N01.01.02.",
		"N01.02.01.02.",
	}
	orgNodes := OrgStructure{nodes: []string{
		"N01.",
		"N01.01.",
		"N01.01.01.",
		"N01.01.02.",
		"N01.02.",
		"N01.02.01.01.",
		"N01.02.01.02.",
	}}
	got := buildIndex(orgNodes, dataNodes).data
	want := map[string]loc{
		"N01.":          loc{0, 4, -1, -1},
		"N01.01.":       loc{0, 3, 0, 1},
		"N01.01.01.":    loc{2, 2, 2, 2},
		"N01.01.02.":    loc{3, 3, 3, 3},
		"N01.02.":       loc{4, 4, -1, -1},
		"N01.02.01.01.": loc{-1, -1, -1, -1},
		"N01.02.01.02.": loc{4, 4, 4, 4},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %v, expected %v", got, want)
	}
}

func Test_sortDataByOrgNode(t *testing.T) {
	type test struct {
		input    [][]string
		expected [][]string
	}

	r1 := []string{"N01.", "A"}
	r2 := []string{"N01.100.", "B"}
	r3 := []string{"N01.02.", "C"}
	r4 := []string{"N01.02.222", "D"}
	r5 := []string{"N01.10.", "E"}
	r6 := []string{"N01.10.10.", "F"}
	r7 := []string{"N04.01.", "G"}

	shuffled := [][]string{{"N02."}, {"N03.01."}, {"N01.01.03."}, {"N01."}, {"N03."}, {"N01.02.02."}, {"N01.01."}, {"N01.01.02."}, {"N01.02."}, {"N01.01.01."}, {"N01.02.01."}, {"N03.02."}}
	expected := [][]string{{"N01."}, {"N01.01."}, {"N01.01.01."}, {"N01.01.02."}, {"N01.01.03."}, {"N01.02."}, {"N01.02.01."}, {"N01.02.02."}, {"N02."}, {"N03."}, {"N03.01."}, {"N03.02."}}
	tests := []test{
		test{[][]string{r7, r1}, [][]string{r1, r7}},
		test{[][]string{r3, r1, r2}, [][]string{r1, r3, r2}},
		test{[][]string{r4, r5, r1}, [][]string{r1, r4, r5}},
		test{[][]string{r6, r1}, [][]string{r1, r6}},
		test{shuffled, expected},
	}

	for _, test := range tests {
		got, err := sortDataByOrgNode(test.input, 0)
		want := test.expected

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Expected %v, got %v", want, got)
		}
	}
}
