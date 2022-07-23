package survey

const (
	Direct FilterType = "direct"
	Rollup            = "rollup"
)

type Cut struct {
	Id           string
	OrgNode      string
	Type         FilterType
	Demographics map[string]int
}

type QuestionAnswersCounts map[int]int

type CutResult struct {
	Id          string
	Respondents int
	Counts      map[string]QuestionAnswersCounts
}

type DataProcessor interface {
	Process(cuts []Cut) []CutResult
}

func newQuestionEmptyResult(c Column) QuestionAnswersCounts {
	cnts := QuestionAnswersCounts{}
	for i := c.MinValue; i <= c.MaxValue; i++ {
		cnts[i] = 0
	}
	return cnts
}

func EmptyCounts(sch Schema) map[string]QuestionAnswersCounts {
	counts := map[string]QuestionAnswersCounts{}
	for _, c := range sch.GetQuestionsColumns() {
		counts[c.Name] = newQuestionEmptyResult(c)
	}
	return counts
}

func NewNoMatchResult(sch Schema, id string) CutResult {
	counts := map[string]QuestionAnswersCounts{}
	for _, c := range sch.GetQuestionsColumns() {
		counts[c.Name] = newQuestionEmptyResult(c)
	}
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

func (s *SynchronousDataProcessor) Process(cuts []Cut) []CutResult {
	results := make([]CutResult, 0)
	survey := s.Survey
	schema := s.Survey.schema
	for _, cut := range cuts {
		loc, exists := survey.index[cut.OrgNode]
		if !exists {
			results = append(results, NewNoMatchResult(schema, cut.Id))
		}

		start, end := -1, -1
		switch cut.Type {
		case Direct:
			{
				start = loc.directStart
				end = loc.directEnd
			}
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
