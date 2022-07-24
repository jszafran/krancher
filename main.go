package main

import (
	"krancher/survey"
	"log"
	"time"
)

func main() {
	programStart := time.Now()
	dataProvider := survey.CSVDataProvider{DataPath: "resources/fake_big_survey.csv"}
	schema := survey.SchemaFromJSON("resources/big_schema.json")
	orgStructure := survey.ReadOrgStructureFromCSV("resources/fake_org.csv", false)
	srv, err := survey.NewSurvey(dataProvider, schema, orgStructure)

	if err != nil {
		log.Fatalf("failed to create the survey, %s", err)
	}
	c1 := survey.Cut{
		Id:           "#1",
		OrgNode:      "N01.",
		Type:         survey.Rollup,
		Demographics: map[string]int{},
	}
	//c2 := survey.Cut{
	//	Id:           "#2",
	//	OrgNode:      "N01.01.",
	//	Type:         survey.Rollup,
	//	Demographics: map[string]int{},
	//}
	//c3 := survey.Cut{
	//	Id:           "#3",
	//	OrgNode:      "N01.",
	//	Type:         survey.Direct,
	//	Demographics: map[string]int{},
	//}

	calcTime := time.Now()
	dataProc := survey.SynchronousDataProcessor{Survey: &srv, Schema: schema}
	wrkl := survey.Workload{Cuts: []survey.Cut{c1}}
	_ = dataProc.Process(wrkl)
	log.Printf("Total time for calculating cuts: %s", time.Since(calcTime))
	log.Printf("Total program time: %s", time.Since(programStart))
}
