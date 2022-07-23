package survey

import (
	"log"
	"reflect"
	"testing"
)

func TestNewNoMatchResult(t *testing.T) {
	sch := Schema{Columns: []Column{
		{"Q1", "Q1", 1, 3, true, Question},
		{"Q2", "Q2", 1, 7, true, Question},
		{"D1", "D1", 1, 2, true, Demography},
	}}

	got := NewNoMatchResult(sch)
	want := CutResult{0, map[string]QuestionAnswersCounts{
		"Q1": {1: 0, 2: 0, 3: 0},
		"Q2": {1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0},
	}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}

func TestCalculateCounts(t *testing.T) {
	org := OrgStructure{nodes: []string{"N01.", "N01.01.", "N01.01.02.", "N01.02.", "N02."}}
	dp := InMemoryDataProvider{Data: [][]string{
		{"org", "Q1", "Q2", "Q3", "D1", "D2"},
		{"N01.", "1", "", "3", "1", ""},
		{"N01.01.", "1", "2", "5", "1", ""},
		{"N01.01.", "2", "2", "5", "2", "1"},
		{"N01.01.02.", "5", "5", "5", "", "1"},
		{"N02.", "3", "2", "1", "1", "2"},
	}}
	sch := Schema{Columns: []Column{
		{"Q1", "Q1", 1, 6, true, Question},
		{"Q2", "Q2", 1, 6, true, Question},
		{"Q3", "Q3", 1, 6, true, Question},
		{"D1", "D1", 1, 6, true, Demography},
		{"D2", "D2", 1, 6, true, Demography},
	}}
	srv, _ := NewSurvey(dp, sch, org)
	c1 := Cut{Type: Rollup, OrgNode: "N01.01.", Demographics: make(map[string]int, 0)}

	got := CalculateCounts(&srv, sch, c1)
	want := CutResult{
		Respondents: 3,
		Counts: map[string]QuestionAnswersCounts{
			"Q1": {1: 1, 2: 1, 3: 0, 4: 0, 5: 1, 6: 0},
			"Q2": {1: 0, 2: 2, 3: 0, 4: 0, 5: 1, 6: 0},
			"Q3": {1: 0, 2: 0, 3: 0, 4: 0, 5: 3, 6: 0},
		},
	}

	if !reflect.DeepEqual(got, want) {
		log.Fatalf("Expected %v but got %v", want, got)
	}
}
