package types

import (
	"encoding/csv"
	"errors"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Survey struct {
	schema Schema
	index  orgNodeIndex
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
			}

		}
		ixData[node] = l
	}
	return orgNodeIndex{ixData}
}

func sortDataByOrgNode(data [][]string, orgColIx int) ([][]string, error) {
	parseNode := func(node string) ([]int, error) {
		node = strings.Replace(node, "N", "", -1)
		spl := strings.Split(node, ".")
		if spl[len(spl)-1] == "" {
			spl = spl[:len(spl)-1]
		}
		nodeInts := make([]int, len(spl))
		for i, el := range spl {
			nodeInt, err := strconv.Atoi(el)
			if err != nil {
				return nil, errors.New("error when parsing node " + node + " to integers")
			}
			nodeInts[i] = nodeInt
		}
		return nodeInts, nil
	}

	sort.SliceStable(data, func(i, j int) bool {
		n1, err1 := parseNode(data[i][orgColIx])
		n2, err2 := parseNode(data[j][orgColIx])
		if err1 != nil || err2 != nil {
			if err1 != nil {
				log.Fatal(err1)
			}
			log.Fatal(err2)
		}
		minLen := int(math.Min(float64(len(n1)), float64(len(n2))))
		for i := 0; i < minLen; i++ {
			if n1[i] == n2[i] {
				continue
			}
			return n1[i] < n2[i]
		}
		return false
	})
	return data, nil
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
		str1 := strings.Replace(data[i][orgIx], "N", "", 0)
		str2 := strings.Replace(data[j][orgIx], "N", "", 0)
		spl1 := strings.Split(str1, ".")
		spl2 := strings.Split(str2, ",")
		mlen := int(math.Min(float64(len(spl1)), float64(len(spl2))))
		for i := 0; i < mlen; i++ {
			n1, err1 := strconv.Atoi(spl1[i])
			n2, err2 := strconv.Atoi(spl2[i])
			if err1 != nil || err2 != nil {
				log.Fatalf("Failed to convert values %v or %v to int.", spl1[i], spl2[i])
			}
			if n1 == n2 {
				continue
			}
			return n1 < n2
		}
		return false
	})
	dataNodes := make([]string, 0)

	for _, row := range data {
		dataNodes = append(dataNodes, row[nmToIx["org"]])
	}

	_ = buildIndex(org, dataNodes)

	return Survey{}
}
