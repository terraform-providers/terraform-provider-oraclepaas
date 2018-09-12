package oraclepaas

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// The Exadata Service resource is a variant of the Database Cloud Service.

func resourceOraclePAASExadataServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPAASExadataServiceInstanceCreate,
		Read:   resourceOPAASExadataServiceInstanceRead,
		Delete: resourceOPAASExadataServiceInstanceDelete,
		// Update: resourceOPAASExadataServiceInstanceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(120 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"edition": {
				// for Exadata only `EE_EP` is supported.
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(database.ServiceInstanceEnterpriseEditionExtremePerformance),
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceEnterpriseEditionExtremePerformance),
				}, true),
			},
			"level": {
				// for Exadata only `PAAS_EXADATA` is supported
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(database.ServiceInstanceLevelEXADATA),
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceLevelEXADATA),
				}, true),
			},
			"exadata_system_name": {
				// Exadata only attribute
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_name": {
				// Exadata only attribute
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_list": {
				// Exadata only attribute
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subscription_type": {
				// for Exadata only `MONTHLY` is supported
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(database.ServiceInstanceSubscriptionTypeMonthly),
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceSubscriptionTypeMonthly),
				}, true),
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_configuration": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_password": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Sensitive:    true,
							ValidateFunc: validateAdminPassword,
						},
						"backup_destination": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(database.ServiceInstanceBackupDestinationNone),
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceBackupDestinationBoth),
								string(database.ServiceInstanceBackupDestinationOSS),
								string(database.ServiceInstanceBackupDestinationNone),
							}, true),
						},
						"character_set": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
							ConflictsWith: []string{"instantiate_from_backup"},
						},
						"failover_database": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
							Default:  false,
						},
						"golden_gate": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"is_rac": {
							// default to true for Exadata
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"national_character_set": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: []string{"instantiate_from_backup"},
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceNCharSetUTF16),
								string(database.ServiceInstanceNCharSetUTF8),
							}, true),
						},
						"oracle_home": {
							// Exadata only attribute
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validateOracleHomeName,
						},
						"pdb_name": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
							ConflictsWith: []string{"instantiate_from_backup"},
							ValidateFunc:  validatePDBName,
						},
						"sid": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "ORCL",
							ValidateFunc: validateSID,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  database.ServiceInstanceTypeDB,
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceTypeDB),
							}, true),
						},
						"snapshot_name": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"instantiate_from_backup"},
						},
						"source_service_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"connect_descriptor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connect_descriptor_with_public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"oracle_home_name": {
							// Exadata only attribute
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"instantiate_from_backup": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				ConflictsWith: []string{
					"database_configuration.0.national_character_set",
					"database_configuration.0.character_set",
					"database_configuration.0.pdb_name",
					"database_configuration.0.snapshot_name",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_storage_container": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"cloud_storage_password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
							Computed:  true,
						},
						"cloud_storage_username": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Computed:  true,
							Sensitive: true,
						},
						"database_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"decryption_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"on_premise": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"service_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"wallet_file_content": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"backups": {
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
			"standby": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"exadata_system_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"node_list": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"notification_email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bring_your_own_license": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"em_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"glassfish_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"networking_info": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"computes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"admin_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"virtual_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"scan_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceOPAASExadataServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()

	// Database and Exadata Common attributes
	input := database.CreateServiceInstanceInput{
		Name:             d.Get("name").(string),
		Edition:          database.ServiceInstanceEdition(d.Get("edition").(string)),
		IsBYOL:           d.Get("bring_your_own_license").(bool),
		Level:            database.ServiceInstanceLevel(d.Get("level").(string)),
		SubscriptionType: database.ServiceInstanceSubscriptionType(d.Get("subscription_type").(string)),
		Version:          database.ServiceInstanceVersion(d.Get("version").(string)),
	}

	// Exadata only attributes

	if v, ok := d.GetOk("exadata_system_name"); ok {
		input.ExadataSystemName = v.(string)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		input.ClusterName = v.(string)
	}
	if _, ok := d.GetOk("node_list"); ok {
		// convert list to comma separated string
		l := d.Get("node_list").([]interface{})
		nodes := make([]string, len(l))
		for i, v := range l {
			nodes[i] = v.(string)
		}
		input.NodeList = strings.Join(nodes[:], ",")
	}

	// Common attributes

	if _, ok := d.GetOk("standby"); ok {
		if input.Parameter.FailoverDatabase != true || input.Parameter.DisasterRecovery != true {
			return fmt.Errorf("Error creating Database Service Instance: `failover_database` and `disaster_recovery` must be set to true inside the `database_configuration` block to use `standby`")
		}
		input.Standbys = expandStandby(d)
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("notification_email"); ok {
		input.EnableNotification = true
		input.NotificationEmail = v.(string)
	}

	input.Parameter, err = expandExadataParameter(d)

	info, err := client.CreateServiceInstance(&input)
	if err != nil {
		return fmt.Errorf("Error creating Exadata Service Instance: %+v", err)
	}

	d.SetId(info.Name)
	return resourceOPAASExadataServiceInstanceUpdate(d, meta)
}

func resourceOPAASExadataServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {

	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()

	getInput := database.GetServiceInstanceInput{
		Name: d.Id(),
	}

	result, err := client.GetServiceInstance(&getInput)
	if err != nil {
		// ExadataServiceInstance does not exist
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Exadata Service Instance %s: %+v", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of Exadata Service Instance %s: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("description", result.Description)
	d.Set("edition", result.Edition)
	d.Set("em_url", result.EMURL)
	d.Set("failover_database", result.FailoverDatabase)
	d.Set("glassfish_url", result.GlassFishURL)
	d.Set("identity_domain", result.IdentityDomain)
	d.Set("bring_your_own_license", result.IsBYOL)
	d.Set("level", result.Level)
	d.Set("uri", result.URI)
	d.Set("status", result.Status)
	d.Set("subscription_type", result.SubscriptionType)
	d.Set("timezone", result.Timezone)
	d.Set("version", result.Version)
	// d.Set("networking_info", result.???) // TODO
	// d.Set("snapshot_service" result.SnapshotService) // TODO add to SDK
	// d.Set("is_clone", result.IsClone) // TODO add to SDK
	// d.Set("is_managed", result.IsManaged) // TODO add to SDK
	// d.Set("use_high_performance_storage", result.UseHighPerformanceStorage)

	// Exadata attributes
	d.Set("cluster_name", result.ClusterName)
	d.Set("exadata_system_name", result.SubscriptionName)
	// d.Set("node_list", d.Get("node_list")) // not returned from API

	dbConfig, err := flattenExadataConfig(d, result)
	if err != nil {
		return err
	}
	if err := d.Set("database_configuration", dbConfig); err != nil {
		return fmt.Errorf("Error setting Exadata Service Instance database configuration: %+v", err)
	}

	if _, ok := d.GetOk("backups"); ok {
		backups, err := flattenBackupConfig(d, result)
		if err != nil {
			return err
		}
		if err := d.Set("backups", backups); err != nil {
			return fmt.Errorf("Error setting Exadata Service Instance backup configuration: %+v", err)
		}
	}

	// TODO standby configuration
	// standbys, err := flattenStandbyConfig(d, result)
	// if err != nil {
	// 	return err
	// }
	// if err := d.Set("standbys", standbys); err != nil {
	// 	return fmt.Errorf("Error setting Exadata Service Instance standby configuration: %+v", err)
	// }

	networking, err := flattenExadataNetworkingInfo(d, &result.NetworkingInfo)
	if err != nil {
		return err
	}
	if err := d.Set("networking_info", networking); err != nil {
		return fmt.Errorf("Error setting Exadata Service Instance Networking Info configuration: %+v", err)
	}

	return nil
}

func resourceOPAASExadataServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()
	name := d.Id()

	client.Timeout = d.Timeout(schema.TimeoutDelete)

	input := database.DeleteServiceInstanceInput{
		Name: name,
	}
	if err := client.DeleteServiceInstance(&input); err != nil {
		return fmt.Errorf("Error deleting Exadata Service Instance: %+v", err)
	}
	return nil
}

func resourceOPAASExadataServiceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	// dbClient, err := getDatabaseClient(meta)
	// if err != nil {
	// 	return err
	// }
	// client := dbClient.ServiceInstanceClient()

	return resourceOPAASExadataServiceInstanceRead(d, meta)
}

