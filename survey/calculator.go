package survey

import "sort"

type FilterType string

const (
	Direct FilterType = "direct"
	Rollup            = "rollup"
)

type Cut struct {
	Id           string         `json:"id"`
	OrgNode      string         `json:"org_node"`
	Type         FilterType     `json:"type"`
	Demographics map[string]int `json:"demographics"`
}

type QuestionAnswersCounts map[int]int

type CutResult struct {
	Id          string
	Respondents int
	Counts      map[string]QuestionAnswersCounts
}

type Workload struct {
	Cuts []Cut
}

func (w *Workload) GetDemographicsSet() []string {
	set := make(map[string]bool)
	for _, c := range w.Cuts {
		for k, _ := range c.Demographics {
			set[k] = true
		}
	}
	demographics := make([]string, 0)
	for k := range set {
		demographics = append(demographics, k)
	}

	sort.SliceStable(demographics, func(i, j int) bool {
		return demographics[i] < demographics[j]
	})
	return demographics
}

type DataProcessor interface {
	Process(w Workload) []CutResult
}

func newQuestionEmptyResult(c Column) QuestionAnswersCounts {
	counts := QuestionAnswersCounts{}
	for i := c.MinValue; i <= c.MaxValue; i++ {
		counts[i] = 0
	}
	return counts
}

func EmptyCounts(sch Schema) map[string]QuestionAnswersCounts {
	counts := map[string]QuestionAnswersCounts{}
	for _, c := range sch.GetQuestionsColumns() {
		counts[c.Name] = newQuestionEmptyResult(c)
	}
	return counts
}

func NewNoMatchResult(sch Schema, id string) CutResult {
	counts := EmptyCounts(sch)
	return CutResult{
		Respondents: 0,
		Counts:      counts,
		Id:          id,
	}
}

type SynchronousDataProcessor struct {
	Survey *Survey
	Schema Schema
}

func (s *SynchronousDataProcessor) Process(w Workload) []CutResult {
	results := make([]CutResult, 0)
	survey := s.Survey
	schema := s.Survey.schema
	for _, cut := range w.Cuts {
		loc, exists := survey.index[cut.OrgNode]
		if !exists {
			results = append(results, NewNoMatchResult(schema, cut.Id))
		}

		start, end := -1, -1
		switch cut.Type {
		case Direct:
			start = loc.directStart
			end = loc.directEnd
		case Rollup:
			start = loc.rollupStart
			end = loc.rollupEnd
		}

		if start == -1 && end == -1 {
			results = append(results, NewNoMatchResult(schema, cut.Id))
		}

		counts := EmptyCounts(schema)
		questions := schema.GetQuestionsCodes()
		respondents := 0
		for i := start; i <= end; i++ {
			eligible := true
			for k, v := range cut.Demographics {
				if survey.demographicData[k][i] != v {
					eligible = false
				}
			}
			if eligible {
				respondents++
				for _, qst := range questions {
					v := survey.answersData[qst][i]
					if v > -1 {
						counts[qst][v]++
					}
				}
			}
		}
		results = append(results, CutResult{
			Id:          cut.Id,
			Respondents: respondents,
			Counts:      counts,
		})

	}
	return results
}
