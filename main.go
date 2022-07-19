package main

import (
	"fmt"
	"krancher/kio"
	"krancher/types"
	"log"
	"strconv"
)

func main() {
	s := "6"
	sInt, err := strconv.Atoi(s)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sInt)
	schema := kio.SchemaFromJSON("resources/example_schema.json")
	orgStructure := types.ReadOrgStructureFromCSV("resources/fake_org.csv", false)
	survey, err := types.NewSurvey("resources/fake_data.csv", schema, orgStructure)

	for _, x := range schema.Columns {
		fmt.Println(x)
	}
	if err != nil {
		log.Fatalf("failed to create the survey, %s", err)
	}

	fmt.Println(survey)
}
