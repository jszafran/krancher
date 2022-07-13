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
	rollupStart int
	rollupEnd   int
	directStart int
	directEnd   int
}

func newLoc() loc {
	return loc{-1, -1, -1, -1}
}

type orgNodeIndex struct {
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

func buildIndex(org OrgStructure, dataNodes []string) orgNodeIndex {
	ixData := map[string]loc{}
	for _, node := range org.nodes {
		l := newLoc()
		for i, dataNode := range dataNodes {
			if len(node) > len(dataNode) {
				continue
			}
			// rollup match
			if node == dataNode[:len(node)] {
				if l.rollupStart == -1 {
					l.rollupStart = i
					l.rollupEnd = i
				} else {
					l.rollupEnd++
				}
			}

			// direct match
			if node == dataNode {
				if l.directStart == -1 {
					l.directStart = i
					l.directEnd = i
				} else {
					l.directEnd++
				}
				continue
			}

		}
		ixData[node] = l
	}
	return orgNodeIndex{ixData}
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
