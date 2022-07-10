package main

import (
	"fmt"
	"krancher/kio"
)

func main() {
	p := "resources/example_schema.jsozn"
	sch := kio.SchemaFromJSON(p)
	fmt.Println(sch)
}
