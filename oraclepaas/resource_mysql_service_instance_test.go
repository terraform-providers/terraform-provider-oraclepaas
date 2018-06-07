package oraclepaas

import (
	"fmt"
	"github.com/hashicorp/go-oracle-terraform/mysql"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccOraclePAASMySQLServiceInstance_EnterpriseMonitor(t *testing.T) {

	ri := acctest.RandInt()
	config := testMySQLServiceInstanceEnterpriseMonitor(ri)
	resourceName := "oraclepaas_mysql_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMySQLServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMySQLServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "Test Service Instance with EM"),
					resource.TestCheckResourceAttr(
						resourceName, "mysql_configuration.0.enterprise_monitor_configuration.#", "1")),
			},
		},
	})
}

func TestAccOPAASMySQLServiceInstance_CloudStorage(t *testing.T) {

	storageEndpoint := os.Getenv("TEST_STORAGE_ENDPOINT")
	storageUsername := os.Getenv("TEST_STORAGE_USERNAME")
	storagePassword := os.Getenv("TEST_STORAGE_PASSWORD")

	if storageEndpoint == "" {
		t.Skip("You will need to set the cloud storage container environment parameter `TEST_STORAGE_ENDPOINT` to run this test. Skipping...")
	}

	if storageUsername == "" {
		t.Skip("You will need to set the cloud storage username environment parameter `TEST_STORAGE_USERNAME` to run this test. Skipping...")
	}

	if storagePassword == "" {
		t.Skip("You will need to set the cloud storage password environment parameter `TEST_STORAGE_PASSWORD` to run this test. Skipping...")
	}

	ri := acctest.RandInt()
	container := fmt.Sprintf("%s/acc-test-%d", storageEndpoint, ri)
	config := testAccMySQLServiceInstanceCloudStorage(ri, container, storageUsername, storagePassword)
	resourceName := "oraclepaas_mysql_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMySQLServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMySQLServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "Test Service Instance with Storage"),
					resource.TestCheckResourceAttr(
						resourceName, "backups.0.cloud_storage_container", container),
					resource.TestCheckResourceAttr(
						resourceName, "mysql_configuration.0.enterprise_monitor_configuration.#", "0")),
			},
		},
	})
}

/* Test with OCI.
 */
func TestAccOPAASMySQLServiceInstance_OCI(t *testing.T) {

	oci_region := os.Getenv("TEST_OCI_REGION")
	oci_availability_domain := os.Getenv("TEST_OCI_AD")
	oci_subnet := os.Getenv("TEST_OCI_SUBNET")

	if oci_region == "" {
		t.Skip("Missing Environment Parameter `TEST_OCI_REGION`. You will need to set the environment parameters `TEST_OCI_REGION`, `TEST_OCI_AD` and `TEST_OCI_SUBNET` to run this test.")
	}

	if oci_availability_domain == "" {
		t.Skip("Missing Environment Parameter `TEST_OCI_AD`. You will need to set the environment parameters `TEST_OCI_REGION`, `TEST_OCI_AD` and `TEST_OCI_SUBNET` to run this test.")
	}

	if oci_subnet == "" {
		t.Skip("Missing Environment Parameter `TEST_OCI_SUBNET`. You will need to set the environment parameters `TEST_OCI_REGION`, `TEST_OCI_AD` and `TEST_OCI_SUBNET` to run this test.")
	}

	ri := acctest.RandInt()
	config := testAccMySQLServiceInstanceOCI(ri, oci_region, oci_availability_domain, oci_subnet)
	t.Logf("Config : %s", config)

	resourceName := "oraclepaas_mysql_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMySQLServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMySQLServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "region", oci_region),
					resource.TestCheckResourceAttr(
						resourceName, "availability_domain", oci_availability_domain),
					resource.TestCheckResourceAttr(
						resourceName, "mysql_configuration.0.enterprise_monitor_configuration.#", "0")),
			},
		},
	})
}

func testAccCheckMySQLServiceInstanceExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPAASClient).mysqlClient.ServiceInstanceClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_mysql_service_instance" {
			continue
		}

		input := mysql.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}

		if _, err := client.GetServiceInstance(&input); err != nil {
			return fmt.Errorf("Error retrieving state of MySQLServiceInstance %s: %+v", input.Name, err)
		}
	}
	return nil
}

func testAccCheckMySQLServiceInstanceDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*OPAASClient).mysqlClient.ServiceInstanceClient()
	if client == nil {
		return fmt.Errorf("MySQL Client is not initialized. Make sure to use `mysql_endpoint` variable or `OPAAS_MYSQL_ENDPOINT` env variable")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_mysql_service_instance" {
			continue
		}

		input := mysql.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}

		if info, err := client.GetServiceInstance(&input); err == nil {
			return fmt.Errorf("MySQLServiceInstance %s (%v) still exists: %v", input.Name, input, info)
		}
	}
	return nil
}

func testMySQLServiceInstanceEnterpriseMonitor(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_mysql_service_instance" "test" {
	description				  = "Test Service Instance with EM"
	name                      = "TestInst%d"
	vm_public_key             = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC0Pspsfu8lUTxILGf+dJnTTbIeFZrL/NKaQNNEvH9jF9aXcr347C5dKlu45LE2jTB8OfjtaExOznn7kKiOErwWPJUzDncDDsmUacDzs5KGbDBGQb6zxEMyYgYCKDiru5V24CrZqam+3QP5AurLopD3JaYmZSikKgP+syu16jBs3WzRLvGzDknIkrUk6t7XjzJ5X/wgMTqepjDDyn9NJ3nG5l4iQe7ULgAbfnRjTM3pRQZ5EM67iN3jc+cIFeNsEwqnxb9ZCJ7avb+Yqdcm/7A5tlX+rMwnTYYCPF/j8bgFdHuO9VHEiQHkM7FuRvZGWkXCryyg9iLM+myG5XdVa3Z2IsfBx3qIfxKMcWsHIk5mmDvWIDbgvBne6JSPKhkB7qM6F10pJSVvt08tGwmlTxZZJPKCkpd0nrfrVChMdMr9yRoYH46bqwMbPFCffNeVkJfj4IMlSSU+A9RGLLEnkdv+Xk3yCS+8RcNA6Zilv9VnJm4hBEJ2LsDVZfwqTvUAeB4evpOCMS+v4YKn/w+R4cB/+SdYDtifBwKW8TYk4ZK3J4wHa6XAI4u3b9C0bIfUmXZs36Gyy4MArtg6QGqrmTzYMa5eI2uB7BnO0JM/Moref8vvQYvGjbnkC5G/yCoLswbt477Gn+Ih96PyZ81qMmTv8qE9S3F3qCqkR3sDJA3oDw=="
	backup_destination        = "NONE"

	mysql_configuration = {
		db_name               = "demo_db"
		db_storage            = 25
		mysql_port            = 3306
		mysql_username        = "root"
		mysql_password        = "MySqlPassword_1"
		shape                 = "oc3"
		mysql_charset         = "utf8"
		mysql_collation       = "utf8_general_ci"

		enterprise_monitor_configuration = {
			em_agent_password = "MySqlPassword_1"
			em_agent_username = "admin"
			em_password 	  = "MySqlPassword_1"
			em_username 	  = "admin"
			em_port 		  = "18443"
		}
	}
}
`, rInt)
}

func testAccMySQLServiceInstanceCloudStorage(rInt int, container string, username string, password string) string {
	return fmt.Sprintf(`
