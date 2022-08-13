package main

import (
	"krancher/survey"
	"log"
	"time"
)

func main() {
	programStart := time.Now()

	opts, optsErr := survey.ParseOpts()
	if optsErr != nil {
		log.Fatal(optsErr)
	}

	dataProvider := survey.CSVDataProvider{DataPath: opts.DataPath}
	schema := survey.SchemaFromJSON(opts.SchemaPath)
	orgStructure := survey.ReadOrgStructureFromCSV(opts.OrgStructurePath, false)
	srv, err := survey.NewSurvey(dataProvider, schema, orgStructure, opts.IndexAlgorithm)

	if err != nil {
		log.Fatalf("failed to create the survey, %s", err)
	}

	cuts, _ := survey.CutsFromJSON(opts.WorkloadPath)
	wrkl := survey.Workload{
		Survey:    &srv,
		Schema:    survey.Schema{},
		Cuts:      cuts,
		Algorithm: opts.WorkloadAlgorithm,
	}
	res := wrkl.Run()

	writeStartTime := time.Now()
	err = opts.ResultsPersistence(res, opts.OutputPath)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Results saved successfully (%s)", time.Since(writeStartTime))
	log.Printf("Total program time: %s\n", time.Since(programStart))
}
