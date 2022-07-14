package main

import "krancher/types"

func main() {
	types.RunBenchmark("resources/bigger_org_structure.csv", "resources/data_nodes_100x.csv")
}
