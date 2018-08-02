package oraclepaas

import (
	"fmt"
	"log"
	"strconv"
	"time"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOraclePAASDatabaseServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPAASDatabaseServiceInstanceCreate,
		Read:   resourceOPAASDatabaseServiceInstanceRead,
		Delete: resourceOPAASDatabaseServiceInstanceDelete,
		Update: resourceOPAASDatabaseServiceInstanceUpdate,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceStandardEdition),
					string(database.ServiceInstanceEnterpriseEdition),
					string(database.ServiceInstanceEnterpriseEditionHighPerformance),
					string(database.ServiceInstanceEnterpriseEditionExtremePerformance),
				}, true),
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(database.ServiceInstanceLevelPAAS),
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceLevelPAAS),
					string(database.ServiceInstanceLevelEXADATA),
					string(database.ServiceInstanceLevelBasic),
				}, true),
			},
			"shape": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"exadata_system_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"shape"},
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_list": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceSubscriptionTypeHourly),
					string(database.ServiceInstanceSubscriptionTypeMonthly),
				}, true),
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssh_public_key": {
				Type:     schema.TypeString,
				Optional: true,
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
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
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
							ConflictsWith: []string{"instantiate_from_backup", "hybrid_disaster_recovery"},
						},
						"db_demo": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disaster_recovery": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
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
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"national_character_set": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: []string{"instantiate_from_backup", "hybrid_disaster_recovery"},
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceNCharSetUTF16),
								string(database.ServiceInstanceNCharSetUTF8),
							}, true),
						},
						"pdb_name": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
							ConflictsWith: []string{"instantiate_from_backup", "hybrid_disaster_recovery"},
						},
						"sid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "ORCL",
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "UTC",
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
						"usable_storage": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(15, 2048),
						},
						"snapshot_name": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"instantiate_from_backup", "hybrid_disaster_recovery"},
						},
						"source_service_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
			"hybrid_disaster_recovery": {
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
						"cloud_storage_username": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cloud_storage_password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
					},
				},
			},
			"default_access_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// All Instances share this
						"enable_ssh": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						// Single Instance rules
						"enable_http": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_http_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_db_console": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_db_express": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_db_listener": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						// RAC Rules
						"enable_em_console": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_rac_db_listener": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_scan_listener": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_rac_ons": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
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
						"availability_domain": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"availability_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_reservations": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"high_performance_storage": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"desired_state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceLifecycleStateStop),
					string(database.ServiceInstanceLifecycleStateRestart),
					string(database.ServiceInstanceLifecycleStateStart),
				}, true),
			},
			"cloud_storage_container": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"compute_site_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"dbaas_monitor_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"em_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"glassfish_url": {
				Type:     schema.TypeString,
				ForceNew: true,
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
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceOPAASDatabaseServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating database service instance")

	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()

	input := database.CreateServiceInstanceInput{
		Name:                      d.Get("name").(string),
		Edition:                   database.ServiceInstanceEdition(d.Get("edition").(string)),
		IPReservations:            getStringList(d, "ip_reservations"),
		IsBYOL:                    d.Get("bring_your_own_license").(bool),
		Level:                     database.ServiceInstanceLevel(d.Get("level").(string)),
		SubscriptionType:          database.ServiceInstanceSubscriptionType(d.Get("subscription_type").(string)),
		UseHighPerformanceStorage: d.Get("high_performance_storage").(bool),
		Version:                   database.ServiceInstanceVersion(d.Get("version").(string)),
	}

	if v, ok := d.GetOk("shape"); ok {
		input.Shape = database.ServiceInstanceShape(v.(string))
	}

	if v, ok := d.GetOk("ssh_public_key"); ok {
		input.VMPublicKey = v.(string)
	}

	if v, ok := d.GetOk("exadata_system_name"); ok {
		input.ExadataSystemName = v.(string)
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		input.ClusterName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}

	if v, ok := d.GetOk("notification_email"); ok {
		input.EnableNotification = true
		input.NotificationEmail = v.(string)
	}

	if v, ok := d.GetOk("ip_network"); ok {
		input.IPNetwork = v.(string)
	}

	if _, ok := d.GetOk("ip_reservations"); ok {
		input.IPReservations = getStringList(d, "ip_reservations")
	}

	if v, ok := d.GetOk("region"); ok {
		input.Region = v.(string)
	}

	if v, ok := d.GetOk("availability_domain"); ok {
		input.AvailabilityDomain = v.(string)
	}

	if v, ok := d.GetOk("subnet"); ok {
		input.Subnet = v.(string)
	}

	// Only the PaaS levels can have a parameter.
	if input.Level != database.ServiceInstanceLevelBasic {
		input.Parameter, err = expandParameter(d)
		if err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("standby"); ok {
		if input.Parameter.FailoverDatabase != true || input.Parameter.DisasterRecovery != true {
			return fmt.Errorf("Error creating Database Service Instance: `failover_database` and `disaster_recovery` must be set to true inside the `database_configuration` block to use `standby`")
		}
		input.Standbys = expandStandby(d)
		if err != nil {
			return err
		}
	}

	info, err := client.CreateServiceInstance(&input)
	if err != nil {
		return fmt.Errorf("Error creating DatabaseServiceInstance: %+v", err)
	}

	d.SetId(info.Name)
	return resourceOPAASDatabaseServiceInstanceUpdate(d, meta)
}

func resourceOPAASDatabaseServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()

	log.Printf("[DEBUG] Reading state of ip reservation %s", d.Id())
	getInput := database.GetServiceInstanceInput{
		Name: d.Id(),
	}

	result, err := client.GetServiceInstance(&getInput)
	if err != nil {
		// DatabaseServiceInstance does not exist
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading database service instance %s: %+v", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of database service instance %s: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("region", result.Region)
	d.Set("availability_domain", result.AvailabilityDomain)
	d.Set("description", result.Description)
	d.Set("backup_destination", result.BackupDestination)
	d.Set("character_set", result.CharSet)
	d.Set("cloud_storage_container", result.CloudStorageContainer)
	d.Set("compute_site_name", result.ComputeSiteName)
	d.Set("connect_descriptor", result.ConnectDescriptor)
	d.Set("desired_state", d.Get("desired_state"))
	d.Set("dbaas_monitor_url", result.DBAASMonitorURL)
	d.Set("edition", result.Edition)
	d.Set("em_url", result.EMURL)
	d.Set("failover_database", result.FailoverDatabase)
	d.Set("high_performance_storage", result.UseHighPerformanceStorage)
	d.Set("glassfish_url", result.GlassFishURL)
	d.Set("identity_domain", result.IdentityDomain)
	d.Set("ip_network", result.IPNetwork)
	d.Set("bring_your_own_license", result.IsBYOL)
	d.Set("level", result.Level)
	d.Set("national_character_set", result.NCharSet)
	d.Set("pdb_name", result.PDBName)
	d.Set("uri", result.URI)
	d.Set("shape", result.Shape)
	d.Set("sid", result.SID)
	d.Set("status", result.Status)
	d.Set("subnet", result.Subnet)
	d.Set("subscription_type", result.SubscriptionType)
	d.Set("timezone", result.Timezone)
	d.Set("version", result.Version)
	d.Set("exadata_system_name", result.SubscriptionName)
	d.Set("cluster_name", result.ClusterName)
	d.Set("node_list", result.NodeList)

	setAttributesFromConfig(d)

	// Obtain and set the default Access Rules
	getDefaultAccessRulesInput := &database.GetDefaultAccessRuleInput{
		ServiceInstanceID: d.Id(),
	}
	defaultAccessRules, err := dbClient.AccessRules().GetDefaultAccessRules(getDefaultAccessRulesInput)
	if err != nil {
		return err
	}
	if err = d.Set("default_access_rules", flattenDefaultAccessRules(defaultAccessRules)); err != nil {
		return fmt.Errorf("Error setting Database Default Access Rules: %+v", err)
	}

	return nil
}

