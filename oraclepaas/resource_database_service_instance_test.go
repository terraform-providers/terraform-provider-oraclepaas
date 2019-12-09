package oraclepaas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func TestAccOPAASDatabaseServiceInstance_UpdateVolumes(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDatabaseServiceInstanceBasic(ri)
	config2 := testAccDatabaseServiceInstanceUpdateVolumes(ri)
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
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "database_configuration.0.backup_storage_volume_size", "10"),
					resource.TestCheckResourceAttr(
						resourceName, "database_configuration.0.data_storage_volume_size", "10"),
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
    #bring_your_own_license = true

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
    #bring_your_own_license = true

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
    #bring_your_own_license = true

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
    #bring_your_own_license = true

    database_configuration {
        admin_password = "Test_String7"
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
    #bring_your_own_license = true

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
    #bring_your_own_license = true

    database_configuration {
        admin_password = "Test_String7"
        backup_destination = "NONE"
        sid = "ORCL"
        usable_storage = 15
    }

    desired_state = "stop"
}`, rInt)
}

func testAccDatabaseServiceInstanceUpdateVolumes(rInt int) string {
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
    #bring_your_own_license = true

    database_configuration {
        admin_password = "Test_String7"
        backup_destination = "NONE"
        sid = "ORCL"
		usable_storage = 15
		backup_storage_volume_size = 10
		data_storage_volume_size = 10
    }
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
