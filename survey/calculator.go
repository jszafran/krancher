package survey

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type FilterType string

const (
	Direct FilterType = "direct"
	Rollup            = "rollup"
)

type Cut struct {
	Id           string         `json:"id"`
	OrgNode      string         `json:"org_node"`
	Type         FilterType     `json:"type"`
	Demographics map[string]int `json:"demographics"`
}

type QuestionAnswersCounts map[int]int

type CutResult struct {
	Id          string
	Respondents int
	Counts      map[string]QuestionAnswersCounts
}

type Workload struct {
	Survey    *Survey
	Schema    Schema
	Cuts      []Cut
	Algorithm func(cuts []Cut, survey *Survey) []CutResult
}

func (w *Workload) GetDemographicsSet() []string {
	set := make(map[string]bool)
	for _, c := range w.Cuts {
		for k := range c.Demographics {
			set[k] = true
		}
	}
	demographics := make([]string, 0)
	for k := range set {
		demographics = append(demographics, k)
	}

	sort.SliceStable(demographics, func(i, j int) bool {
		return demographics[i] < demographics[j]
	})
	return demographics
}

func (w *Workload) Run() []CutResult {
	t1 := time.Now()
	results := w.Algorithm(w.Cuts, w.Survey)
	log.Printf("Workload (%d cuts) completed! It took %s.", len(w.Cuts), time.Since(t1))
	return results
}

type DataProcessor interface {
	Process(w Workload) []CutResult
}

func newQuestionEmptyResult(c Column) QuestionAnswersCounts {
	counts := QuestionAnswersCounts{}
	for i := c.MinValue; i <= c.MaxValue; i++ {
		counts[i] = 0
	}
	return counts
}

func EmptyCounts(sch Schema) map[string]QuestionAnswersCounts {
	counts := map[string]QuestionAnswersCounts{}
	for _, c := range sch.GetQuestionsColumns() {
		counts[c.Name] = newQuestionEmptyResult(c)
	}
	return counts
}

func NewNoMatchResult(sch Schema, id string) CutResult {
	counts := EmptyCounts(sch)
	return CutResult{
		Respondents: 0,
		Counts:      counts,
		Id:          id,
	}
}

func calculateCut(cut Cut, survey *Survey) CutResult {
	schema := survey.schema
	loc, exists := survey.index[cut.OrgNode]
	if !exists {
		return NewNoMatchResult(schema, cut.Id)
	}

	start, end := -1, -1

	switch cut.Type {
	case Direct:
		start = loc.directStart
		end = loc.directEnd
	case Rollup:
		start = loc.rollupStart
		end = loc.rollupEnd
	}

	if start == -1 && end == -1 {
		return NewNoMatchResult(schema, cut.Id)
	}

	counts := EmptyCounts(schema)
	questions := schema.GetQuestionsCodes()
	respondents := 0
	for i := start; i <= end; i++ {
		eligible := true
		for k, v := range cut.Demographics {
			if survey.demographicData[k][i] != v {
				eligible = false
			}
		}
		if eligible {
			respondents++
			for _, qst := range questions {
				v := survey.answersData[qst][i]
				if v > -1 {
					counts[qst][v]++
				}
			}
		}
	}
	return CutResult{
		Id:          cut.Id,
		Respondents: respondents,
		Counts:      counts,
	}
}

func SequentialCutProcessor(cuts []Cut, survey *Survey) []CutResult {
	res := make([]CutResult, 0)
	for _, cut := range cuts {
		cr := calculateCut(cut, survey)
		res = append(res, cr)
	}
	return res
}

func ConcurrentCutProcessor(cuts []Cut, survey *Survey) []CutResult {
	var wg sync.WaitGroup
	ch := make(chan CutResult, len(cuts))

	processCut := func(c Cut, s *Survey, ch chan<- CutResult, wg *sync.WaitGroup) {
		defer wg.Done()
		ch <- calculateCut(c, s)
	}

	for _, cut := range cuts {
		wg.Add(1)
		go processCut(cut, survey, ch, &wg)
	}

	wg.Wait()
	close(ch)

	res := make([]CutResult, 0)

	for cr := range ch {
		res = append(res, cr)
	}

	return res
}

func CutsFromJSON(path string) ([]Cut, error) {
	jsonFile, err1 := os.Open(path)
	defer jsonFile.Close()
	if err1 != nil {
		log.Fatal(err1)
	}

	bytes, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		log.Fatal(err2)
	}

	var cuts []Cut
	err3 := json.Unmarshal(bytes, &cuts)
	if err3 != nil {
		log.Fatal(err3)
	}

	return cuts, nil
}

func PersistResults(res []CutResult, path string) error {
	file, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, file, fs.ModePerm)
	return err
}
