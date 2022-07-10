package types

type ColumnType string
type FilterType string

const (
	Question   ColumnType = "question"
	Demography            = "demography"
)

const (
	Direct FilterType = "direct"
	Rollup            = "rollup"
)

type Column struct {
	code     string
	minValue uint
	maxValue uint
	nullable bool
	ofType   ColumnType
}

type OrgNodesColumn struct {
	code string
}

type Schema struct {
	columns        []Column
	OrgNodesColumn OrgNodesColumn
}

func filterColumnsByType(s *Schema, ct ColumnType) []Column {
	cols := make([]Column, 0)
	for _, i := range s.columns {
		if i.ofType == ct {
			cols = append(cols, i)
		}
	}
	return cols
}

func (s *Schema) GetQuestionsColumns() []Column {
	return filterColumnsByType(s, Question)
}

func (s *Schema) GetDemographicsColumns() []Column {
	return filterColumnsByType(s, Demography)
}

func (s *Schema) getCodes(ct ColumnType) []string {
	codes := make([]string, 0)
	for _, c := range filterColumnsByType(s, Question) {
		codes = append(codes, c.code)
	}
	return codes
}

func (s *Schema) GetQuestionsCodes() []string {
	return s.getCodes(Question)
}

func (s *Schema) GetDemographicsCodes() []string {
	return s.getCodes(Demography)
}
