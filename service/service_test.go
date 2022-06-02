package service_test

import (
	"fmt"
	"testing"

	"github.com/plgd-dev/client-application/test"
)

func TestServiceServe(t *testing.T) {
	fmt.Printf("%v\n\n", test.MakeConfig(t))

	shutDown := test.SetUp(t)
	defer shutDown()
}
