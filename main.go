package main

import (
	"fmt"
	"krancher/survey"
	"log"
	"time"
)

func main() {
	programStart := time.Now()
	dataProvider := survey.CSVDataProvider{DataPath: "resources/itest_data_2x.csv"}
	schema := survey.SchemaFromJSON("resources/itest_schema.json")
	orgStructure := survey.ReadOrgStructureFromCSV("resources/itest_org.csv", false)
	srv, err := survey.NewSurvey(dataProvider, schema, orgStructure)

	if err != nil {
		log.Fatalf("failed to create the survey, %s", err)
	}
	//c1 := survey.Cut{
	//	Id:           "#1",
	//	OrgNode:      "N00.01.02.",
	//	Type:         survey.Direct,
	//	Demographics: map[string]int{},
	//}

	calcTime := time.Now()
	dataProc := survey.SynchronousDataProcessor{Survey: &srv, Schema: schema}
	wrkl, _ := survey.WorkloadFromJSON("resources/itest_all_cuts.json")
	res := dataProc.Process(wrkl)
	log.Printf("Total time for calculating cuts: %s", time.Since(calcTime))
	log.Printf("Total program time: %s", time.Since(programStart))
	fmt.Println(len(res))
}
