package survey

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
)

type ResultsPersistenceFunc func(res []CutResult, outputPath string) error

// StandardJSONPersistence full JSON - no optimization
func StandardJSONPersistence(res []CutResult, outputPath string) error {
	file, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputPath, file, fs.ModePerm)
	return err
}

// EmptyCutOptimizedJSONPersistence counts for cuts with 0 respondents are not included in the JSON
func EmptyCutOptimizedJSONPersistence(res []CutResult, outputPath string) error {
	resOpt := make([]CutResult, 0)
	for _, r := range res {
		if r.Respondents == 0 {
			resOpt = append(resOpt, CutResult{
				Id:          r.Id,
				Respondents: 0,
				Counts:      nil,
			})
			continue
		}
		resOpt = append(resOpt, r)
	}

	return StandardJSONPersistence(resOpt, outputPath)
}
