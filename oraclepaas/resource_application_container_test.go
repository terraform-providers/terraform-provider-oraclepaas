package oraclepaas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/application"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOraclePAASApplicationContainer_Basic(t *testing.T) {
	ri := acctest.RandIntRange(1, 10000)
	config := testAccApplicationContainerBasic(ri)
	resourceName := "oraclepaas_application_container.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationContainerExists,
					resource.TestCheckResourceAttrSet(
						resourceName, "web_url"),
					resource.TestCheckResourceAttrSet(
						resourceName, "app_url"),
				),
			},
		},
	})
}

func TestAccOraclePAASApplicationContainer_Manifest(t *testing.T) {
	ri := acctest.RandIntRange(1, 10000)
	config := testAccApplicationContainerManifest(ri)
	resourceName := "oraclepaas_application_container.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationContainerExists,
					resource.TestCheckResourceAttrSet(
						resourceName, "web_url"),
					resource.TestCheckResourceAttrSet(
						resourceName, "app_url"),
				),
			},
		},
	})
}

func TestAccOraclePAASApplicationContainer_ManifestAttr(t *testing.T) {
	ri := acctest.RandIntRange(1, 10000)
	config := testAccApplicationContainerManifestAttr(ri)
	resourceName := "oraclepaas_application_container.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationContainerExists,
					resource.TestCheckResourceAttrSet(
						resourceName, "web_url"),
				),
			},
		},
	})
}

func TestAccOraclePAASApplicationContainer_Deployment(t *testing.T) {
	ri := acctest.RandIntRange(1, 10000)
	config := testAccApplicationContainerDeployment(ri)
	resourceName := "oraclepaas_application_container.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationContainerExists,
					resource.TestCheckResourceAttrSet(
						resourceName, "web_url"),
					resource.TestCheckResourceAttrSet(
						resourceName, "app_url"),
				),
			},
		},
	})
}

func TestAccOraclePAASApplicationContainer_DeploymentAttr(t *testing.T) {
	ri := acctest.RandIntRange(1, 10000)
	config := testAccApplicationContainerDeploymentAttr(ri)
	resourceName := "oraclepaas_application_container.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationContainerExists,
					resource.TestCheckResourceAttrSet(
						resourceName, "web_url"),
					resource.TestCheckResourceAttrSet(
						resourceName, "app_url"),
				),
			},
		},
	})
}

func TestAccOraclePAASApplicationContainer_ManifestDeploymentAttr(t *testing.T) {
	ri := acctest.RandIntRange(1, 10000)
	config := testAccApplicationContainerManifestDeploymentAttr(ri)
	resourceName := "oraclepaas_application_container.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationContainerExists,
					resource.TestCheckResourceAttrSet(
						resourceName, "app_url"),
				),
			},
		},
	})
}

func testAccCheckApplicationContainerExists(s *terraform.State) error {
	aClient, err := getApplicationClient(testAccProvider.Meta())
	if err != nil {
		return err
	}
	client := aClient.ContainerClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_application_container" {
			continue
		}

		input := application.GetApplicationContainerInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetApplicationContainer(&input); err != nil {
			return fmt.Errorf("Error retrieving state of Application Container %s: %+v", input.Name, err)
		}
	}

	return nil
}

func testAccCheckApplicationContainerDestroy(s *terraform.State) error {
	aClient, err := getApplicationClient(testAccProvider.Meta())
	if err != nil {
		return err
	}
	client := aClient.ContainerClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oraclepaas_application_container" {
			continue
		}

		input := application.GetApplicationContainerInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetApplicationContainer(&input); err == nil {
			return fmt.Errorf("Application Container %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}

func testAccApplicationContainerBasic(rInt int) string {
	return fmt.Sprintf(`resource "oraclepaas_application_container" "test" {
    name        = "testappcontainer%d"
  }`, rInt)
}

func testAccApplicationContainerManifest(rInt int) string {
	return fmt.Sprintf(`resource "oraclepaas_application_container" "test" {
    name        = "testappcontainer%d"
	manifest_file = "testdata/manifest.json"
  }`, rInt)
}

func testAccApplicationContainerManifestAttr(rInt int) string {
	return fmt.Sprintf(`resource "oraclepaas_application_container" "test" {
    name        = "testappcontainer%d"
	manifest_attributes {
		runtime {
			major_version = 7
		}
		command = "sh target/bin/start"
		release {
			build = 150520.1154
			commit = "d8c2596364d9584050461"
			version = "15.1.0"
		}
		notes = "notes related to release"
		mode =  "rolling"
	}
  }`, rInt)
}

func testAccApplicationContainerDeployment(rInt int) string {
	return fmt.Sprintf(`resource "oraclepaas_application_container" "test" {
    name        = "testappcontainer%d"
	deployment_file = "testdata/deployment.json"
  }`, rInt)
}

func testAccApplicationContainerDeploymentAttr(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_application_container" "test" {
	name        = "testappcontainer%d"
	deployment_attributes {
		memory = "2G"
    	instances = 1
    	environment {
        	"NO_OF_CONNECTIONS" = "25"
        	"TWITTER_ID" = "JAVA"
    	}
	}
}`, rInt)
}

func testAccApplicationContainerManifestDeploymentAttr(rInt int) string {
	return fmt.Sprintf(`
resource "oraclepaas_application_container" "test" {
	name        = "testappcontainer%d"
	manifest_attributes {
		runtime {
			major_version = 7
		}
		command = "sh target/bin/start"
		release {
			build = 150520.1154
			commit = "d8c2596364d9584050461"
			version = "15.1.0"
		}
		notes = "notes related to release"
		mode =  "rolling"
	}
	deployment_attributes {
		memory = "2G"
    	instances = 1
    	environment {
        	"NO_OF_CONNECTIONS" = "25"
        	"TWITTER_ID" = "JAVA"
    	}
	}
}`, rInt)
}
