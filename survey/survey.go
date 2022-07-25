package survey

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Survey struct {
	schema          Schema
	index           OrgNodeIndex
	answersData     map[string][]int
	demographicData map[string][]int
}

type loc struct {
	rollupStart int
	rollupEnd   int
	directStart int
	directEnd   int
}

type IndexBuilder string

const (
	Concurrent IndexBuilder = "concurrent"
	Sequential              = "sequential"
)

func newLoc() loc {
	return loc{-1, -1, -1, -1}
}

type OrgNodeIndex map[string]loc

type DataProvider interface {
	GetData() ([][]string, error)
}

type CSVDataProvider struct {
	DataPath string
}

func (c CSVDataProvider) GetData() ([][]string, error) {
	file, err := os.Open(c.DataPath)
	defer file.Close()

	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	lines, err1 := reader.ReadAll()

	if err1 != nil {
		return nil, err1
	}

	return lines, nil
}

type InMemoryDataProvider struct {
	Data [][]string
}

func (i InMemoryDataProvider) GetData() ([][]string, error) {
	return i.Data, nil
}

type IndexElement struct {
	OrgNode  string
	Location loc
}

func ConcurrentIndex(org OrgStructure, dataNodes []string) OrgNodeIndex {
	var wg sync.WaitGroup
	ch := make(chan IndexElement, len(org.nodes))

	calculateIndexElement := func(node string, dataNodes []string, ch chan<- IndexElement, wg *sync.WaitGroup) {
		defer wg.Done()
		ch <- calculateLocation(node, dataNodes)
	}

	for _, node := range org.nodes {
		wg.Add(1)
		go calculateIndexElement(node, dataNodes, ch, &wg)
	}

	wg.Wait()
	close(ch)

	res := make(map[string]loc)

	for element := range ch {
		res[element.OrgNode] = element.Location
	}

	return res
}

func SequentialIndex(org OrgStructure, dataNodes []string) OrgNodeIndex {
	res := make(map[string]loc)

	for _, node := range org.nodes {
		res[node] = calculateLocation(node, dataNodes).Location
	}

	return res
}

func calculateLocation(node string, dataNodes []string) IndexElement {
	l := newLoc()
	for i, dataNode := range dataNodes {
		if len(node) > len(dataNode) {
			continue
		}

		if node == dataNode[:len(node)] {
			if l.rollupStart == -1 {
				l.rollupStart = i
				l.rollupEnd = i
			} else {
				l.rollupEnd++
			}
		}

		if node == dataNode {
			if l.directStart == -1 {
				l.directStart = i
				l.directEnd = i
			} else {
				l.directEnd++
			}
		}
	}
	return IndexElement{OrgNode: node, Location: l}
}
func buildHeaderColumnMaps(columns []string) map[string]int {
	nmToIx := map[string]int{}

	for i, col := range columns {
		nmToIx[col] = i
	}
	return nmToIx
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

	slicesOfIntEqual := func(s1 []int, s2 []int) bool {
		if len(s1) != len(s2) {
			return false
		}
		for i, v1 := range s1 {
			if v1 != s2[i] {
				return false
			}
		}
		return true
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
		if slicesOfIntEqual(n1[:minLen], n2[:minLen]) {
			return len(n1) < len(n2)
		}
		for x := 0; x < minLen; x++ {
			if n1[x] == n2[x] {
			} else {
				return n1[x] < n2[x]
			}
		}
		return true
	})
	return data, nil
}

func parseColumnsData(data [][]string, codes []string, nmToIxMap map[string]int) (map[string][]int, error) {
	parsedData := make(map[string][]int, 0)
	for _, row := range data {
		for _, code := range codes {
			ix, ok := nmToIxMap[code]
			if !ok {
				return nil, errors.New(fmt.Sprintf("failed to find %s in the name to column index map", code))
			}

			rawVal := row[ix]

			if rawVal == "" {
				parsedData[code] = append(parsedData[code], -1)
				continue
			}
			parsedFloat, err := strconv.ParseFloat(rawVal, 64)
			parsedInt := int(parsedFloat)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed when converting %s to float", rawVal))
			}
			parsedData[code] = append(parsedData[code], parsedInt)
		}
	}
	return parsedData, nil
}

func NewSurvey(
	dataProvider DataProvider,
	s Schema,
	org OrgStructure,
	ib IndexBuilder,
) (Survey, error) {
	lines, err := dataProvider.GetData()
	if err != nil {
		return Survey{}, err
	}

	nmToIx := buildHeaderColumnMaps(lines[0])
	orgIx := nmToIx[s.OrgNodesColumn.Name]

	sortDataStart := time.Now()
	data, _ := sortDataByOrgNode(lines[1:], orgIx)
	log.Printf("Sorting data took %s\n", time.Since(sortDataStart))

	dataNodes := make([]string, 0)
	for _, row := range data {
		dataNodes = append(dataNodes, row[orgIx])
	}

	ixBuildStart := time.Now()
	var ixImpl func(org OrgStructure, dataNodes []string) OrgNodeIndex
	switch ib {
	case Concurrent:
		ixImpl = ConcurrentIndex
	case Sequential:
		ixImpl = SequentialIndex
	}
	index := ixImpl(org, dataNodes)
	log.Printf("Building index took %s\n", time.Since(ixBuildStart))

	demogsStart := time.Now()
	demogs, err1 := parseColumnsData(data, s.GetDemographicsCodes(), nmToIx)
	if err1 != nil {
		return Survey{}, err1
	}
	demogsTotalTime := time.Since(demogsStart)

	qstStart := time.Now()
	answers, err2 := parseColumnsData(data, s.GetQuestionsCodes(), nmToIx)
	if err2 != nil {
		return Survey{}, err2
	}
	qstTotalTime := time.Since(qstStart)
	totalTime := time.Since(demogsStart)

	log.Printf("Reading demogs took %s\n", demogsTotalTime)
	log.Printf("Reading questions took %s\n", qstTotalTime)
	log.Printf("Total time for parsing questions + demogs: %s\n", totalTime)

	return Survey{
		schema:          s,
		index:           index,
		answersData:     answers,
		demographicData: demogs,
	}, nil
}
