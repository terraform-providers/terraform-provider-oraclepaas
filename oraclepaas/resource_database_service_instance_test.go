package oraclepaas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOraclePAASDatabaseServiceInstance_Basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceBasic(ri)
	resourceName := "oraclepaas_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "description", "test service instance"),
					resource.TestCheckResourceAttr(
						resourceName, "edition", "EE"),
					resource.TestCheckResourceAttr(
						resourceName, "version", "12.2.0.1"),
				),
			},
		},
	})
}

func TestAccOPAASDatabaseServiceInstance_CloudStorage(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceCloudStorage(ri)
	resourceName := "oraclepaas_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "cloud_storage_container", fmt.Sprintf("%sacctest-%d", os.Getenv("OPC_STORAGE_URL"), ri)),
				),
			},
		},
	})
}

func TestAccOPAASDatabaseServiceInstance_FromBackup(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceFromBackup(ri)
	resourceName := "oraclepaas_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "instantiate_from_backup.0.cloud_storage_container", fmt.Sprintf("%stest-db-java-instance", os.Getenv("OPC_STORAGE_URL"))),
				),
			},
		},
	})
}

func TestAccOPAASDatabaseServiceInstance_DefaultAccessRule(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceDefaultAccessRule(ri)
	resourceName := "oraclepaas_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "default_access_rules.0.enable_ssh", "false"),
				),
			},
		},
	})
}

func TestAccOraclePAASDatabaseServiceInstance_DesiredState(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceDesiredState(ri)
	resourceName := "oraclepaas_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "status", "Stopped"),
				),
			},
		},
	})
}

func TestAccOPAASDatabaseServiceInstance_UpdateShape(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceBasic(ri)
	config2 := testAccDatabaseServiceInstanceUpdateShape(ri)
	resourceName := "oraclepaas_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "shape", "oc3"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "shape", "oc4"),
				),
			},
		},
	})
}

// An OCI account is need to test this
/*
func TestAccOPAASDatabaseServiceInstance_HDG(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceHDG(ri)
	resourceName := "oraclepaas_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "hybrid_disaster_recovery.0.cloud_storage_container", fmt.Sprintf("Storage-%s/test-db-java-instance", os.Getenv("OPC_IDENTITY_DOMAIN"))),
				),
			},
		},
	})
} */

func testAccCheckDatabaseServiceInstanceExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPAASClient).databaseClient.ServiceInstanceClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_database_service_instance" {
			continue
		}

		input := database.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetServiceInstance(&input); err != nil {
			return fmt.Errorf("Error retrieving state of DatabaseServiceInstance %s: %+v", input.Name, err)
		}
	}

	return nil
}

func testAccCheckDatabaseServiceInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPAASClient).databaseClient.ServiceInstanceClient()
	if client == nil {
		return fmt.Errorf("Database Client is not initialized. Make sure to use `database_endpoint` variable or `OPAAS_DATABASE_ENDPOINT` env variable")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_database_service_instance" {
			continue
		}

		input := database.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetServiceInstance(&input); err == nil {
			return fmt.Errorf("DatabaseServiceInstance %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccDatabaseServiceInstanceBasic(rInt int) string {
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
    bring_your_own_license = true

    database_configuration {
        admin_password = "Test_String7"
        backup_destination = "NONE"
        sid = "ORCL"
        usable_storage = 15
    }
}`, rInt)
}

func testAccDatabaseServiceInstanceUpdateShape(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_database_service_instance" "test" {
    name        = "test-service-instance-%d"
    description = "test service instance"
    edition = "EE"
    level = "PAAS"
    shape = "oc4"
    subscription_type = "HOURLY"
    version = "12.2.0.1"
    ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
    bring_your_own_license = true

    database_configuration {
        admin_password = "Test_String7"
        backup_destination = "NONE"
		sid = "ORCL"
        usable_storage = 15
    }
}`, rInt)
}

func testAccDatabaseServiceInstanceCloudStorage(rInt int) string {
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
    bring_your_own_license = true

    database_configuration {
        admin_password = "Test_String7"
        backup_destination = "OSS"
        failover_database = false
        sid = "ORCL"
        usable_storage = 15
    }

    backups {
        cloud_storage_container = "%sacctest-%d"
        create_if_missing = true
    }
}`, rInt, os.Getenv("OPC_STORAGE_URL"), rInt)
}

func testAccDatabaseServiceInstanceFromBackup(rInt int) string {
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
    bring_your_own_license = true

    database_configuration {
        admin_password = "Test_String7"
        backup_destination = "OSS"
        failover_database = false
        sid = "ORCL"
        usable_storage = 15
        backup_destination = "NONE"
    }

    instantiate_from_backup {
        cloud_storage_container = "%stest-db-java-instance"
        on_premise     = false
        database_id    = "1"
        service_id     = "ORCL"
        decryption_key = "Test_String7"
    }
}`, rInt, os.Getenv("OPC_STORAGE_URL"))
}

