package oraclepaas

import (
	"time"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// The Exadata Service resource is a variant of the Database Cloud Service. The exadata
// resource defines the *schema only* to enforce Exadata Cloud Service
// specific attributes, valiations, and removes unsupported attributes.
// All CRUD logic is part of the `database_service_instance` implementation

func resourceOraclePAASExadataServiceInstance() *schema.Resource {
	return &schema.Resource{
		// use Database Service CRUD operations
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
				// for Exadata on `MONTHLY` is supported
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
						"db_demo": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
			"desired_state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceLifecycleStateStop),
					string(database.ServiceInstanceLifecycleStateRestart),
					string(database.ServiceInstanceLifecycleStateStart),
				}, true),
			},
			"compute_site_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbaas_monitor_url": {
				Type:     schema.TypeString,
				Computed: true,
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
		},
	}
}

// see `resource_database_service_instance.go` for CRUD operations
