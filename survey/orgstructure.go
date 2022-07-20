package survey

import (
	"encoding/csv"
	"os"
	"sort"
)

type OrgStructure struct {
	nodes []string
}

func (orgStr *OrgStructure) sortNodesAsc() {
	sort.Strings(orgStr.nodes)
}

func ReadOrgStructureFromCSV(path string, hasHeader bool) OrgStructure {
	file, _ := os.Open(path)
	defer file.Close()

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
