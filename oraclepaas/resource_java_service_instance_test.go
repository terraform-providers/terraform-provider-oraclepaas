package oraclepaas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOraclePAASJavaServiceInstance_Basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccJavaServiceInstanceBasic(ri)
	resourceName := "oraclepaas_java_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJavaServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJavaServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "level", "PAAS"),
					resource.TestCheckResourceAttr(
						resourceName, "weblogic_server.0.admin.0.username", "terraform-user"),
					resource.TestCheckResourceAttrSet(
						resourceName, "weblogic_server.0.admin.0.hostname"),
					resource.TestCheckResourceAttrSet(
						resourceName, "weblogic_server.0.root_url"),
				),
			},
		},
	})
}

// TODO Add support for OTD in the provider
/*
func TestAccOraclePAASJavaServiceInstance_OTD(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccJavaServiceInstanceOTD(ri)
	resourceName := "oraclepaas_java_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJavaServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJavaServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "level", "PAAS"),
					resource.TestCheckResourceAttr(
						resourceName, "otd.0.shape", "oc1m"),
				),
			},
		},
	})
} */

func TestAccOraclePAASJavaServiceInstance_Clusters(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccJavaServiceInstanceClusters(ri)
	resourceName := "oraclepaas_java_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJavaServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJavaServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "weblogic_server.0.cluster.#", "2"),
				),
			},
		},
	})
}

func TestAccOraclePAASJavaServiceInstance_UpdateShape(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccJavaServiceInstanceBasic(ri)
	config2 := testAccJavaServiceInstanceUpdateShape(ri)
	resourceName := "oraclepaas_java_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJavaServiceInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJavaServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "weblogic_server.0.shape", "oc3"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJavaServiceInstanceExists,
					resource.TestCheckResourceAttr(
						resourceName, "weblogic_server.0.shape", "oc5"),
				),
			},
		},
	})
}

func testAccCheckJavaServiceInstanceExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPAASClient).javaClient.ServiceInstanceClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_java_service_instance" {
			continue
		}

		input := java.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetServiceInstance(&input); err != nil {
			return fmt.Errorf("Error retrieving state of JavaServiceInstance %s: %+v", input.Name, err)
		}
	}

	return nil
}

func testAccCheckJavaServiceInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OPAASClient).javaClient.ServiceInstanceClient()
	if client == nil {
		return fmt.Errorf("Java Client is not initialized. Make sure to use `java_endpoint` variable or `ORACLEPAAS_JAVA_ENDPOINT` env variable")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_java_service_instance" {
			continue
		}

		input := java.GetServiceInstanceInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetServiceInstance(&input); err == nil {
			return fmt.Errorf("JavaServiceInstance %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccJavaServiceInstanceBasic(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_database_service_instance" "test" {
   bring_your_own_license = true

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
}

resource "oraclepaas_java_service_instance" "test" {
    bring_your_own_license = true

    name = "tfinstance%d"
    edition = "SUITE"
    service_version = "12cRelease212"
    ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
    force_delete = true	
    bring_your_own_license = true

    weblogic_server {
        shape = "oc3"
        database {
            name = "${oraclepaas_database_service_instance.test.name}"
			username = "sys"
            password = "Test_String7"
        }
        admin {
            username = "terraform-user"
            password = "Test_String7"
        }
    }
    backups {
        cloud_storage_container = "%sacctest-%d"
        auto_generate = true
    }
}`, rInt, os.Getenv("OPC_STORAGE_URL"), rInt, rInt, os.Getenv("OPC_STORAGE_URL"), rInt)
}

func testAccJavaServiceInstanceUpdateShape(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_database_service_instance" "test" {
    bring_your_own_license = true

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
    }
  
    backups {
        cloud_storage_container = "%sacctest-%d"
        create_if_missing = true
    }
}

resource "oraclepaas_java_service_instance" "test" {
    bring_your_own_license = true

    name = "tfinstance%d"
    edition = "SUITE"
    service_version = "12cRelease212"
    ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
    force_delete = true	
    weblogic_server {
        shape = "oc5"
        database {
            name = "${oraclepaas_database_service_instance.test.name}"
            username = "sys"
            password = "Test_String7"
        }
        admin {
            username = "terraform-user"
            password = "Test_String7"
        }
    }
    backups {
        cloud_storage_container = "%sacctest-%d"
        auto_generate = true
    }
}`, rInt, os.Getenv("OPC_STORAGE_URL"), rInt, rInt, os.Getenv("OPC_STORAGE_URL"), rInt)
}

func testAccJavaServiceInstanceOTD(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_database_service_instance" "test" {
    bring_your_own_license = true

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
}

resource "oraclepaas_java_service_instance" "test" {
    bring_your_own_license = true

    name = "tfinstance%d"
    edition = "SUITE"
    service_version = "12cRelease212"
    ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
    force_delete = true
    bring_your_own_license = true

    weblogic_server {
        shape = "oc3"
        database {
            name = "${oraclepaas_database_service_instance.test.name}"
            username = "sys"
            password = "Test_String7"
        }
        admin {
            username = "terraform-user"
            password = "Test_String7"
        }
    }

    backups {
        cloud_storage_container = "%sacctest-%d"
        auto_generate = true
    }

    otd {
        admin {
            username = "terraform-user"
            password = "Test_String7"
        }
        shape = "oc1m"
    }
}`, rInt, os.Getenv("OPC_STORAGE_URL"), rInt, rInt, os.Getenv("OPC_STORAGE_URL"), rInt)
}

func testAccJavaServiceInstanceClusters(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_database_service_instance" "test" {
    bring_your_own_license = true

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
}

resource "oraclepaas_java_service_instance" "test" {
    bring_your_own_license = true

    name = "tfinstance%d"
    edition = "SUITE"
    service_version = "12cRelease212"
    ssh_public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
    force_delete = true	
    bring_your_own_license = true

    weblogic_server {
        shape = "oc3"
        database {
            name = "${oraclepaas_database_service_instance.test.name}"
            username = "sys"
            password = "Test_String7"
        }
        admin {
            username = "terraform-user"
            password = "Test_String7"
        }
        cluster {
            name = "matthew-test"
            type = "APPLICATION_CLUSTER"
        }
        cluster {
            name = "matthew-test2"
            type = "CACHING_CLUSTER"
            shape = "oc3"
        }
    }
	
    backups {
        cloud_storage_container = "%sacctest-%d"
        auto_generate = true
    }
}`, rInt, os.Getenv("OPC_STORAGE_URL"), rInt, rInt, os.Getenv("OPC_STORAGE_URL"), rInt)
}
