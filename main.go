package main

import (
	"krancher/survey"
	"log"
	"time"
)

func main() {
	programStart := time.Now()
	dataProvider := survey.CSVDataProvider{DataPath: "resources/itest_data_50x.csv"}
	schema := survey.SchemaFromJSON("resources/itest_schema.json")
	orgStructure := survey.ReadOrgStructureFromCSV("resources/itest_org.csv", false)
	srv, err := survey.NewSurvey(dataProvider, schema, orgStructure, survey.Concurrent)

	if err != nil {
		log.Fatalf("failed to create the survey, %s", err)
	}

	cuts, _ := survey.CutsFromJSON("resources/itest_all_cuts.json")
	wrkl := survey.Workload{
		Survey:    &srv,
		Schema:    survey.Schema{},
		Cuts:      cuts,
		Algorithm: survey.ConcurrentCutProcessor,
	}
	wrkl.Run()
	log.Printf("Total program time: %s\n", time.Since(programStart))
}
