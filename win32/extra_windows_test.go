package win32_test

import (
	"fmt"
	"log"

	"github.com/taskcluster/runlib/win32"
)

func ExampleMergeEnvLists() {
	lists := []*[]string{
		{
			"a=dog",
			"Pete=man",
			"x=ray",
		}, {
			"food=good",
			"PETE=person",
		},
	}
	res, err := win32.MergeEnvLists(lists...)
	if err != nil {
		log.Fatalf("Hit error: %v", err)
	}
	fmt.Println(*res)
	// Output:
	// hfjdhfd
}
