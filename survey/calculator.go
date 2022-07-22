package survey

import (
	"log"
)

const (
	Direct FilterType = "direct"
	Rollup            = "rollup"
)

type Cut struct {
	OrgNode      string
	Type         FilterType
	Demographics map[string]int
}

type QuestionAnswersCounts map[int]int

type CutResult struct {
	Respondents int
	Counts      map[string]QuestionAnswersCounts
}

func newQuestionEmptyResult(c Column) QuestionAnswersCounts {
	cnts := QuestionAnswersCounts{}
	for i := c.MinValue; i <= c.MaxValue; i++ {
		cnts[i] = 0
	}
	return cnts
}

func EmptyCounts(sch Schema) map[string]QuestionAnswersCounts {
	res := make(map[string]QuestionAnswersCounts)
	for _, qst := range sch.GetDemographicsCodes() {
		res[qst] = make(QuestionAnswersCounts)
	}
	return res
}

func NewNoMatchResult(sch Schema) CutResult {
	counts := map[string]QuestionAnswersCounts{}
	for _, c := range sch.GetQuestionsColumns() {
		counts[c.Name] = newQuestionEmptyResult(c)
	}
	return CutResult{
		Respondents: 0,
		Counts:      counts,
	}
}

func CalculateCounts(srv *Survey, sch Schema, c Cut) CutResult {
	loc, exists := srv.index[c.OrgNode]
	if !exists {
		log.Fatal("corrupted index")
	}

	var start, end = -1, -1

	if c.Type == Direct {
		start = loc.directStart
		end = loc.directEnd
	} else if c.Type == Rollup {
		start = loc.rollupStart
		end = loc.rollupEnd
	}

	if start == -1 && end == -1 {
		return NewNoMatchResult(sch)
	}

	counts := EmptyCounts(srv.schema)
	questions := sch.GetQuestionsCodes()
	if len(c.Demographics) == 0 {
		for i := start; i <= end; i++ {
			for _, qst := range questions {
				v := srv.answersData[qst][i]
				if v > -1 {
					counts[qst][v]++
				}
			}
		}
		return CutResult{Respondents: end - start + 1, Counts: counts}
	}

	respondents := 0
	for i := start; i <= end; i++ {
		// check demographic data
		eligible := true
		for k, v := range c.Demographics {
			if srv.demographicData[k][i] != v {
				eligible = false
			}
		}
		if eligible {
			respondents++
			for _, qst := range questions {
				v := srv.answersData[qst][i]
				if v > 1 {
					counts[qst][v]++
				}
			}
		}
	}
	return CutResult{Respondents: respondents, Counts: counts}
}
