package types

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
