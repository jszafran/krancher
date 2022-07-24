package survey

import (
	"errors"
	"fmt"
)

func orgStructureContains(org OrgStructure, node string) bool {
	contains := false
	for _, orgNode := range org.nodes {
		if orgNode == node {
			contains = true
		}
	}
	return contains
}

func demogsContains(schema Schema, demog string) bool {
	contains := false
	for _, d := range schema.GetDemographicsCodes() {
		if d == demog {
			contains = true
		}
	}
	return contains
}

func ValidateWorkload(wrkl Workload, org OrgStructure, schema Schema) []error {
	// check if there are cuts pointing to a non-existing org node
	errs := make([]error, 0)

	for _, cut := range wrkl.Cuts {
		if !orgStructureContains(org, cut.OrgNode) {
			err := errors.New(fmt.Sprintf("%s: %s", CutPointingToNonExistingOrgNode, cut.OrgNode))
			errs = append(errs, err)
		}
	}

	// check if there are any demographics that are not included in survey schema

	for _, demog := range wrkl.GetDemographicsSet() {
		if !demogsContains(schema, demog) {

			errs = append(errs, errors.New(""))
		}
	}
	return errs
}
