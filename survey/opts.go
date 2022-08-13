package survey

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"
)

type WorkloadOpts struct {
	DataPath         string
	SchemaPath       string
	OrgStructurePath string
	WorkloadPath     string
	OutputPath       string
	IndexBuilder     IndexBuilder
	Processor        Processor
	Persistor        Persistor
}

func getCurrentTimeResultsFilename() string {
	return strings.Replace(fmt.Sprintf("%s.json", time.Now().Format(time.RFC3339)[:19]), ":", "", -1)
}

func ProcessorAlgorithmFactory(algoName string) (Processor, error) {
	var f Processor
	algs := map[string]Processor{
		"sequential": SequentialCutProcessor,
		"concurrent": ConcurrentCutProcessor,
	}

	alg, exists := algs[algoName]
	if !exists {
		return f, fmt.Errorf("%s processor algorithm not supported", algoName)
	}
	return alg, nil
}

func PersistenceAlgorithmFactory(algoName string) (Persistor, error) {
	var f Persistor
	algs := map[string]Persistor{
		"standard":               StandardJSONPersistor,
		"skip_empty_cuts_counts": SkipEmptyCutsCountsJSONPersistor,
	}
	alg, exists := algs[algoName]
	if !exists {
		return f, fmt.Errorf("%s persistence algorithm not supported", algoName)
	}
	return alg, nil
}

func IndexBuilderAlgorithmFactory(algoName string) (IndexBuilder, error) {
	var f IndexBuilder
	algs := map[string]IndexBuilder{
		"standard":   SequentialIndexBuilder,
		"concurrent": ConcurrentIndexBuilder,
	}

	alg, exists := algs[algoName]
	if !exists {
		return f, fmt.Errorf("%s index builder algorithm not supported", algoName)
	}
	return alg, nil
}

func GetOpts() (WorkloadOpts, error) {
	var opts WorkloadOpts
	var dataPath string
	var schemaPath string
	var orgStructurePath string
	var workloadPath string
	var outputPath string
	var indexBuilderAlgName string
	var processorAlgName string
	var persistorAlgName string
	flag.StringVar(&dataPath, "data", "", "Path to CSV with survey data.")
	flag.StringVar(&schemaPath, "schema", "", "Path to JSON file with survey schema.")
	flag.StringVar(&orgStructurePath, "org_structure", "", "Path to CSV file with org structure nodes.")
	flag.StringVar(&workloadPath, "workload", "", "Path to JSON file containing workload definition (cuts).")
	flag.StringVar(&outputPath, "output_path", "", "Path to JSON file to which results will be saved.")
	flag.StringVar(&indexBuilderAlgName, "index_builder", "", "Name of the index builder algorithm.")
	flag.StringVar(&processorAlgName, "processor", "", "Name of the processor algorithm.")
	flag.StringVar(&persistorAlgName, "persistor", "", "Name of the persistor algorithm.")
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

	processor, err := ProcessorAlgorithmFactory(processorAlgName)
	if err != nil {
		return opts, err
	}

	indexBuilder, err := IndexBuilderAlgorithmFactory(indexBuilderAlgName)
	if err != nil {
		return opts, err
	}

	persistor, err := PersistenceAlgorithmFactory(persistorAlgName)
	if err != nil {
		return opts, err
	}
	return WorkloadOpts{
		DataPath:         dataPath,
		SchemaPath:       schemaPath,
		OrgStructurePath: orgStructurePath,
		WorkloadPath:     workloadPath,
		IndexBuilder:     indexBuilder,
		Processor:        processor,
		OutputPath:       outputPath,
		Persistor:        persistor,
	}, nil
}
