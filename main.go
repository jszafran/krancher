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
	_, err := survey.NewSurvey(dataProvider, schema, orgStructure)

	if err != nil {
		log.Fatalf("failed to create the survey, %s", err)
	}

	log.Printf("Total program time: %s\n", time.Since(programStart))
}
