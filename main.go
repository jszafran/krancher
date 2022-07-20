package main

import (
	"krancher/kio"
	"krancher/types"
	"log"
	"time"
)

func main() {
	programStart := time.Now()

	schema := kio.SchemaFromJSON("resources/example_schema.json")
	orgStructure := types.ReadOrgStructureFromCSV("resources/fake_org.csv", false)
	_, err := types.NewSurvey("resources/fake_big_survey.csv", schema, orgStructure)

	if err != nil {
		log.Fatalf("failed to create the survey, %s", err)
	}

	log.Printf("Total program time: %s\n", time.Since(programStart))
}
