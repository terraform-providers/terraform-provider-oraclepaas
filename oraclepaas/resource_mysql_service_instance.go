package oraclepaas

import (
	"fmt"
	"log"
	"strconv"
	"time"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/mysql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOraclePAASMySQLServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOraclePAASMySQLServiceInstanceCreate,
		Read:   resourceOraclePAASMySQLServiceInstanceRead,
		Update: resourceOraclePAASMySQLServiceInstanceUpdate,
		Delete: resourceOraclePAASMySQLServiceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(120 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"service_description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"vm_public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"availability_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"backup_destination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(mysql.ServiceInstanceBackupDestinationBoth),
					string(mysql.ServiceInstanceBackupDestinationOSS),
					string(mysql.ServiceInstanceBackupDestinationNone),
				}, true),
			},

			"metering_frequency": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(mysql.ServiceInstanceMeteringFrequencyHourly),
					string(mysql.ServiceInstanceMeteringFrequencyMonthly),
				}, true),
			},

			"enable_notification": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"notification_email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Use for OCI configuration (not OCI-Classic)
			"ip_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vm_user": {
				// default to opc
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"cloud_storage_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_storage_container": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"cloud_storage_username": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Computed:  true,
							Sensitive: true,
						},
						"cloud_storage_password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Computed:  true,
							Sensitive: true,
						},
						"create_if_missing": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
					},
				},
			},

			"mysql_configuration": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"db_storage": {
							// integer. default 25
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(25, 1024),
						},
						"mysql_charset": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"mysql_collation": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"mysql_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(3200, 3399),
							Default:      3306,
						},
						"mysql_username": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"mysql_password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"shape": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						/* Couldn't get these to work with the current API. I've commented them out for now
						"mysql_options" : {
							Type: schema.TypeString,
							Optional: true,
							ForceNew: true,
						}
						"mysql_timezone" : {
							Type: schema.TypeString,
							Optional: true,
							ForceNew: true,
						}
						*/
						"snapshot_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"source_service_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"enterprise_monitor": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"enterprise_monitor_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"em_agent_password": {
										Type:      schema.TypeString,
										Optional:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"em_agent_username": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringLenBetween(2, 32),
									},
									"em_password": {
										Type:      schema.TypeString,
										Optional:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"em_username": {
										Type:      schema.TypeString,
										Optional:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"em_port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ForceNew:     true,
										Computed:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
								},
							},
						},
						// this comes from the service.
						"connect_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"service_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_release_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"em_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}, // end declaration
	} // end return
}

func resourceOraclePAASMySQLServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] Resource state: %#v", d.State())
	log.Print("[DEBUG] Creating mySQL service instance")

	mySQLClient, err := getMySQLClient(meta)
	if err != nil {
		return err
	}
	client := mySQLClient.ServiceInstanceClient()

	input := mysql.CreateServiceInstanceInput{}
	input.ServiceParameters, err = getServiceParameters(d)
	if err != nil {
		log.Printf("[Error] : Error while extracting MySQL Service Instance information : %s", err)
		return err
	}

	input.ComponentParameters, err = getComponentParameters(d)
	if err != nil {
		log.Printf("[Error] : Error while extracting MySQL component information from TF file. : %s", err)
		return err
	}

	log.Printf("[DEBUG] : Testing Create :%v : ", &input)

	_, err = client.CreateServiceInstance(&input)

	if err != nil {
		log.Printf("[Error] : Error while creating MySQL Service Instance : %v", err)
		return err
	}

	d.SetId(input.ServiceParameters.ServiceName)
	return resourceOraclePAASMySQLServiceInstanceUpdate(d, meta)
}

/**
getServiceParameters gets the values from the terraform resource file, and updates the inputParameter
with the respective values for calling the "Create"
*/
func getServiceParameters(d *schema.ResourceData) (mysql.ServiceParameters, error) {

	input := &mysql.ServiceParameters{
		ServiceName:       d.Get("service_name").(string),
		BackupDestination: d.Get("backup_destination").(string),
		VMPublicKeyText:   d.Get("vm_public_key").(string),
	}

	if value, ok := d.GetOk("metering_frequency"); ok {
		input.MeteringFrequency = value.(string)
	}

	if value, ok := d.GetOk("region"); ok {
		input.Region = value.(string)
	}

	if value, ok := d.GetOk("service_description"); ok {
		input.ServiceDescription = value.(string)
	}

	err := expandCloudStorage(d, input)
	if err != nil {
		return *input, err
	}

	return *input, nil
}

