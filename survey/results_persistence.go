package survey

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
)

// StandardJSONPersistor full JSON - no optimization
func StandardJSONPersistor(res []CutResult, outputPath string) error {
	file, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputPath, file, fs.ModePerm)
	return err
}

// SkipEmptyCutsCountsJSONPersistor counts for cuts with 0 respondents are not included in the JSON
func SkipEmptyCutsCountsJSONPersistor(res []CutResult, outputPath string) error {
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

	return StandardJSONPersistor(resOpt, outputPath)
}
