package survey

import (
	"encoding/csv"
	"log"
	"os"
)

type OrgStructure struct {
	nodes []string
}

func ReadOrgStructureFromCSV(path string, hasHeader bool) OrgStructure {
	file, _ := os.Open(path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	csvReader := csv.NewReader(file)

	lines, _ := csvReader.ReadAll()
	nodes := make([]string, 0)

	for i, line := range lines {
		if i == 0 && hasHeader {
			continue
		}
		nodes = append(nodes, line[0])
	}
	return OrgStructure{nodes}
}
