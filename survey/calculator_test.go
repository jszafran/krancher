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

	got := NewNoMatchResult(sch, "#1")
	want := CutResult{"#1", 0, map[string]QuestionAnswersCounts{
		"Q1": {1: 0, 2: 0, 3: 0},
		"Q2": {1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0},
	}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}

func TestSynchronousDataProcessor_Process(t *testing.T) {
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
	srv, _ := NewSurvey(dp, sch, org, Sequential)
	emptyDemog := make(map[string]int)
	cuts := []Cut{
		{Id: "C1", Type: Rollup, OrgNode: "N01.01.", Demographics: emptyDemog},
		{Id: "C2", Type: Direct, OrgNode: "N02.", Demographics: emptyDemog},
		{Id: "C3", Type: Direct, OrgNode: "N02.", Demographics: map[string]int{"D1": 1, "D2": 1}},
		{Id: "C4", Type: Direct, OrgNode: "N02.", Demographics: map[string]int{"D1": 1, "D2": 2}},
		{Id: "C5", Type: Direct, OrgNode: "N02.", Demographics: map[string]int{"D2": 2}},
		{Id: "C6", Type: Rollup, OrgNode: "N01.", Demographics: map[string]int{"D1": 1}},
	}

	syncDataProc := SynchronousDataProcessor{
		Survey: &srv,
		Schema: sch,
	}

	got := syncDataProc.Process(Workload{cuts})
	want := []CutResult{
		{
			Id:          "C1",
			Respondents: 3,
			Counts: map[string]QuestionAnswersCounts{
				"Q1": {1: 1, 2: 1, 3: 0, 4: 0, 5: 1, 6: 0},
				"Q2": {1: 0, 2: 2, 3: 0, 4: 0, 5: 1, 6: 0},
				"Q3": {1: 0, 2: 0, 3: 0, 4: 0, 5: 3, 6: 0},
			},
		},
		{
			Id:          "C2",
			Respondents: 1,
			Counts: map[string]QuestionAnswersCounts{
				"Q1": {1: 0, 2: 0, 3: 1, 4: 0, 5: 0, 6: 0},
				"Q2": {1: 0, 2: 1, 3: 0, 4: 0, 5: 0, 6: 0},
				"Q3": {1: 1, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0},
			},
		},
		{
			Id:          "C3",
			Respondents: 0,
			Counts: map[string]QuestionAnswersCounts{
				"Q1": {1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0},
				"Q2": {1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0},
				"Q3": {1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0},
			},
		},
		{
			Id:          "C4",
			Respondents: 1,
			Counts: map[string]QuestionAnswersCounts{
				"Q1": {1: 0, 2: 0, 3: 1, 4: 0, 5: 0, 6: 0},
				"Q2": {1: 0, 2: 1, 3: 0, 4: 0, 5: 0, 6: 0},
				"Q3": {1: 1, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0},
			},
		},
		{
			Id:          "C5",
			Respondents: 1,
			Counts: map[string]QuestionAnswersCounts{
				"Q1": {1: 0, 2: 0, 3: 1, 4: 0, 5: 0, 6: 0},
				"Q2": {1: 0, 2: 1, 3: 0, 4: 0, 5: 0, 6: 0},
				"Q3": {1: 1, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0},
			},
		},
		{
			Id:          "C6",
			Respondents: 2,
			Counts: map[string]QuestionAnswersCounts{
				"Q1": {1: 2, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0},
				"Q2": {1: 0, 2: 1, 3: 0, 4: 0, 5: 0, 6: 0},
				"Q3": {1: 0, 2: 0, 3: 1, 4: 0, 5: 1, 6: 0},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}

func TestWorkload_GetDemographicsSet(t *testing.T) {
	cuts := []Cut{
		{
			Demographics: nil,
		},
		{
			Demographics: map[string]int{"D5": 1, "D2": 2},
		},
		{
			Demographics: map[string]int{"D3": 1, "D4": 2},
		},
		{
			Demographics: map[string]int{"D1": 1},
		},
	}
	wrkl := Workload{cuts}

	got := wrkl.GetDemographicsSet()
	want := []string{"D1", "D2", "D3", "D4", "D5"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %v but expected %v", got, want)
	}
}