/**
Expands and reads the values in the Cloud_Storage list specified in the terraform file.
*/
func expandCloudStorage(d *schema.ResourceData, parameter *mysql.ServiceParameters) error {

	cloudStorageInfo := d.Get("cloud_storage_configuration").([]interface{})

	if parameter.BackupDestination == string(mysql.ServiceInstanceBackupDestinationBoth) || parameter.BackupDestination == string(mysql.ServiceInstanceBackupDestinationOSS) {
		if len(cloudStorageInfo) == 0 {
			return fmt.Errorf("`cloud_storage_configuration` must be set if `backup_destination` is set to `OSS` or `BOTH`")
		}
	}

	if len(cloudStorageInfo) > 0 {
		attrs := cloudStorageInfo[0].(map[string]interface{})
		parameter.CloudStorageContainer = attrs["cloud_storage_container"].(string)
		parameter.CloudStorageContainerAutoGenerate = attrs["create_if_missing"].(bool)
		if val, ok := attrs["cloud_storage_username"].(string); ok && val != "" {
			parameter.CloudStorageUsername = val
		}
		if val, ok := attrs["cloud_storage_password"].(string); ok && val != "" {
			parameter.CloudStoragePassword = val
		}
	}
	return nil
}

func expandEM(input map[string]interface{}, parameter *mysql.MySQLParameters) error {

	emInfo := input["enterprise_monitor_configuration"].([]interface{})

	log.Printf("[DEBUG] input                       : %v", input)
	log.Printf("[DEBUG] parameter.EnterpriseMonitor : %v", parameter.EnterpriseMonitor)
	log.Printf("[DEBUG] emInfo                      : %d", len(emInfo))

	if parameter.EnterpriseMonitor == "Yes" {
		if len(emInfo) == 0 {
			return fmt.Errorf("`enterprise_monitor_configuration` must be set if `enterprise_monitor` is set to `Yes`")
		}
	}

	if len(emInfo) > 0 {
		attrs := emInfo[0].(map[string]interface{})

		if val, ok := attrs["em_agent_password"].(string); ok && val != "" {
			parameter.EnterpriseMonitorAgentPassword = attrs["em_agent_password"].(string)
		}

		if val, ok := attrs["em_agent_user"].(string); ok && val != "" {
			parameter.EnterpriseMonitorAgentUser = attrs["em_agent_user"].(string)
		}

		if val, ok := attrs["em_password"].(string); ok && val != "" {
			parameter.EnterpriseMonitorManagerPassword = attrs["em_password"].(string)
		}

		if val, ok := attrs["em_user"].(string); ok && val != "" {
			parameter.EnterpriseMonitorManagerUser = attrs["em_user"].(string)
		}

		if val, ok := attrs["em_port"].(string); ok && val != "" {
			parameter.MysqlEMPort = attrs["em_port"].(string)
		}
	}

	return nil
}

func getComponentParameters(d *schema.ResourceData) (mysql.ComponentParameters, error) {

	result := mysql.ComponentParameters{}

	if v, ok := d.GetOk("mysql_configuration"); ok {
		mysqlList := v.(*schema.Set).List()

		// get the first entry.
		mysqlItem := mysqlList[0].(map[string]interface{})
		MysqlInput := &mysql.MySQLParameters{
			DBName:    mysqlItem["db_name"].(string),
			DBStorage: strconv.Itoa(mysqlItem["db_storage"].(int)),
		}

		log.Printf("[DEBUG] Enterprise Monitor : %v", mysqlItem["enterprise_monitor"])

		if mysqlItem["enterprise_monitor"] != nil {
			if mysqlItem["enterprise_monitor"] == true {
				MysqlInput.EnterpriseMonitor = "Yes"
			} else {
				MysqlInput.EnterpriseMonitor = "No"
			}
		}

		err := expandEM(mysqlItem, MysqlInput)
		if err != nil {
			return result, err
		}

		if mysqlItem["mysql_charset"] != nil {
			MysqlInput.MysqlCharset = mysqlItem["mysql_charset"].(string)
		}

		if mysqlItem["mysql_collation"] != nil {
			MysqlInput.MysqlCollation = mysqlItem["mysql_collation"].(string)
		}

		if mysqlItem["mysql_em_port"] != nil {
			MysqlInput.MysqlEMPort = strconv.Itoa(mysqlItem["mysql_em_port"].(int))
		}

		if mysqlItem["mysql_port"] != nil {
			MysqlInput.MysqlPort = strconv.Itoa(mysqlItem["mysql_port"].(int))
		}

		if mysqlItem["mysql_username"] != nil {
			MysqlInput.MysqlUserName = mysqlItem["mysql_username"].(string)
		}

		if mysqlItem["mysql_password"] != nil {
			MysqlInput.MysqlUserPassword = mysqlItem["mysql_password"].(string)
		}

		if mysqlItem["shape"] != nil {
			MysqlInput.Shape = mysqlItem["shape"].(string)
		}

		if mysqlItem["snapshot_name"] != nil {
			MysqlInput.SnapshotName = mysqlItem["snapshot_name"].(string)
		}

		if mysqlItem["source_service_name"] != nil {
			MysqlInput.SourceServiceName = mysqlItem["source_service_name"].(string)
		}

		result.Mysql = *MysqlInput
	}

	return result, nil
}

func resourceOraclePAASMySQLServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] Resource state: %#v", d.State())
	mysqlClient, err := getMySQLClient(meta)
	if err != nil {
		return err
	}
	client := mysqlClient.ServiceInstanceClient()

	input := mysql.GetServiceInstanceInput{
		Name: d.Id(),
	}

	result, err := client.GetServiceInstance(&input)

	if err != nil {
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading database service instance %s: %+v", d.Id(), err)
	}

	// if there is not result, there was an earlier issue. We set the ID of the mysql instance to blank.
	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of mysql service instance %s: %#v", d.Id(), result)

	d.Set("backup_destination", result.BackupDestination)
	d.Set("metering_frequency", result.MeteringFrequency)
	d.Set("service_name", result.ServiceName)
	d.Set("service_description", result.ServiceDescription)
	d.Set("service_id", result.ServiceId)
	d.Set("service_type", result.ServiceType)
	d.Set("release_version", result.ReleaseVersion)
	d.Set("service_version", result.ServiceVersion)
	d.Set("base_release_version", result.BaseReleaseVersion)
	d.Set("creator", result.Creator)
	d.Set("creation_date", result.CreationDate)
	d.Set("state", result.Status)

	if err := updateMySQLAttributesFromAttachments(d, result.Components.Mysql); err != nil {
		return err
	}

	return nil
}

func updateMySQLAttributesFromAttachments(d *schema.ResourceData, instanceInfo mysql.MysqlInfo) error {

	result := make([]map[string]interface{}, 0)

	if v, ok := d.GetOk("mysql_configuration"); ok {
		mysqlList := v.(*schema.Set).List()

		if len(mysqlList) != 1 {
			return fmt.Errorf("Invalid mySQL Instance info")
		}
		newState := mysqlList[0].(map[string]interface{})
		attributeMap := instanceInfo.Attributes

		if attr, ok := attributeMap["MYSQL_CHARACTER_SET"]; ok {
			newState["mysql_charset"] = attr.Value
		}

		if attr, ok := attributeMap["MYSQL_COLLATION"]; ok {
			newState["mysql_collation"] = attr.Value
		}

		if attr, ok := attributeMap["MYSQL_DBNAME"]; ok {
			newState["db_name"] = attr.Value
		}

		if attr, ok := attributeMap["shape"]; ok {
			newState["shape"] = attr.Value
		}

		if attr, ok := attributeMap["CONNECT_STRING"]; ok {
			newState["connect_string"] = attr.Value
		}

		if attr, ok := attributeMap["MYSQL_ENTERPRISE_MONITOR"]; ok {
			if attr.Value == "Yes" || attr.Value == "YES" {
				newState["enterprise_monitor"] = true
			} else {
				newState["enterprise_monitor"] = false
			}
		}

		/* Temporarily commented out. Base service has some issues with Timezone
		if attr, ok := attributeMap["MYSQL_TIMEZONE"]; ok {
			newState["enterprise_monitor"] = attr.Value
		}
		*/

		// Update from the VM Map
		vmInstancesMap := instanceInfo.VMInstances

		if len(vmInstancesMap) != 1 {
			return fmt.Errorf("Error. Failed to detect correct mySQL Instance information.")
		}

		for _, vmInstance := range vmInstancesMap {
			newState["hostname"] = vmInstance.HostName
			newState["ip_address"] = vmInstance.IPAddress
			newState["public_ip_address"] = vmInstance.PublicIPAddress
			newState["component_type"] = vmInstance.ComponentType
			newState["state"] = vmInstance.State
		}

		result = append(result, newState)
	}

	return d.Set("mysql_configuration", result)
}

func resourceOraclePAASMySQLServiceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceOraclePAASMySQLServiceInstanceRead(d, meta)
}

func resourceOraclePAASMySQLServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] Resource state: %#v", d.State())
	log.Print("[DEBUG] Deleting mySQL service instance")

	mySQLClient, err := getMySQLClient(meta)

	if err != nil {
		return err
	}

	client := mySQLClient.ServiceInstanceClient()
	serviceParams, err := getServiceParameters(d)

	if err != nil {
		log.Printf("Error : %s", err)
		return err
	}

	jobID := serviceParams.ServiceName

	log.Printf("[DEBUG] Deleting DatabaseServiceInstance: %v", jobID)

	if err := client.DeleteServiceInstance(jobID); err != nil {
		return fmt.Errorf("Error deleting MySQL instance %s: %s", jobID, err)
	}

	return nil
}
