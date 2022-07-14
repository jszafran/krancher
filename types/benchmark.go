package types

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func RunBenchmark(orgNodesPath string, dataNodesPath string) {
	orgStructure := ReadOrgStructureFromCSV(orgNodesPath, true)
	dataNodes := make([]string, 0)
	now := time.Now()
	f, _ := os.Open(dataNodesPath)
	defer f.Close()
	csvReader := csv.NewReader(f)
	lines, _ := csvReader.ReadAll()
	for _, line := range lines {
		dataNodes = append(dataNodes, line[0])
	}
	_ = buildIndex(orgStructure, dataNodes)
	elapsed := time.Since(now)
	fmt.Println("Benchmark done!")
	log.Printf("Building index took %s", elapsed)
}
