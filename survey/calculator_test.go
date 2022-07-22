package survey

import (
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
