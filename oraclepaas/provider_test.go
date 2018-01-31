package oraclepaas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"oraclepaas": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("Error creating Provider: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	required := []string{"OPC_USERNAME", "OPC_PASSWORD", "OPC_IDENTITY_DOMAIN", "ORACLEPAAS_DATABASE_ENDPOINT"}
	for _, prop := range required {
		if os.Getenv(prop) == "" {
			t.Fatalf("%s must be set for acceptance test", prop)
		}
	}
	config := Config{
		User:             os.Getenv("OPC_USERNAME"),
		Password:         os.Getenv("OPC_PASSWORD"),
		IdentityDomain:   os.Getenv("OPC_IDENTITY_DOMAIN"),
		MaxRetries:       1,
		Insecure:         false,
		DatabaseEndpoint: os.Getenv("ORACLEPAAS_DATABASE_ENDPOINT"),
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(fmt.Sprintf("%+v", err))
	}
	if client.databaseClient == nil {
		t.Fatalf("Database Client is nil. Make sure your Oracle Cloud Account has access to the Database Cloud")
	}
}

type OPAASResourceState struct {
	*database.DatabaseClient
	*terraform.InstanceState
}

func oraclepaasResourceCheck(resourceName string, f func(checker *OPAASResourceState) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		state := &OPAASResourceState{
			DatabaseClient: testAccProvider.Meta().(*OPAASClient).databaseClient,
			InstanceState:  rs.Primary,
		}

		return f(state)
	}
}
