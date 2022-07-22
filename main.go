package main

import (
	"krancher/survey"
	"log"
	"time"
)

func main() {
	programStart := time.Now()
	dataProvider := survey.CSVDataProvider{DataPath: "resources/fake_big_survey.csv"}
	schema := survey.SchemaFromJSON("resources/example_schema.json")
	orgStructure := survey.ReadOrgStructureFromCSV("resources/fake_org.csv", false)
	srv, err := survey.NewSurvey(dataProvider, schema, orgStructure)

	if err != nil {
		log.Fatalf("failed to create the survey, %s", err)
	}
	c := survey.Cut{
		OrgNode:      "N01.",
		Type:         survey.Rollup,
		Demographics: nil,
	}
	_ = survey.CalculateCounts(&srv, schema, c)
	log.Printf("Total program time: %s", time.Since(programStart))
}
