package main

import (
	"fmt"
	"krancher/types"
)

func main() {
	p := "resources/org.csv"
	org1 := types.ReadOrgStructureFromCSV(p, false)
	org2 := types.ReadOrgStructureFromCSV(p, true)
	fmt.Println(org1)
	fmt.Println(org2)

}