resource "oraclepaas_mysql_service_instance" "test" {
		
	description					= "Test Service Instance with Storage"
	name                    	= "TestInst%d"
	vm_public_key				= "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC0Pspsfu8lUTxILGf+dJnTTbIeFZrL/NKaQNNEvH9jF9aXcr347C5dKlu45LE2jTB8OfjtaExOznn7kKiOErwWPJUzDncDDsmUacDzs5KGbDBGQb6zxEMyYgYCKDiru5V24CrZqam+3QP5AurLopD3JaYmZSikKgP+syu16jBs3WzRLvGzDknIkrUk6t7XjzJ5X/wgMTqepjDDyn9NJ3nG5l4iQe7ULgAbfnRjTM3pRQZ5EM67iN3jc+cIFeNsEwqnxb9ZCJ7avb+Yqdcm/7A5tlX+rMwnTYYCPF/j8bgFdHuO9VHEiQHkM7FuRvZGWkXCryyg9iLM+myG5XdVa3Z2IsfBx3qIfxKMcWsHIk5mmDvWIDbgvBne6JSPKhkB7qM6F10pJSVvt08tGwmlTxZZJPKCkpd0nrfrVChMdMr9yRoYH46bqwMbPFCffNeVkJfj4IMlSSU+A9RGLLEnkdv+Xk3yCS+8RcNA6Zilv9VnJm4hBEJ2LsDVZfwqTvUAeB4evpOCMS+v4YKn/w+R4cB/+SdYDtifBwKW8TYk4ZK3J4wHa6XAI4u3b9C0bIfUmXZs36Gyy4MArtg6QGqrmTzYMa5eI2uB7BnO0JM/Moref8vvQYvGjbnkC5G/yCoLswbt477Gn+Ih96PyZ81qMmTv8qE9S3F3qCqkR3sDJA3oDw=="
	backup_destination      	= "BOTH"
	
	backups {
		cloud_storage_container = "%s"
		cloud_storage_username  = "%s"
		cloud_storage_password  = "%s"
		create_if_missing 		= true				
	}
			
	mysql_configuration = {
		db_name                 = "demo_db"
		db_storage              = 25
		mysql_port              = 3306
		mysql_username          = "root"
		mysql_password          = "MySqlPassword_1"
		shape                   = "oc3"
		mysql_charset           = "utf8"
		mysql_collation         = "utf8_general_ci"	    
	}
}`, rInt, container, username, password)
}

func testAccMySQLServiceInstanceOCI(rInt int, oci_region string, oci_availability_domain string, oci_subnet string) string {

	return fmt.Sprintf(`
resource "oraclepaas_mysql_service_instance" "test" {
		
	description			= "Test Service Instance Creation on OCI"
	name                = "TestInst%d"
	vm_public_key       = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC0Pspsfu8lUTxILGf+dJnTTbIeFZrL/NKaQNNEvH9jF9aXcr347C5dKlu45LE2jTB8OfjtaExOznn7kKiOErwWPJUzDncDDsmUacDzs5KGbDBGQb6zxEMyYgYCKDiru5V24CrZqam+3QP5AurLopD3JaYmZSikKgP+syu16jBs3WzRLvGzDknIkrUk6t7XjzJ5X/wgMTqepjDDyn9NJ3nG5l4iQe7ULgAbfnRjTM3pRQZ5EM67iN3jc+cIFeNsEwqnxb9ZCJ7avb+Yqdcm/7A5tlX+rMwnTYYCPF/j8bgFdHuO9VHEiQHkM7FuRvZGWkXCryyg9iLM+myG5XdVa3Z2IsfBx3qIfxKMcWsHIk5mmDvWIDbgvBne6JSPKhkB7qM6F10pJSVvt08tGwmlTxZZJPKCkpd0nrfrVChMdMr9yRoYH46bqwMbPFCffNeVkJfj4IMlSSU+A9RGLLEnkdv+Xk3yCS+8RcNA6Zilv9VnJm4hBEJ2LsDVZfwqTvUAeB4evpOCMS+v4YKn/w+R4cB/+SdYDtifBwKW8TYk4ZK3J4wHa6XAI4u3b9C0bIfUmXZs36Gyy4MArtg6QGqrmTzYMa5eI2uB7BnO0JM/Moref8vvQYvGjbnkC5G/yCoLswbt477Gn+Ih96PyZ81qMmTv8qE9S3F3qCqkR3sDJA3oDw=="
	backup_destination  = "NONE"
	region              = "%s"
	availability_domain	= "%s"

	mysql_configuration = {
		db_name			= "demo_db"
		db_storage      = 25
		mysql_port      = 3306
		mysql_username  = "root"
		mysql_password  = "MySqlPassword_1"
		shape           = "oc3"
		mysql_charset   = "utf8"
		mysql_collation = "utf8_general_ci"	    
		subnet          = "%s"
	}
}`, rInt, oci_region, oci_availability_domain, oci_subnet)
}
