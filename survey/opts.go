package survey

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"
)

type ProgramOpts struct {
	DataPath           string
	SchemaPath         string
	OrgStructurePath   string
	WorkloadPath       string
	OutputPath         string
	IndexAlgorithm     func(org OrgStructure, dataNodes []string) OrgNodeIndex
	WorkloadAlgorithm  func(cuts []Cut, survey *Survey) []CutResult
	ResultsPersistence ResultsPersistenceFunc
}

func getCurrentTimeResultsFilename() string {
	return strings.Replace(fmt.Sprintf("%s.json", time.Now().Format(time.RFC3339)[:19]), ":", "", -1)
}

func getPersistenceAlgorithm(algoName string) func(res []CutResult, outputPath string) error {
	var impl func(res []CutResult, outputPath string) error
	switch algoName {
	case "standard":
		impl = StandardJSONPersistence
	case "empty_cuts_optimized":
		impl = EmptyCutOptimizedJSONPersistence
	default:
		impl = StandardJSONPersistence
	}
	return impl
}

func GetOpts() (ProgramOpts, error) {
	var opts ProgramOpts
	var dataPath string
	var schemaPath string
	var orgStructurePath string
	var workloadPath string
	var outputPath string
	var concurrentIndex bool
	var concurrentWorkload bool
	var indexImpl func(org OrgStructure, dataNodes []string) OrgNodeIndex
	var workloadImpl func(cuts []Cut, survey *Survey) []CutResult
	var persistResultsImplName string
	flag.StringVar(&dataPath, "data", "", "Path to CSV with survey data.")
	flag.StringVar(&schemaPath, "schema", "", "Path to JSON file with survey schema.")
	flag.StringVar(&orgStructurePath, "org_structure", "", "Path to CSV file with org structure nodes.")
	flag.StringVar(&workloadPath, "workload", "", "Path to JSON file containing workload definition (cuts).")
	flag.BoolVar(&concurrentIndex, "concurrent_index", false, "Use concurrency for building the index.")
	flag.BoolVar(&concurrentWorkload, "concurrent_workload", false, "Use concurrency for processing workload.")
	flag.StringVar(&outputPath, "output_path", "", "Path to JSON file to which results will be saved.")
	flag.StringVar(&persistResultsImplName, "persistence_algorithm", "", "Name for the algorithm for results persistence.")
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

	if outputPath == "" {
		outputPath = getCurrentTimeResultsFilename()
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

	persistResultsImpl := getPersistenceAlgorithm(persistResultsImplName)

	return ProgramOpts{
		DataPath:           dataPath,
		SchemaPath:         schemaPath,
		OrgStructurePath:   orgStructurePath,
		WorkloadPath:       workloadPath,
		IndexAlgorithm:     indexImpl,
		WorkloadAlgorithm:  workloadImpl,
		OutputPath:         outputPath,
		ResultsPersistence: persistResultsImpl,
	}, nil
}