// Certain values aren't received from the get call and need to be specified from the config
func setAttributesFromConfig(d *schema.ResourceData) {
	d.Set("disaster_recovery", d.Get("disaster_recovery"))

}

func resourceOPAASDatabaseServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()
	name := d.Id()

	client.Timeout = d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Deleting DatabaseServiceInstance: %v", name)

	input := database.DeleteServiceInstanceInput{
		Name: name,
	}
	if err := client.DeleteServiceInstance(&input); err != nil {
		return fmt.Errorf("Error deleting DatabaseServiceInstance: %+v", err)
	}
	return nil
}

func resourceOPAASDatabaseServiceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()

	if d.HasChange("desired_state") {
		updateInput := &database.DesiredStateInput{
			Name:           d.Id(),
			LifecycleState: database.ServiceInstanceLifecycleState(d.Get("desired_state").(string)),
		}

		_, err := client.UpdateDesiredState(updateInput)
		if err != nil {
			return fmt.Errorf("Unable to update Service Instance %q: %+v", d.Id(), err)
		}
	}

	if old, new := d.GetChange("shape"); old.(string) != "" && old.(string) != new.(string) {
		updateInput := &database.UpdateServiceInstanceInput{
			Name:  d.Id(),
			Shape: database.ServiceInstanceShape(new.(string)),
		}

		_, err := client.UpdateServiceInstance(updateInput)
		if err != nil {
			return fmt.Errorf("Unable to update Service Instance %q: %+v", d.Id(), err)
		}
	}

	err = updateDefaultAccessRules(d, meta)
	if err != nil {
		return fmt.Errorf("Unable to update Default Access Rules: %+v", err)
	}
	return resourceOPAASDatabaseServiceInstanceRead(d, meta)
}

func updateDefaultAccessRules(d *schema.ResourceData, meta interface{}) error {
	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.AccessRules()

	defaultAccessRuleConfig := d.Get("default_access_rules").([]interface{})
	if len(defaultAccessRuleConfig) == 0 {
		return nil
	}
	updateDefaultAccessRuleInput := &database.DefaultAccessRuleInfo{
		ServiceInstanceID: d.Id(),
	}

	defaultAccessRuleInfo := defaultAccessRuleConfig[0].(map[string]interface{})
	if val, ok := defaultAccessRuleInfo["enable_ssh"]; ok {
		updateDefaultAccessRuleInput.EnableSSH = helper.Bool(val.(bool))
	}
	if val, ok := defaultAccessRuleInfo["enable_http"]; ok {
		updateDefaultAccessRuleInput.EnableHTTP = helper.Bool(val.(bool))
	}
	if val, ok := defaultAccessRuleInfo["enable_http_ssl"]; ok {
		updateDefaultAccessRuleInput.EnableHTTPSSL = helper.Bool(val.(bool))
	}
	if val, ok := defaultAccessRuleInfo["enable_db_console"]; ok {
		updateDefaultAccessRuleInput.EnableDBConsole = helper.Bool(val.(bool))
	}
	if val, ok := defaultAccessRuleInfo["enable_db_express"]; ok {
		updateDefaultAccessRuleInput.EnableDBExpress = helper.Bool(val.(bool))
	}
	if val, ok := defaultAccessRuleInfo["enable_db_listener"]; ok {
		updateDefaultAccessRuleInput.EnableDBListener = helper.Bool(val.(bool))
	}
	if val, ok := defaultAccessRuleInfo["enable_em_console"]; ok {
		updateDefaultAccessRuleInput.EnableEMConsole = helper.Bool(val.(bool))
	}
	if val, ok := defaultAccessRuleInfo["enable_rac_db_listener"]; ok {
		updateDefaultAccessRuleInput.EnableRACDBListener = helper.Bool(val.(bool))
	}
	if val, ok := defaultAccessRuleInfo["enable_scan_listener"]; ok {
		updateDefaultAccessRuleInput.EnableScanListener = helper.Bool(val.(bool))
	}
	if val, ok := defaultAccessRuleInfo["enable_rac_ons"]; ok {
		updateDefaultAccessRuleInput.EnableRACOns = helper.Bool(val.(bool))
	}

	_, err = client.UpdateDefaultAccessRules(updateDefaultAccessRuleInput)
	if err != nil {
		return err
	}

	return nil
}