func testAccDatabaseServiceInstanceDefaultAccessRule(rInt int) string {
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
    bring_your_own_license = true

    database_configuration {
        admin_password = "Test_String7"
        backup_destination = "NONE"
        sid = "ORCL"
        usable_storage = 15
    }

    default_access_rules {
        enable_ssh = false
    }
}`, rInt)
}

func testAccDatabaseServiceInstanceDesiredState(rInt int) string {
	return fmt.Sprintf(`resource "oraclepaas_database_service_instance" "test" {
    name        = "test-service-instance-%d"
    description = "test service instance"
    edition = "EE"
    level = "PAAS"
    shape = "oc3"
    subscription_type = "HOURLY"
    version = "12.2.0.1"
    ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
    bring_your_own_license = true

    database_configuration {
        admin_password = "Test_String7"
        backup_destination = "NONE"
        sid = "ORCL"
        usable_storage = 15
    }

    desired_state = "stop"
}`, rInt)
}

// An OCI account is need to test this
/*
func testAccDatabaseServiceInstanceHDG(rInt int) string {
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
		backup_destination = "OSS"
		failover_database = false
		sid = "ORCL"
		usable_storage = 15
		backup_destination = "NONE"
	}

	hybrid_disaster_recovery {
		cloud_storage_container = "Storage-%s/test-db-java-instance"
  	}
}`, rInt, os.Getenv("OPC_IDENTITY_DOMAIN"))
} */

func TestValidateSID(t *testing.T) {
	validPrefixes := []string{
		"ORCL",
		"MYsid01",
		"ABCDEFGH",
		"A1234567",
	}

	for _, v := range validPrefixes {
		_, errors := validateSID(v, "sid")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid SID: %q", v, errors)
		}
	}

	invalidPrefixes := []string{
		"",
		"12345678",
		"ThisSIDisTooLong",
		"NO_CHAR",
	}

	for _, v := range invalidPrefixes {
		_, errors := validateSID(v, "sid")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid SID", v)
		}
	}
}

func TestValidatePDBName(t *testing.T) {
	validPrefixes := []string{
		"pdb1",
		"MyPDBExample",
		"ABCDEFGHIJKLMabcdefghijklm1234", // 30 chars
	}

	for _, v := range validPrefixes {
		_, errors := validatePDBName(v, "pdb_name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid PDB Name: %q", v, errors)
		}
	}

	invalidPrefixes := []string{
		"",
		"12345678",
		"ABCDEFGHIJKLMabcdefghijklm12345", // 31 chars
		"no_underscore",
	}

	for _, v := range invalidPrefixes {
		_, errors := validatePDBName(v, "pdb_name")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid PDB Name", v)
		}
	}
}

func TestValidateOracleHome(t *testing.T) {
	validPrefixes := []string{
		"OracleHome",
		"MyOracleHome",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789012", // 64 chars
	}

	for _, v := range validPrefixes {
		_, errors := validateOracleHomeName(v, "pdb_name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid PDB Name: %q", v, errors)
		}
	}

	invalidPrefixes := []string{
		"",
		"12345678",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890123", // 65 chars
		"no_underscore_at_end_",
	}

	for _, v := range invalidPrefixes {
		_, errors := validateOracleHomeName(v, "pdb_name")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid PDB Name", v)
		}
	}
}

func TestValidateAdminPassword(t *testing.T) {
	validPrefixes := []string{
		"Pa55_Word",
		"Eight8_8",                       // 8 chars
		"Thirty_30_Thirty_30_Thirty_30_", // 30 chars
	}

	for _, v := range validPrefixes {
		_, errors := validateAdminPassword(v, "admin_password")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Admin Password: %q", v, errors)
		}
	}

	invalidPrefixes := []string{
		"",
		"Seven_7",                         // 7 chars
		"ThirtyOne_31_ThirtyOne_31_Thirt", // 31 chars
		"NoSpecial1",
		"nouppercase_1",
		"NOLOWERCASE_1",
		"No_Number",
		"Has Whitespace_1",
		"Containssys_1",
		"Containssystem_1",
		"Containsroot_1",
		"Containsoracle_1",
		"Containsdbsnmp_1",
	}

	for _, v := range invalidPrefixes {
		_, errors := validateAdminPassword(v, "admin_password")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid Admin Password", v)
		}
	}
}
