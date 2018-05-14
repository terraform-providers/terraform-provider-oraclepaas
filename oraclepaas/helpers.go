package oraclepaas

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-oracle-terraform/application"
	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/terraform/helper/schema"
)

func javaServiceInstanceShapes() []string {
	return []string{string(java.ServiceInstanceShapeOC3), string(java.ServiceInstanceShapeOC4), string(java.ServiceInstanceShapeOC5),
		string(java.ServiceInstanceShapeOC6), string(java.ServiceInstanceShapeOC7), string(java.ServiceInstanceShapeOC1M), string(java.ServiceInstanceShapeOC2M),
		string(java.ServiceInstanceShapeOC3M), string(java.ServiceInstanceShapeOC4M), string(java.ServiceInstanceShapeOC5M),
		string(java.ServiceInstanceShapeVMStandard1_1), string(java.ServiceInstanceShapeVMStandard1_2), string(java.ServiceInstanceShapeVMStandard1_4),
		string(java.ServiceInstanceShapeVMStandard1_8), string(java.ServiceInstanceShapeVMStandard1_16), string(java.ServiceInstanceShapeVMStandard2_1),
		string(java.ServiceInstanceShapeVMStandard2_2), string(java.ServiceInstanceShapeVMStandard2_2), string(java.ServiceInstanceShapeVMStandard2_4),
		string(java.ServiceInstanceShapeVMStandard2_8), string(java.ServiceInstanceShapeVMStandard2_16), string(java.ServiceInstanceShapeVMStandard2_24),
		string(java.ServiceInstanceShapeBMStandard1_36), string(java.ServiceInstanceShapeBMStandard2_52)}
}

// Helper function to get a string list from the schema, and alpha-sort it
func getStringList(d *schema.ResourceData, key string) []string {
	if _, ok := d.GetOk(key); !ok {
		return nil
	}
	l := d.Get(key).([]interface{})
	res := make([]string, len(l))
	for i, v := range l {
		res[i] = v.(string)
	}
	sort.Strings(res)
	return res
}

// Helper function to set a string list in the schema, in an alpha-sorted order.
func setStringList(d *schema.ResourceData, key string, value []string) error {
	sort.Strings(value)
	return d.Set(key, value)
}

// Helper function to get an int list from the schema, and numerically sort it
func getIntList(d *schema.ResourceData, key string) []int {
	if _, ok := d.GetOk(key); !ok {
		return nil
	}

	l := d.Get(key).([]interface{})
	res := make([]int, len(l))
	for i, v := range l {
		res[i] = v.(int)
	}
	sort.Ints(res)
	return res
}

func setIntList(d *schema.ResourceData, key string, value []int) error {
	sort.Ints(value)
	return d.Set(key, value)
}

// A user may inadvertently call the database service without passing in the required parameters (because it's optional)
// so we check to make sure that the database client has been initialized
func getDatabaseClient(meta interface{}) (*database.DatabaseClient, error) {
	client := meta.(*OPAASClient).databaseClient
	if client == nil {
		return nil, fmt.Errorf("Database Client is not initialized. Make sure to use `database_endpoint` variable or `ORACLEPAAS_DATABASE_ENDPOINT` env variable")
	}
	return client, nil
}

// A user may inadvertently call the java without passing in the required parameters to use that service
// (because it's optional) so we check to make sure that the java client has been initialized
func getJavaClient(meta interface{}) (*java.JavaClient, error) {
	client := meta.(*OPAASClient).javaClient
	if client == nil {
		return nil, fmt.Errorf("Java Client is not initialized. Make sure to use `java_endpoint` variable or `ORACLEPAAS_JAVA_ENDPOINT` env variable")
	}
	return client, nil
}

// A user may inadvertently call the application cloud without passing in the required parameters to use that service
// (because it's optional) so we check to make sure that the application client has been initialized
func getApplicationClient(meta interface{}) (*application.Client, error) {
	client := meta.(*OPAASClient).applicationClient
	if client == nil {
		return nil, fmt.Errorf("Application Client is not initialized. Make sure to use `application_endpoint` variable or `ORACLEPAAS_APPLICAITON_ENDPOINT` env variable")
	}
	return client, nil
}