func expandStandby(d *schema.ResourceData) []database.StandBy {
	standbyConfig := d.Get("standby").([]interface{})
	attrs := standbyConfig[0].(map[string]interface{})
	standby := database.StandBy{
		AvailabilityDomain: attrs["availability_domain"].(string),
		Subnet:             attrs["subnet"].(string),
	}

	return []database.StandBy{standby}
}

func expandParameter(d *schema.ResourceData) (database.ParameterInput, error) {
	databaseConfigInfo := d.Get("database_configuration").([]interface{})
	attrs := databaseConfigInfo[0].(map[string]interface{})
	parameter := database.ParameterInput{
		AdminPassword:     attrs["admin_password"].(string),
		BackupDestination: database.ServiceInstanceBackupDestination(attrs["backup_destination"].(string)),
		DisasterRecovery:  attrs["disaster_recovery"].(bool),
		FailoverDatabase:  attrs["failover_database"].(bool),
		GoldenGate:        attrs["golden_gate"].(bool),
		IsRAC:             attrs["is_rac"].(bool),
		SID:               attrs["sid"].(string),
		Timezone:          attrs["timezone"].(string),
		Type:              database.ServiceInstanceType(attrs["type"].(string)),
	}

	if val, ok := attrs["usable_storage"].(int); ok && val != 0 {
		parameter.UsableStorage = strconv.Itoa(val)
	}
	if val, ok := attrs["snapshot_name"].(string); ok && val != "" {
		parameter.SnapshotName = val
	}
	if val, ok := attrs["source_service_name"].(string); ok && val != "" {
		parameter.SourceServiceName = val
	}
	if val, ok := attrs["character_set"].(string); ok && val != "" {
		parameter.CharSet = val
	}
	if val, ok := attrs["national_character_set"].(string); ok && val != "" {
		parameter.NCharSet = database.ServiceInstanceNCharSet(val)
	}
	if val, ok := attrs["pdb_name"].(string); ok && val != "" {
		parameter.PDBName = val
	}
	if val, ok := attrs["db_demo"].(string); ok && val != "" {
		addParam := database.AdditionalParameters{
			DBDemo: val,
		}
		parameter.AdditionalParameters = addParam
	}
	expandIbkup(d, &parameter)
	err := expandBackups(d, &parameter)
	if err != nil {
		return parameter, err
	}
	expandHDG(d, &parameter)

	return parameter, nil
}

func expandIbkup(d *schema.ResourceData, parameter *database.ParameterInput) {
	ibkupInfo := d.Get("instantiate_from_backup").([]interface{})
	if len(ibkupInfo) > 0 {
		attrs := ibkupInfo[0].(map[string]interface{})

		parameter.IBKUP = true
		parameter.IBKUPCloudStorageContainer = attrs["cloud_storage_container"].(string)
		parameter.IBKUPOnPremise = attrs["on_premise"].(bool)
		if val, ok := attrs["cloud_storage_username"]; ok {
			parameter.IBKUPCloudStorageUser = val.(string)
		}
		if val, ok := attrs["cloud_storage_password"]; ok {
			parameter.IBKUPCloudStoragePassword = val.(string)
		}
		if val, ok := attrs["decryption_key"]; ok {
			parameter.IBKUPDecryptionKey = val.(string)
		}
		if val, ok := attrs["service_id"]; ok {
			parameter.IBKUPServiceID = val.(string)
		}
		if val, ok := attrs["database_id"]; ok {
			parameter.IBKUPDatabaseID = val.(string)
		}
		if val, ok := attrs["wallet_file_content"]; ok {
			parameter.IBKUPWalletFileContent = val.(string)
		}
	}
}

