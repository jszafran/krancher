package types

import (
	"encoding/csv"
	"os"
	"sort"
)

type Survey struct {
	schema Schema
}

type loc struct {
	start int
	end   int
}

type orgNodeRollupIndex struct {
	data map[string]loc
}

type orgNodeDirectIndex struct {
	data map[string]loc
}

func buildHeaderColumnMaps(columns []string) (map[string]int, map[int]string) {
	nmToIx := map[string]int{}
	ixToNm := map[int]string{}

	for i, col := range columns {
		nmToIx[col] = i
		ixToNm[i] = col
	}

	return nmToIx, ixToNm
}

func NewSurvey(dataPath string, s Schema, org OrgStructure) Survey {
	file, _ := os.Open(dataPath)
	defer file.Close()

	csvReader := csv.NewReader(file)
	lines, _ := csvReader.ReadAll()

	nmToIx, _ := buildHeaderColumnMaps(lines[0])
	orgIx := nmToIx[s.OrgNodesColumn.Name]
	data := lines[1:]

	sort.SliceStable(data, func(i, j int) bool {
		return data[i][orgIx] < data[j][orgIx]
	})
	return Survey{}
}
