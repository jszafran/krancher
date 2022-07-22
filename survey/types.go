package survey

import (
	"errors"
	"fmt"
)

type ColumnType string
type FilterType string

const (
	Question   ColumnType = "question"
	Demography            = "demography"
)

type Column struct {
	Name     string     `json:"name,omitempty"`
	Text     string     `json:"text,omitempty"`
	MinValue int        `json:"min_value,omitempty"`
	MaxValue int        `json:"max_value,omitempty"`
	Nullable bool       `json:"nullable,omitempty"`
	OfType   ColumnType `json:"of_type,omitempty"`
}

type OrgNodesColumn struct {
	Name string `json:"name,omitempty"`
}

type Schema struct {
	Columns        []Column       `json:"columns,omitempty"`
	OrgNodesColumn OrgNodesColumn `json:"org_nodes_column"`
}

func filterColumnsByType(s Schema, ct ColumnType) []Column {
	cols := make([]Column, 0)
	for _, i := range s.Columns {
		if i.OfType == ct {
			cols = append(cols, i)
		}
	}
	return cols
}

func (s Schema) GetQuestionsColumns() []Column {
	return filterColumnsByType(s, Question)
}

func (s Schema) GetDemographicsColumns() []Column {
	return filterColumnsByType(s, Demography)
}

func (s Schema) getNames(ct ColumnType) []string {
	codes := make([]string, 0)
	for _, c := range filterColumnsByType(s, ct) {
		codes = append(codes, c.Name)
	}
	return codes
}

func (s Schema) GetQuestionsCodes() []string {
	return s.getNames(Question)
}

func (s Schema) GetDemographicsCodes() []string {
	return s.getNames(Demography)
}

func (s Schema) GetColumnForName(name string) (Column, error) {
	var column Column
	for _, c := range s.Columns {
		if c.Name == name {
			column = c
		}
	}
	return column, errors.New(fmt.Sprintf("Column with %s name not found", name))
}