func expandBackups(d *schema.ResourceData, parameter *database.ParameterInput) error {
	cloudStorageInfo := d.Get("backups").([]interface{})

	if parameter.BackupDestination == database.ServiceInstanceBackupDestinationBoth || parameter.BackupDestination == database.ServiceInstanceBackupDestinationOSS {
		if len(cloudStorageInfo) == 0 {
			return fmt.Errorf("`backups` must be set if `backup_destination` is set to `OSS` or `BOTH`")
		}
	}

	if len(cloudStorageInfo) > 0 {
		attrs := cloudStorageInfo[0].(map[string]interface{})
		parameter.CloudStorageContainer = attrs["cloud_storage_container"].(string)
		parameter.CreateStorageContainerIfMissing = attrs["create_if_missing"].(bool)
		if val, ok := attrs["cloud_storage_username"].(string); ok && val != "" {
			parameter.CloudStorageUsername = val
		}
		if val, ok := attrs["cloud_storage_password"].(string); ok && val != "" {
			parameter.CloudStoragePassword = val
		}
	}
	return nil
}

func expandHDG(d *schema.ResourceData, parameter *database.ParameterInput) error {
	hdgInfo := d.Get("hybrid_disaster_recovery").([]interface{})

	if len(hdgInfo) > 0 {
		if parameter.FailoverDatabase == true || parameter.IsRAC == true {
			return fmt.Errorf("`hybrid_disaster_recovery` cannot be set if `is_rac` or `failover_database` is set to true")
		}
		attrs := hdgInfo[0].(map[string]interface{})

		parameter.HDG = true
		parameter.HDGCloudStorageContainer = attrs["cloud_storage_container"].(string)
		// TODO read these values in the sdk like we do with cloud storage
		if val, ok := attrs["cloud_storage_username"].(string); ok && val != "" {
			parameter.HDGCloudStorageUser = val
		}
		if val, ok := attrs["cloud_storage_password"].(string); ok && val != "" {
			parameter.HDGCloudStoragePassword = val
		}
	}

	return nil
}

func flattenDefaultAccessRules(defaultAccessRules *database.DefaultAccessRuleInfo) []interface{} {
	result := make(map[string]interface{})

	if defaultAccessRules.EnableSSH != nil {
		result["enable_ssh"] = *defaultAccessRules.EnableSSH
	}
	if defaultAccessRules.EnableHTTP != nil {
		result["enable_http"] = *defaultAccessRules.EnableHTTP
	}
	if defaultAccessRules.EnableHTTPSSL != nil {
		result["enable_http_ssl"] = *defaultAccessRules.EnableHTTPSSL
	}
	if defaultAccessRules.EnableDBConsole != nil {
		result["enable_db_console"] = *defaultAccessRules.EnableDBConsole
	}
	if defaultAccessRules.EnableDBExpress != nil {
		result["enable_db_express"] = *defaultAccessRules.EnableDBExpress
	}
	if defaultAccessRules.EnableDBListener != nil {
		result["enable_db_listener"] = *defaultAccessRules.EnableDBListener
	}
	if defaultAccessRules.EnableEMConsole != nil {
		result["enable_em_console"] = *defaultAccessRules.EnableEMConsole
	}
	if defaultAccessRules.EnableRACDBListener != nil {
		result["enable_rac_db_listener"] = *defaultAccessRules.EnableRACDBListener
	}
	if defaultAccessRules.EnableScanListener != nil {
		result["enable_scan_listener"] = *defaultAccessRules.EnableScanListener
	}
	if defaultAccessRules.EnableRACOns != nil {
		result["enable_rac_ons"] = *defaultAccessRules.EnableRACOns
	}
	return []interface{}{result}
}
