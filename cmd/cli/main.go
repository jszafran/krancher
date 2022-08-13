package main

import (
	"krancher/survey"
	"log"
)

func main() {
	opts, err := survey.ParseOpts()
	if err != nil {
		log.Fatal(err)
	}

	err = survey.RunWorkloadFromOpts(opts)
	if err != nil {
		log.Fatal(err)
	}
}
