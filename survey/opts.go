package survey

import (
	"errors"
	"flag"
)

type ProgramOpts struct {
	DataPath          string
	SchemaPath        string
	OrgStructurePath  string
	WorkloadPath      string
	IndexAlgorithm    func(org OrgStructure, dataNodes []string) OrgNodeIndex
	WorkloadAlgorithm func(cuts []Cut, survey *Survey) []CutResult
}

func GetOpts() (ProgramOpts, error) {
	var opts ProgramOpts
	var dataPath string
	var schemaPath string
	var orgStructurePath string
	var workloadPath string
	var concurrentIndex bool
	var concurrentWorkload bool
	var indexImpl func(org OrgStructure, dataNodes []string) OrgNodeIndex
	var workloadImpl func(cuts []Cut, survey *Survey) []CutResult
	flag.StringVar(&dataPath, "data", "", "Path to CSV with survey data.")
	flag.StringVar(&schemaPath, "schema", "", "Path to JSON file with survey schema.")
	flag.StringVar(&orgStructurePath, "org_structure", "", "Path to CSV file with org structure nodes.")
	flag.StringVar(&workloadPath, "workload", "", "Path to JSON file containing workload definition (cuts).")
	flag.BoolVar(&concurrentIndex, "concurrent_index", false, "Use concurrency for building the index.")
	flag.BoolVar(&concurrentWorkload, "concurrent_workload", false, "Use concurrency for processing workload.")
	flag.Parse()

	if dataPath == "" {
		return opts, errors.New("survey data path not provided")
	}

	if schemaPath == "" {
		return opts, errors.New("schema path not provided")
	}

	if orgStructurePath == "" {
		return opts, errors.New("org structure path not provided")
	}

	if workloadPath == "" {
		return opts, errors.New("workload path not provided")
	}

	switch concurrentIndex {
	case true:
		indexImpl = ConcurrentIndex
	case false:
		indexImpl = SequentialIndex
	}

	switch concurrentWorkload {
	case true:
		workloadImpl = ConcurrentCutProcessor
	case false:
		workloadImpl = SequentialCutProcessor
	}

	return ProgramOpts{
		DataPath:          dataPath,
		SchemaPath:        schemaPath,
		OrgStructurePath:  orgStructurePath,
		WorkloadPath:      workloadPath,
		IndexAlgorithm:    indexImpl,
		WorkloadAlgorithm: workloadImpl,
	}, nil
}
