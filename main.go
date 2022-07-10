package main

import (
	"fmt"
	"krancher/kio"
)

func main() {
	p := "resources/example_schema.json"
	sch := kio.SchemaFromJSON(p)
	fmt.Println(sch)
}
