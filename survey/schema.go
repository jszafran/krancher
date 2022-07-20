package survey

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func SchemaFromJSON(path string) Schema {
	// TODO: research what is idiomatic way of handling multiple errors in one func?
	jsonFile, err1 := os.Open(path)
	var schema Schema

	if err1 != nil {
		log.Fatal(err1)
	}

	bytes, err2 := ioutil.ReadAll(jsonFile)

	if err2 != nil {
		log.Fatal(err2)
	}

	err3 := json.Unmarshal(bytes, &schema)

	if err3 != nil {
		log.Fatal(err3)
	}
	defer jsonFile.Close()

	return schema
}