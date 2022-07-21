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

type QuestionAnswersCounts map[uint]uint

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

	var start, end int = -1, -1

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

	return CutResult{}
}
