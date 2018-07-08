package helper

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-oracle-terraform/opc"
)

// TestEnvVar is the constant to determine whether to run acceptance tests
const TestEnvVar = "ORACLE_ACC"

// TestCase Test suite helpers
type TestCase struct {
	// Fields to test stuff with
}

// Test sets up the test framework
func Test(t TestT, c TestCase) {
	if os.Getenv(TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' is set", TestEnvVar))
		return
	}

	// Setup logging Output
	logWriter, err := opc.LogOutput()
	if err != nil {
		t.Error(fmt.Sprintf("Error setting up log writer: %s", err))
	}
	log.SetOutput(logWriter)
}

// TestT supports errors, fatals, and skips for tests
type TestT interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Skip(args ...interface{})
}
