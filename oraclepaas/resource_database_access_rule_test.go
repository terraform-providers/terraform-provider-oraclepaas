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
						resourceName, "enabled", "false"),
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
						resourceName, "enabled", "false"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseAccessRuleExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "test-access-rule"),
					resource.TestCheckResourceAttr(
						resourceName, "enabled", "true"),
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

func testAccDatabaseAccessRuleBasic(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_database_service_instance" "test" {
	name        = "test-service-instance-%d"
	description = "test service instance"
	edition = "EE"
	level = "PAAS"
	shape = "oc3"
	subscription_type = "HOURLY"
	version = "12.2.0.1"
	ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="

	database_configuration {
		admin_password = "Test_String7"
		backup_destination = "NONE"
		sid = "ORCL"
		usable_storage = 15
	}
}

resource "oraclepaas_database_access_rule" "test" {
	name = "test-access-rule-%d"
	service_instance_id = "${oraclepaas_database_service_instance.test.name}"
	description = "test-access-rule"
	ports = "8000"
	source = "PUBLIC-INTERNET"
	enabled = false
}
`, rInt, rInt)
}

func testAccDatabaseAccessRuleUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_database_service_instance" "test" {
	name        = "test-service-instance-%d"
	description = "test service instance"
	edition = "EE"
	level = "PAAS"
	shape = "oc3"
	subscription_type = "HOURLY"
	version = "12.2.0.1"
	ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="

	database_configuration {
		admin_password = "Test_String7"
		backup_destination = "NONE"
		sid = "ORCL"
		usable_storage = 15
	}
}

resource "oraclepaas_database_access_rule" "test" {
	name = "test-access-rule-%d"
	service_instance_id = "${oraclepaas_database_service_instance.test.name}"
	description = "test-access-rule"
	ports = "8000"
	source = "PUBLIC-INTERNET"
	enabled = true
}
`, rInt, rInt)
}
