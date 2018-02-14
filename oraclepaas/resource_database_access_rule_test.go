package oraclepaas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOPAASDatabaseAccessRule_Basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseAccessRuleBasic(ri)
	resourceName := "oraclepaas_database_access_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseAccessRuleExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "test-access-rule"),
					resource.TestCheckResourceAttr(
						resourceName, "status", "disabled"),
				),
			},
		},
	})
}

func TestAccOPAASDatabaseAccessRule_Update(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseAccessRuleBasic(ri)
	config2 := testAccDatabaseAccessRuleUpdate(ri)
	resourceName := "oraclepaas_database_access_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseAccessRuleExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "test-access-rule"),
					resource.TestCheckResourceAttr(
						resourceName, "status", "disabled"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseAccessRuleExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "test-access-rule"),
					resource.TestCheckResourceAttr(
						resourceName, "status", "enabled"),
				),
			},
		},
	})
}

func testAccCheckDatabaseAccessRuleExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPAASClient).databaseClient.AccessRules()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_database_access_rule" {
			continue
		}

		input := database.GetAccessRuleInput{
			Name:              rs.Primary.Attributes["name"],
			ServiceInstanceID: rs.Primary.Attributes["service_instance_id"],
		}
		if _, err := client.GetAccessRule(&input); err != nil {
			return fmt.Errorf("Error retrieving state of Database Access Rule %q: %+v", input.Name, err)
		}
	}

	return nil
}

func testAccCheckDatabaseAccessRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPAASClient).databaseClient.AccessRules()
	if client == nil {
		return fmt.Errorf("Database Client is not initialized. Make sure to use `database_endpoint` variable or `OPAAS_DATABASE_ENDPOINT` env variable")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_database_access_rule" {
			continue
		}

		input := database.GetAccessRuleInput{
			Name:              rs.Primary.Attributes["name"],
			ServiceInstanceID: rs.Primary.Attributes["service_instance_id"],
		}
		if info, err := client.GetAccessRule(&input); err == nil && info != nil {
			return fmt.Errorf("Database Access Rule %q still exists: %#v", input.Name, info)
		}
	}

	return nil
}

// TODO add database service instance
func testAccDatabaseAccessRuleBasic(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_database_access_rule" "test" {
	name = "test-access-rule-%d"
	service_instance_id = "matthew-test2"
	description = "test-access-rule"
	ports = "8000"
	source = "PUBLIC-INTERNET"
	status = "disabled"
}
`, rInt)
}

func testAccDatabaseAccessRuleUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_database_access_rule" "test" {
	name = "test-access-rule-%d"
	service_instance_id = "matthew-test2"
	description = "test-access-rule"
	ports = "8000"
	source = "PUBLIC-INTERNET"
	status = "enabled"
}
`, rInt)
}
