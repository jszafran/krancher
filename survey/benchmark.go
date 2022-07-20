package survey

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func RunBenchmark(orgNodesPath string, dataNodesPath string) {
	orgStructure := ReadOrgStructureFromCSV(orgNodesPath, true)
	dataNodes := make([]string, 0)
	startTime := time.Now()
	f, _ := os.Open(dataNodesPath)
	defer f.Close()
	csvReader := csv.NewReader(f)
	lines, _ := csvReader.ReadAll()
	fmt.Printf("Benchmark will be ran on %d rows of data", len(lines)-1)
	for _, line := range lines {
		dataNodes = append(dataNodes, line[0])
	}

	sortStartTime := time.Now()
	sort.SliceStable(dataNodes, func(i, j int) bool {
		excludeEmpty := func(x []string) []string {
			if x[len(x)-1] == "" {
				return x[:len(x)-1]
			}
			return x
		}
		str1 := strings.Replace(dataNodes[i], "N", "", -1)
		str2 := strings.Replace(dataNodes[j], "N", "", -1)
		spl1 := excludeEmpty(strings.Split(str1, "."))
		spl2 := excludeEmpty(strings.Split(str2, "."))
		mlen := int(math.Min(float64(len(spl1)), float64(len(spl2))))
		for i := 0; i < mlen; i++ {
			n1, err1 := strconv.Atoi(spl1[i])
			n2, err2 := strconv.Atoi(spl2[i])
			if err1 != nil || err2 != nil {
				log.Fatalf("Failed to convert values %v or %v to int.", spl1[i], spl2[i])
			}
			if n1 == n2 {
				continue
			}
			return n1 < n2
		}
		return false
	})
	sortingElapsed := time.Since(sortStartTime)
	indexBuildStartTime := time.Now()
	_ = buildIndex(orgStructure, dataNodes)
	indexBuildElapsed := time.Since(indexBuildStartTime)
	totalElapsed := time.Since(startTime)
	fmt.Println("Benchmark done!")
	fmt.Printf("Sorting data nodes took %s\n", sortingElapsed)
	log.Printf("Building index took %s\n", indexBuildElapsed)
	log.Printf("Total time: %s\n", totalElapsed)
}
