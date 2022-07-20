package survey

import "log"

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

func NewNoMatchResult(sch Schema) CutResult {
	counts := map[string]QuestionAnswersCounts{}
	for _, qst := range sch.GetQuestionsColumns() {
		qstCounts := QuestionAnswersCounts{}
		for i := qst.MinValue; i <= qst.MaxValue; i++ {
			qstCounts[i] = 0
		}
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

	var start, end int

	if c.Type == Direct {
		start = loc.directStart
		end = loc.directEnd
	} else {
		start = loc.rollupStart
		end = loc.rollupEnd
	}


	for _, resp := range
	return CutResult{}
}
