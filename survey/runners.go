package survey

import (
	"fmt"
	"log"
	"time"
)

func RunWorkloadFromOpts(opts WorkloadOpts) error {
	programStart := time.Now()

	dataProvider := CSVDataProvider{DataPath: opts.DataPath}
	schema := SchemaFromJSON(opts.SchemaPath)
	orgStructure := ReadOrgStructureFromCSV(opts.OrgStructurePath, false)
	srv, err := NewSurvey(dataProvider, schema, orgStructure, opts.IndexBuilder)

	if err != nil {
		return fmt.Errorf("failed to create the survey, %s", err)
	}

	cuts, _ := CutsFromJSON(opts.WorkloadPath)
	wrkl := Workload{
		Survey:    &srv,
		Schema:    schema,
		Cuts:      cuts,
		Algorithm: opts.Processor,
	}
	res := wrkl.Run()

	writeStartTime := time.Now()
	err = opts.Persistor(res, opts.OutputPath)

	if err != nil {
		return err
	}

	log.Printf("Results saved successfully (%s)", time.Since(writeStartTime))
	log.Printf("Total program time: %s\n", time.Since(programStart))
	return nil
}
