package oraclepaas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOPAASMySQLAccessRule_Basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccMySQLAccessRuleBasic(ri)
	resourceName := "oraclepaas_mysql_access_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMySQLAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMySQLAccessRuleExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "test-access-rule"),
					resource.TestCheckResourceAttr(
						resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccOPAASMySQLAccessRule_Update(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccMySQLAccessRuleBasic(ri)
	config2 := testAccMySQLAccessRuleUpdate(ri)
	resourceName := "oraclepaas_mysql_access_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMySQLAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMySQLAccessRuleExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "test-access-rule"),
					resource.TestCheckResourceAttr(
						resourceName, "enabled", "false"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMySQLAccessRuleExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "test-access-rule"),
					resource.TestCheckResourceAttr(
						resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckMySQLAccessRuleExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPAASClient).mysqlClient.AccessRules()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_mysql_access_rule" {
			continue
		}

		input := mysql.GetAccessRuleInput{
			Name:              rs.Primary.Attributes["name"],
			ServiceInstanceID: rs.Primary.Attributes["service_instance_id"],
		}
		if _, err := client.GetAccessRule(&input); err != nil {
			return fmt.Errorf("Error retrieving state of MySQL Access Rule %q: %+v", input.Name, err)
		}
	}

	return nil
}

func testAccCheckMySQLAccessRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPAASClient).mysqlClient.AccessRules()
	if client == nil {
		return fmt.Errorf("MySQL Client is not initialized. Make sure to use `mysql_endpoint` variable or `OPAAS_DATABASE_ENDPOINT` env variable")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_mysql_access_rule" {
			continue
		}

		input := mysql.GetAccessRuleInput{
			Name:              rs.Primary.Attributes["name"],
			ServiceInstanceID: rs.Primary.Attributes["service_instance_id"],
		}
		if info, err := client.GetAccessRule(&input); err == nil && info != nil {
			return fmt.Errorf("MySQL Access Rule %q still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccMySQLAccessRuleBasic(rInt int) string {

	return fmt.Sprintf(`
resource "oraclepaas_mysql_service_instance" "test" {
	description					= "Test Service Instance with Storage"
	name                    	= "TestInst%d"
	ssh_public_key				= "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC0Pspsfu8lUTxILGf+dJnTTbIeFZrL/NKaQNNEvH9jF9aXcr347C5dKlu45LE2jTB8OfjtaExOznn7kKiOErwWPJUzDncDDsmUacDzs5KGbDBGQb6zxEMyYgYCKDiru5V24CrZqam+3QP5AurLopD3JaYmZSikKgP+syu16jBs3WzRLvGzDknIkrUk6t7XjzJ5X/wgMTqepjDDyn9NJ3nG5l4iQe7ULgAbfnRjTM3pRQZ5EM67iN3jc+cIFeNsEwqnxb9ZCJ7avb+Yqdcm/7A5tlX+rMwnTYYCPF/j8bgFdHuO9VHEiQHkM7FuRvZGWkXCryyg9iLM+myG5XdVa3Z2IsfBx3qIfxKMcWsHIk5mmDvWIDbgvBne6JSPKhkB7qM6F10pJSVvt08tGwmlTxZZJPKCkpd0nrfrVChMdMr9yRoYH46bqwMbPFCffNeVkJfj4IMlSSU+A9RGLLEnkdv+Xk3yCS+8RcNA6Zilv9VnJm4hBEJ2LsDVZfwqTvUAeB4evpOCMS+v4YKn/w+R4cB/+SdYDtifBwKW8TYk4ZK3J4wHa6XAI4u3b9C0bIfUmXZs36Gyy4MArtg6QGqrmTzYMa5eI2uB7BnO0JM/Moref8vvQYvGjbnkC5G/yCoLswbt477Gn+Ih96PyZ81qMmTv8qE9S3F3qCqkR3sDJA3oDw=="
	backup_destination      	= "BOTH"
	shape                   	= "oc3"
	
	backups {
		cloud_storage_container =  "%sacctest-%d"
		create_if_missing 		= true				
	}
			
	mysql_configuration {
		db_name                 = "demo_db"
		db_storage              = 25
		mysql_port              = 3306
		mysql_username          = "root"
		mysql_password          = "MySqlPassword_1"
		mysql_charset           = "utf8"
		mysql_collation         = "utf8_general_ci"	    
	}
}

resource "oraclepaas_mysql_access_rule" "test" {
	name                = "TestRule%d"
	service_instance_id = "${oraclepaas_mysql_service_instance.test.name}"
	description         = "test-access-rule"
	ports               = "8000"
	source              = "0.0.0.0/24"
	destination         = "mysql_MASTER"
	enabled             = false
}`, rInt, os.Getenv("OPC_STORAGE_URL"), rInt, rInt)
}

func testAccMySQLAccessRuleUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_mysql_service_instance" "test" {
	description					= "Test Service Instance with Storage"
	name                    	= "TestInst%d"
	ssh_public_key				= "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC0Pspsfu8lUTxILGf+dJnTTbIeFZrL/NKaQNNEvH9jF9aXcr347C5dKlu45LE2jTB8OfjtaExOznn7kKiOErwWPJUzDncDDsmUacDzs5KGbDBGQb6zxEMyYgYCKDiru5V24CrZqam+3QP5AurLopD3JaYmZSikKgP+syu16jBs3WzRLvGzDknIkrUk6t7XjzJ5X/wgMTqepjDDyn9NJ3nG5l4iQe7ULgAbfnRjTM3pRQZ5EM67iN3jc+cIFeNsEwqnxb9ZCJ7avb+Yqdcm/7A5tlX+rMwnTYYCPF/j8bgFdHuO9VHEiQHkM7FuRvZGWkXCryyg9iLM+myG5XdVa3Z2IsfBx3qIfxKMcWsHIk5mmDvWIDbgvBne6JSPKhkB7qM6F10pJSVvt08tGwmlTxZZJPKCkpd0nrfrVChMdMr9yRoYH46bqwMbPFCffNeVkJfj4IMlSSU+A9RGLLEnkdv+Xk3yCS+8RcNA6Zilv9VnJm4hBEJ2LsDVZfwqTvUAeB4evpOCMS+v4YKn/w+R4cB/+SdYDtifBwKW8TYk4ZK3J4wHa6XAI4u3b9C0bIfUmXZs36Gyy4MArtg6QGqrmTzYMa5eI2uB7BnO0JM/Moref8vvQYvGjbnkC5G/yCoLswbt477Gn+Ih96PyZ81qMmTv8qE9S3F3qCqkR3sDJA3oDw=="
	backup_destination      	= "BOTH"
	shape                   	= "oc3"

	backups {
		cloud_storage_container =  "%sacctest-%d"
		create_if_missing 		= true				
	}
			
	mysql_configuration {
		db_name                 = "demo_db"
		db_storage              = 25
		mysql_port              = 3306
		mysql_username          = "root"
		mysql_password          = "MySqlPassword_1"
		mysql_charset           = "utf8"
		mysql_collation         = "utf8_general_ci"	    
	}
}

resource "oraclepaas_mysql_access_rule" "test" {
	name                = "TestRule%d"
	service_instance_id = "${oraclepaas_mysql_service_instance.test.name}"
	description         = "test-access-rule"
	ports               = "8000"
	source              = "0.0.0.0/24"
	destination         = "mysql_MASTER"
	enabled             = true
}`, rInt, os.Getenv("OPC_STORAGE_URL"), rInt, rInt)
}