func expandExadataParameter(d *schema.ResourceData) (database.ParameterInput, error) {
	databaseConfigInfo := d.Get("database_configuration").([]interface{})
	attrs := databaseConfigInfo[0].(map[string]interface{})

	// Database and Exadata Cloud Service common attributes
	parameter := database.ParameterInput{
		AdminPassword:     attrs["admin_password"].(string),
		BackupDestination: database.ServiceInstanceBackupDestination(attrs["backup_destination"].(string)),
		FailoverDatabase:  attrs["failover_database"].(bool),
		GoldenGate:        attrs["golden_gate"].(bool),
		IsRAC:             attrs["is_rac"].(bool),
		SID:               attrs["sid"].(string),
		Type:              database.ServiceInstanceType(attrs["type"].(string)),
	}

	// Exadata Cloud Service only attributes
	if val, ok := attrs["oracle_home"].(string); ok && val != "" {
		parameter.OracleHomeName = val
	}

	// Common attributes
	if val, ok := attrs["character_set"].(string); ok && val != "" {
		parameter.CharSet = val
	}
	if val, ok := attrs["national_character_set"].(string); ok && val != "" {
		parameter.NCharSet = database.ServiceInstanceNCharSet(val)
	}
	if val, ok := attrs["pdb_name"].(string); ok && val != "" {
		parameter.PDBName = val
	}
	if val, ok := attrs["snapshot_name"].(string); ok && val != "" {
		parameter.SnapshotName = val
	}
	if val, ok := attrs["source_service_name"].(string); ok && val != "" {
		parameter.SourceServiceName = val
	}

	expandIbkup(d, &parameter)
	err := expandBackups(d, &parameter)
	if err != nil {
		return parameter, err
	}

	return parameter, nil
}

func flattenExadataConfig(d *schema.ResourceData, result *database.ServiceInstance) ([]interface{}, error) {
	dbConfig := make(map[string]interface{})
	dbConfig["backup_destination"] = result.BackupDestination
	dbConfig["connect_descriptor"] = result.ConnectDescriptor
	dbConfig["connect_descriptor_with_public_ip"] = result.ConnectDescriptorWithPublicIP
	dbConfig["pdb_name"] = result.PDBName
	dbConfig["sid"] = result.SID
	dbConfig["is_rac"] = result.RACDatabase
	dbConfig["character_set"] = result.CharSet
	dbConfig["national_character_set"] = result.NCharSet
	dbConfig["listener_port"] = result.ListenerPort
	dbConfig["oracle_home_name"] = result.OracleHomeName

	// attributes not returned from API
	dbConfig["oracle_home"] = d.Get("database_configuration.0.oracle_home")
	dbConfig["type"] = d.Get("database_configuration.0.type")
	dbConfig["admin_password"] = d.Get("database_configuration.0.admin_password")

	return []interface{}{dbConfig}, nil
}

func flattenBackupConfig(d *schema.ResourceData, result *database.ServiceInstance) ([]interface{}, error) {
	backupsConfig := make(map[string]interface{})
	backupsConfig["cloud_storage_container"] = result.CloudStorageContainer

	return []interface{}{backupsConfig}, nil
}

func flattenExadataNetworkingInfo(d *schema.ResourceData, result *database.NetworkingInfo) ([]interface{}, error) {
	networkConfig := make(map[string]interface{})
	networkConfig["admin_network"] = result.AdminNetwork
	networkConfig["backup_network"] = result.BackupNetwork
	networkConfig["client_network"] = result.ClientNetwork
	networkConfig["scan_ips"] = result.ScanIPs

	computes, err := flattenExadataComputes(d, &result.Computes)
	if err != nil {
		return nil, err
	}
	networkConfig["computes"] = computes

	return []interface{}{networkConfig}, nil
}

func flattenExadataComputes(d *schema.ResourceData, result *[]database.ComputesInfo) ([]interface{}, error) {
	flattenedComputes := make([]interface{}, 0)

	for _, info := range *result {
		compute := make(map[string]interface{})
		compute["admin_ip"] = info.AdminIP
		compute["client_ip"] = info.ClientIP
		compute["hostname"] = info.Hostname
		compute["virtual_ip"] = info.VirtualIP
		flattenedComputes = append(flattenedComputes, compute)
	}
	return flattenedComputes, nil
}

// Validate the Oracle Home name.
// up to 64 characters; must start with a letter and can contain only letters, numbers and underscores (_); can not end with an underscore (_).
func validateOracleHomeName(v interface{}, k string) (ws []string, errors []error) {
	if len(v.(string)) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q must not exceed 64 characters", k))
	}
	if match, _ := regexp.MatchString("^([a-zA-Z])((([a-zA-Z0-9_]*)([a-zA-Z0-9]+))?)$", v.(string)); match != true {
		errors = append(errors, fmt.Errorf(
			"%q must start with a letter and can contain only letters, numbers and underscores (_); can not end with an underscore (_)", k))
	}
	return
}
