package oraclepaas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOPAASDataSourceDatabaseServiceInstance_Basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDataSourceDatabaseServiceInstanceBasic(ri)
	resourceName := "data.oraclepaas_database_service_instance.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func testAccDataSourceDatabaseServiceInstanceBasic(rInt int) string {
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
}

data "oraclepaas_database_service_instance" "test" {
	name = "${oraclepaas_database_service_instance.test.name}"
}`, rInt)
}
