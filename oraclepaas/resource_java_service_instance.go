package oraclepaas

import (
	"fmt"
	"log"
	"strings"
	"time"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOraclePAASJavaServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOraclePAASJavaServiceInstanceCreate,
		Read:   resourceOraclePAASJavaServiceInstanceRead,
		Delete: resourceOraclePAASJavaServiceInstanceDelete,
		Update: resourceOraclePAASJavaServiceInstanceUpdate,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssh_public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(java.ServiceInstanceLevelPAAS),
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceLevelPAAS),
					string(java.ServiceInstanceLevelBasic),
				}, false),
			},
			"edition": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceEditionEE),
					string(java.ServiceInstanceEditionSE),
					string(java.ServiceInstanceEditionSuite),
				}, false),
			},
			"service_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backups": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_storage_container": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"auto_generate": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
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
					},
				},
			},
			"metering_frequency": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(java.ServiceInstanceSubscriptionTypeHourly),
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceSubscriptionTypeHourly),
					string(java.ServiceInstanceSubscriptionTypeMonthly),
				}, false),
			},
			"availability_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"weblogic_server": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  7001,
										ForceNew: true,
									},
									"secured_port": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  7002,
										ForceNew: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"application_database": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"password": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"pdb_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"backup_volume_size": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cluster_name": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
							ConflictsWith: []string{"weblogic_server.0.cluster"},
						},
						"cluster": {
							Type:          schema.TypeList,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"weblogic_server.0.cluster_name"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(java.ServiceInstanceClusterTypeApplication),
											string(java.ServiceInstanceClusterTypeCaching),
										}, false),
									},
									"server_count": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  1,
									},
									"servers_per_node": {
										Type:         schema.TypeInt,
										Optional:     true,
										ForceNew:     true,
										Default:      1,
										ValidateFunc: validation.IntBetween(1, 8),
									},
									"path_prefixes": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"shape": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(javaServiceInstanceShapes(), false),
									},
								},
							},
						},
						"connect_string": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"database": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pdb_name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"domain": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Default:  string(java.ServiceInstanceDomainModePro),
										ValidateFunc: validation.StringInSlice([]string{
											string(java.ServiceInstanceDomainModeDev),
											string(java.ServiceInstanceDomainModePro),
										}, false),
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Computed: true,
									},
									"partition_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.IntBetween(0, 4),
									},
									"volume_size": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"ip_reservations": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"managed_servers": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.IntBetween(1, 8),
										Default:      1,
									},
									"initial_heap_size": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"max_heap_size": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"jvm_args": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"initial_permanent_generation": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"max_permanent_generation": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"overwrite_jvm_args": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
										Default:  false,
									},
								},
							},
						},
						"middleware_volume_size": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"node_manager": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:      schema.TypeString,
										Optional:  true,
										ForceNew:  true,
										Computed:  true,
										Sensitive: true,
									},
									"password": {
										Type:      schema.TypeString,
										Optional:  true,
										ForceNew:  true,
										Computed:  true,
										Sensitive: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  5556,
									},
								},
							},
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"privileged_content_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  80,
									},
									"privileged_secured_content_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  443,
									},
									"deployment_channel_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  9001,
									},
									"content_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  8001,
									},
								},
							},
						},
						"shape": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(javaServiceInstanceShapes(), false),
						},
						"upper_stack_product_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(java.ServiceInstanceUpperStackProductNameODI),
								string(java.ServiceInstanceUpperStackProductNameWCP),
							}, false),
						},
						"root_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"oracle_traffic_director": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Default:  8989,
										Optional: true,
										ForceNew: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"high_availability": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"ip_reservations": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"listener": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"secured_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Computed: true,
									},
									"privileged_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Computed: true,
									},
									"privileged_secured_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Computed: true,
									},
								},
							},
						},
						"load_balancing_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(java.ServiceInstanceLoadBalancingPolicyLCC),
							ValidateFunc: validation.StringInSlice([]string{
								string(java.ServiceInstanceLoadBalancingPolicyLCC),
								string(java.ServiceInstanceLoadBalancingPolicyLRT),
								string(java.ServiceInstanceLoadBalancingPolicyRR),
							}, false),
						},
						"shape": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"root_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"backup_destination": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(java.ServiceInstanceBackupDestinationBoth),
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceBackupDestinationBoth),
					string(java.ServiceInstanceBackupDestinationNone),
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_admin_console": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"notification_email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bring_your_own_license": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
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
			"use_identity_service": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
		},
	}
}

func resourceOraclePAASJavaServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating JavaServiceInstance")

	jClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := jClient.ServiceInstanceClient()

	input := java.CreateServiceInstanceInput{
		ServiceName:        d.Get("name").(string),
		ServiceLevel:       java.ServiceInstanceLevel(d.Get("level").(string)),
		VMPublicKeyText:    d.Get("ssh_public_key").(string),
		Edition:            java.ServiceInstanceEdition(d.Get("edition").(string)),
		BackupDestination:  java.ServiceInstanceBackupDestination(d.Get("backup_destination").(string)),
		EnableAdminConsole: d.Get("enable_admin_console").(bool),
		UseIdentityService: d.Get("use_identity_service").(bool),
	}

	if val, ok := d.GetOk("service_version"); ok {
		input.ServiceVersion = val.(string)
	}

	if val, ok := d.GetOk("metering_frequency"); ok {
		input.MeteringFrequency = java.ServiceInstanceSubscriptionType(val.(string))
	}

	if val, ok := d.GetOk("description"); ok {
		input.ServiceDescription = val.(string)
	}
	if val, ok := d.GetOk("ip_network"); ok {
		input.IPNetwork = val.(string)
	}
	if val, ok := d.GetOk("region"); ok {
		input.Region = val.(string)
	}
	if val, ok := d.GetOk("bring_your_own_license"); ok {
		input.IsBYOL = val.(bool)
	}
	if val, ok := d.GetOk("notification_email"); ok {
		input.EnableNotification = true
		input.NotificationEmail = val.(string)
	}
	if val, ok := d.GetOk("availability_domain"); ok {
		input.AvailabilityDomain = val.(string)
	}
	if val, ok := d.GetOk("snapshot_name"); ok {
		input.SnapshotName = val.(string)
	}
	if val, ok := d.GetOk("source_service_name"); ok {
		input.SourceServiceName = val.(string)
	}
	if val, ok := d.GetOk("use_identity_service"); ok {
		input.UseIdentityService = val.(bool)
	}
	if val, ok := d.GetOk("subnet"); ok {
		input.Subnet = val.(string)
	}

	expandJavaCloudStorage(d, &input)
	expandWebLogicConfig(d, &input)
	expandOTDConfig(d, &input)

	info, err := client.CreateServiceInstance(&input)
	if err != nil {
		return fmt.Errorf("Error creating JavaServiceInstance: %s", err)
	}

	d.SetId(info.ServiceName)
	return resourceOraclePAASJavaServiceInstanceRead(d, meta)
}

func resourceOraclePAASJavaServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	jClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := jClient.ServiceInstanceClient()

	log.Printf("[DEBUG] Reading state of Java Service Instance %s", d.Id())
	getInput := java.GetServiceInstanceInput{
		Name: d.Id(),
	}

	result, err := client.GetServiceInstance(&getInput)
	if err != nil {
		// Java Service Instance does not exist
		if opcClient.WasNotFoundError(err) || strings.Contains(err.Error(), "No such service") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading JavaServiceInstance %s: %s", d.Id(), err)
	}

	if result == nil {
		log.Printf("[DEBUG] Unable to find Java Service Instance %s", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of JavaServiceInstance %s: %#v", d.Id(), result)
	d.Set("name", result.ServiceName)
	d.Set("level", result.ServiceLevel)
	d.Set("edition", result.Edition)
	d.Set("version", result.ServiceVersion)
	d.Set("metering_frequency", result.MeteringFrequency)
	d.Set("force_delete", d.Get("force_delete"))

	wlsConfig, err := flattenWebLogicConfig(d, result.Components.WLS, result.WLSRoot)
	if err != nil {
		return err
	}
	if err := d.Set("weblogic_server", wlsConfig); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Java Service instance WebLogic Server: %+v", err)
	}
	/* if err := d.Set("oracle_traffic_director", flattenOTDConfig(d, result.Components.OTD, result.OTDRoot)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Java Service Instance Oracle Traffic Director: %+v", err)
	} */

	return nil
}

func resourceOraclePAASJavaServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	jClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := jClient.ServiceInstanceClient()
	name := d.Id()

	log.Printf("[DEBUG] Deleting JavaServiceInstance: %q", name)

	// Need to get the dba username and password to delete the service instance
	webLogicConfig := d.Get("weblogic_server").([]interface{})
	webLogicAttrs := webLogicConfig[0].(map[string]interface{})
	dbaInfo := webLogicAttrs["database"].([]interface{})

	attrs := dbaInfo[0].(map[string]interface{})
	username := attrs["username"].(string)
	password := attrs["password"].(string)

	input := java.DeleteServiceInstanceInput{
		Name:        name,
		DBAUsername: username,
		DBAPassword: password,
		ForceDelete: d.Get("force_delete").(bool),
	}

	if err := client.DeleteServiceInstance(&input); err != nil {
		return fmt.Errorf("Error deleting JavaServiceInstance: %+v", err)
	}
	return nil
}

func resourceOraclePAASJavaServiceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	jClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := jClient.ServiceInstanceClient()

	// Updating the shape refers to changing the shape of the admin cluster for the weblogic server
	if old, new := d.GetChange("weblogic_server.0.shape"); old.(string) != "" && old.(string) != new.(string) {
		wlsComponent := java.ScaleUpDownWLS{
			Hosts: []string{d.Get("weblogic_server.0.admin.0.hostname").(string)},
			Shape: java.ServiceInstanceShape(old.(string)),
		}
		updateInput := &java.ScaleUpDownServiceInstanceInput{
			Name:       d.Id(),
			Components: java.ScaleUpDownComponent{WLS: wlsComponent},
		}

		err := client.ScaleUpDownServiceInstance(updateInput)
		if err != nil {
			return err
		}
	}

	return resourceOraclePAASJavaServiceInstanceRead(d, meta)
}

func expandWebLogicConfig(d *schema.ResourceData, input *java.CreateServiceInstanceInput) {
	webLogicConfig := d.Get("weblogic_server").([]interface{})
	webLogicServer := &java.CreateWLS{}

	attrs := webLogicConfig[0].(map[string]interface{})
	webLogicServer.Shape = java.ServiceInstanceShape(attrs["shape"].(string))
	expandWLSAdmin(webLogicServer, attrs)
	expandAppDBs(webLogicServer, attrs)
	expandClusters(d, webLogicServer, attrs)
	expandDB(webLogicServer, attrs)
	expandDomain(webLogicServer, attrs)
	expandManagedServers(webLogicServer, attrs)
	expandNodeManager(webLogicServer, attrs)
	expandWLSPorts(webLogicServer, attrs)

	if v := attrs["backup_volume_size"]; v != nil {
		webLogicServer.BackupVolumeSize = v.(string)
	}
	if v := attrs["cluster_name"]; v != nil {
		webLogicServer.ClusterName = v.(string)
	}
	if v := attrs["ip_reservations"]; v != nil {
		webLogicServer.IPReservations = getStringList(d, "weblogic_server.0.ip_reservations")
	}
	if v := attrs["middleware_volume_size"]; v != nil {
		webLogicServer.MWVolumeSize = v.(string)
	}
	if v := attrs["upper_stack_product_name"]; v != nil {
		webLogicServer.UpperStackProductName = java.ServiceInstanceUpperStackProductName(v.(string))
	}
	input.Components = java.CreateComponents{WLS: webLogicServer}
}

func expandOTDConfig(d *schema.ResourceData, input *java.CreateServiceInstanceInput) {
	otdConfig := d.Get("oracle_traffic_director").([]interface{})

	if len(otdConfig) == 0 {
		return
	}

	otdInfo := &java.CreateOTD{}
	attrs := otdConfig[0].(map[string]interface{})

	otdInfo.Shape = java.ServiceInstanceShape(attrs["shape"].(string))
	expandOTDAdmin(otdInfo, attrs)
	expandListener(otdInfo, attrs)

	if v := attrs["high_availability"]; v != nil {
		otdInfo.HAEnabled = v.(bool)
	}
	if v := attrs["ip_reservations"]; v != nil {
		otdInfo.IPReservations = getStringList(d, "oracle_traffic_director.0.ip_reservations")
	}
	if v := attrs["load_balancing_policy"]; v != nil {
		otdInfo.LoadBalancingPolicy = java.ServiceInstanceLoadBalancingPolicy(v.(string))
	}

	input.Components.OTD = otdInfo
}

func expandWLSAdmin(webLogicServer *java.CreateWLS, config map[string]interface{}) {
	adminInfo := config["admin"].([]interface{})
	attrs := adminInfo[0].(map[string]interface{})

	webLogicServer.AdminUsername = attrs["username"].(string)
	webLogicServer.AdminPassword = attrs["password"].(string)
	if v := attrs["port"]; v != nil {
		webLogicServer.AdminPort = v.(int)
	}
	if v := attrs["secured_port"]; v != nil {
		webLogicServer.SecuredAdminPort = v.(int)
	}
}

func expandClusters(d *schema.ResourceData, weblogicServer *java.CreateWLS, config map[string]interface{}) {
	clusterInfo := config["cluster"].([]interface{})
	for i, clusterConfig := range clusterInfo {
		cluster := java.CreateCluster{}
		attrs := clusterConfig.(map[string]interface{})

		cluster.ClusterName = attrs["name"].(string)
		cluster.Type = java.ServiceInstanceClusterType(attrs["type"].(string))
		cluster.ServerCount = attrs["server_count"].(int)
		cluster.ServerPerNode = attrs["servers_per_node"].(int)
		if v := attrs["shape"]; v != nil {
			cluster.Shape = java.ServiceInstanceShape(v.(string))
		}
		if v := attrs["path_prefixes"]; v != nil {
			cluster.PathPrefixes = getStringList(d, fmt.Sprintf("weblogic_server.0.clusters.%d.path_prefixes", i))
		}
		weblogicServer.Clusters = append(weblogicServer.Clusters, cluster)
	}
}

func expandDB(webLogicServer *java.CreateWLS, config map[string]interface{}) {
	dbaInfo := config["database"].([]interface{})
	if len(dbaInfo) == 0 {
		return
	}

	attrs := dbaInfo[0].(map[string]interface{})
	webLogicServer.DBServiceName = attrs["name"].(string)
	webLogicServer.DBAName = attrs["username"].(string)
	webLogicServer.DBAPassword = attrs["password"].(string)
	if v := attrs["pdb_name"]; v != nil {
		webLogicServer.PDBServiceName = v.(string)
	}
}

func expandOTDAdmin(otdServer *java.CreateOTD, config map[string]interface{}) {
	adminInfo := config["admin"].([]interface{})
	attrs := adminInfo[0].(map[string]interface{})

	otdServer.AdminUsername = attrs["username"].(string)
	otdServer.AdminPassword = attrs["password"].(string)
	if v := attrs["port"]; v != nil {
		otdServer.AdminPort = v.(int)
	}
}

func expandJavaCloudStorage(d *schema.ResourceData, input *java.CreateServiceInstanceInput) {
	cloudStorageInfo := d.Get("backups").([]interface{})

	attrs := cloudStorageInfo[0].(map[string]interface{})
	input.CloudStorageContainer = attrs["cloud_storage_container"].(string)
	input.CloudStorageContainerAutoGenerate = attrs["auto_generate"].(bool)
	if val, ok := attrs["cloud_storage_username"].(string); ok && val != "" {
		input.CloudStorageUsername = val
	}
	if val, ok := attrs["cloud_storage_password"].(string); ok && val != "" {
		input.CloudStoragePassword = val
	}
}

func expandAppDBs(webLogicServer *java.CreateWLS, config map[string]interface{}) {
	appDBInfo := config["application_database"].(*schema.Set)
	appDBs := make([]java.AppDB, appDBInfo.Len())
	for i, val := range appDBInfo.List() {
		attrs := val.(map[string]interface{})
		appDB := java.AppDB{
			DBAName:       attrs["username"].(string),
			DBAPassword:   attrs["password"].(string),
			DBServiceName: attrs["name"].(string),
		}
		if v := attrs["pdb_name"]; v != nil {
			appDB.PDBServiceName = v.(string)
		}
		appDBs[i] = appDB
	}
	webLogicServer.AppDBs = appDBs
}

func expandDomain(webLogicServer *java.CreateWLS, config map[string]interface{}) {
	domainInfo := config["domain"].([]interface{})

	if len(domainInfo) == 0 {
		return
	}

	attrs := domainInfo[0].(map[string]interface{})

	webLogicServer.DomainMode = java.ServiceInstanceDomainMode(attrs["mode"].(string))
	if val, ok := attrs["name"].(string); ok && val != "" {
		webLogicServer.DomainName = val
	}
	if val, ok := attrs["partition_count"].(int); ok {
		webLogicServer.DomainPartitionCount = val
	}
	if val, ok := attrs["volume_size"].(string); ok && val != "" {
		webLogicServer.DomainVolumeSize = val
	}
}

func expandListener(otdInfo *java.CreateOTD, config map[string]interface{}) {
	listenerInfo := config["listener"].([]interface{})

	if len(listenerInfo) == 0 {
		return
	}

	attrs := listenerInfo[0].(map[string]interface{})
	if v := attrs["port"]; v != nil {
		otdInfo.ListenerPort = v.(int)
	}
	if v := attrs["secured_port"]; v != nil {
		otdInfo.SecuredListenerPort = v.(int)
	}
	if v := attrs["privileged_port"]; v != nil {
		otdInfo.PrivilegedListenerPort = v.(int)
	}
	if v := attrs["privileged_secured_port"]; v != nil {
		otdInfo.PrivilegedSecuredListenerPort = v.(int)
	}
	if otdInfo.ListenerPort != 0 || otdInfo.PrivilegedListenerPort != 0 {
		otdInfo.ListenerPortEnabled = true
	} else {
		otdInfo.ListenerPortEnabled = false
	}
}

func expandManagedServers(webLogicServer *java.CreateWLS, config map[string]interface{}) {
	msInfo := config["managed_servers"].([]interface{})

	if len(msInfo) == 0 {
		return
	}

	attrs := msInfo[0].(map[string]interface{})
	if val, ok := attrs["server_count"]; ok {
		webLogicServer.ManagedServerCount = val.(int)
	}
	if val, ok := attrs["initial_heap_size"]; ok {
		webLogicServer.MSInitialHeapMB = val.(int)
	}
	if val, ok := attrs["max_heap_size"]; ok {
		webLogicServer.MSMaxHeapMB = val.(int)
	}
	if val, ok := attrs["jvm_args"]; ok && val != "" {
		webLogicServer.MSJvmArgs = val.(string)
	}
	if val, ok := attrs["initial_permanent_generation"]; ok {
		webLogicServer.MSPermMB = val.(int)
	}
	if val, ok := attrs["max_permanent_generation"]; ok {
		webLogicServer.MSMaxPermMB = val.(int)
	}
	if val, ok := attrs["overwrite_jvm_args"]; ok {
		webLogicServer.OverwriteMsJvmArgs = val.(bool)
	}
}

func expandNodeManager(webLogicServer *java.CreateWLS, config map[string]interface{}) {
	nmInfo := config["node_manager"].([]interface{})

	if len(nmInfo) == 0 {
		return
	}

	attrs := nmInfo[0].(map[string]interface{})
	webLogicServer.NodeManagerPort = attrs["port"].(int)
	if val, ok := attrs["password"].(string); ok && val != "" {
		webLogicServer.NodeManagerPassword = val
	}
	if val, ok := attrs["username"].(string); ok && val != "" {
		webLogicServer.NodeManagerUserName = val
	}
}

func expandWLSPorts(webLogicServer *java.CreateWLS, config map[string]interface{}) {
	portInfo := config["ports"].([]interface{})

	if len(portInfo) == 0 {
		return
	}

	attrs := portInfo[0].(map[string]interface{})
	if v := attrs["privileged_content_port"]; v != nil {
		webLogicServer.PrivilegedContentPort = v.(int)
	}
	if v := attrs["privileged_secured_content_port"]; v != nil {
		webLogicServer.PrivilegedSecuredContentPort = v.(int)
	}
	if v := attrs["deployment_channel_port"]; v != nil {
		webLogicServer.DeploymentChannelPort = v.(int)
	}
	if v := attrs["content_port"]; v != nil {
		webLogicServer.ContentPort = v.(int)
	}
}

func flattenWebLogicConfig(d *schema.ResourceData, webLogicConfig java.WLS, rootURL string) ([]interface{}, error) {
	result := make(map[string]interface{})

	result["shape"] = d.Get("weblogic_server.0.shape")
	// Hostname is the only thing related to the Admin block that is returned
	result["admin"] = flattenWLSAdmin(d, webLogicConfig.AdminHostName)
	result["database"] = flattenDatabase(d)
	result["domain"] = flattenDomain(d)
	result["managed_servers"] = flattenManagedServers(d)
	result["node_manager"] = flattenNodeManager(d)
	result["ports"] = flattenWLSPorts(d)

	v := flattenAppDB(d)
	if v != nil {
		result["application_database"] = v
	}
	if v, ok := d.GetOk("weblogic_server.0.cluster_name"); ok {
		result["cluster_name"] = webLogicConfig.Clusters[v.(string)].ClusterName
	}
	if _, ok := d.GetOk("weblogic_server.0.cluster"); ok {
		clusters, err := flattenClusters(d, webLogicConfig.Clusters)
		if err != nil {
			return nil, err
		}
		result["cluster"] = clusters
	}
	if v, ok := d.GetOk("weblogic_server.0.backup_volume_size"); ok {
		result["backup_volume_size"] = v
	}
	if v, ok := d.GetOk("weblogic_server.0.connect_string"); ok {
		result["connect_string"] = v
	}
	if rootURL != "" {
		result["root_url"] = rootURL
	}
	if _, ok := d.GetOk("weblogic_server.0.ip_reservations"); ok {
		result["ip_reservations"] = getStringList(d, "weblogic_server.0.ip_reservations")
	}
	if v, ok := d.GetOk("middleware_volume_size"); ok {
		result["middleware_volume_size"] = v
	}
	if v, ok := d.GetOk("weblogic_server.0.upper_stack_product_name"); ok {
		result["upper_stack_product_name"] = v
	}

	return []interface{}{result}, nil
}

/*
func flattenOTDConfig(d *schema.ResourceData, otdConfig java.OTD, rootURL string) []interface{} {
	result := make(map[string]interface{})

	if d.Get("otd.0.shape") == nil {
		return []interface{}{}
	} else {
		result["root_url"] = rootURL
		result["admin"] = flattenOTDAdmin(d, otdConfig.AdminHostName)
		result["high_availability"] = d.Get("otd.0.high_availability")
		result["listener"] = flattenListener(d)
		result["load_balancing_policy"] = d.Get("otd.0.load_balancing_policy")
		result["shape"] = d.Get("otd.0.shape")

		if _, ok := d.GetOk("otd.0.ip_reservations"); ok {
			result["ip_reservations"] = getStringList(d, "otd.0.ip_reservations")
		}
	}

	return []interface{}{result}
}*/

// Only adminHostname is returned from the api forcing the other attributes to be reset
// here from the schema.
func flattenWLSAdmin(d *schema.ResourceData, adminHostname string) []interface{} {
	admin := make(map[string]interface{})
	admin["hostname"] = adminHostname

	// Setting variables that don't get returned from the api
	admin["username"] = d.Get("weblogic_server.0.admin.0.username")
	admin["password"] = d.Get("weblogic_server.0.admin.0.password")
	admin["port"] = d.Get("weblogic_server.0.admin.0.port")
	if v, ok := d.GetOk("weblogic_server.0.admin.0.secured_port"); ok {
		admin["secured_port"] = v
	}

	return []interface{}{admin}
}

// Only adminHostname is returned from the api forcing the other attributes to be reset
// here from the schema.
func flattenOTDAdmin(d *schema.ResourceData, adminHostname string) []interface{} {
	admin := make(map[string]interface{})
	admin["hostname"] = adminHostname

	// Setting variables that don't get returned from the api
	if v, ok := d.GetOk("otd.0.admin.0.username"); ok {
		admin["username"] = v
	}
	if v, ok := d.GetOk("otd.0.admin.0.password"); ok {
		admin["password"] = v
	}
	if v, ok := d.GetOk("otd.0.admin.0.port"); ok {
		admin["port"] = v
	}

	return []interface{}{admin}
}

// AppDBs are not returned by the api forcing them to be reset here from the schema.
func flattenAppDB(d *schema.ResourceData) []interface{} {
	appDBInfo := make([]map[string]interface{}, 0)

	appDBs := d.Get("weblogic_server.0.application_database").(*schema.Set)

	if len(appDBs.List()) == 0 {
		return nil
	}
	for _, val := range appDBs.List() {
		appDBInfo = append(appDBInfo, val.(map[string]interface{}))
	}

	return []interface{}{appDBInfo}
}

func flattenClusters(d *schema.ResourceData, clusters map[string]java.Clusters) ([]interface{}, error) {
	flattenedClusters := make([]interface{}, 0)

	clustersInfo := d.Get("weblogic_server.0.cluster").([]interface{})
	for i, clusterConfig := range clustersInfo {
		attrs := clusterConfig.(map[string]interface{})
		clusterInfo := clusters[attrs["name"].(string)]
		if clusterInfo.ClusterName != "" {
			cluster := make(map[string]interface{})
			cluster["name"] = clusterInfo.ClusterName
			cluster["type"] = d.Get(fmt.Sprintf("weblogic_server.0.cluster.%d.type", i))
			cluster["server_count"] = d.Get(fmt.Sprintf("weblogic_server.0.cluster.%d.server_count", i))
			cluster["servers_per_node"] = d.Get(fmt.Sprintf("weblogic_server.0.cluster.%d.servers_per_node", i))

			if v, ok := d.GetOk(fmt.Sprintf("weblogic_server.0.cluster.%d.shape", i)); ok {
				cluster["shape"] = v.(string)
			}
			if _, ok := d.GetOk(fmt.Sprintf("weblogic_server.0.cluster.%d.path_prefixes", i)); ok {
				cluster["path_prefixes"] = getStringList(d, fmt.Sprintf("weblogic_server.0.cluster.%d.path_prefixes", i))
			}
			flattenedClusters = append(flattenedClusters, cluster)
		}
	}

	return flattenedClusters, nil
}

func flattenDatabase(d *schema.ResourceData) []interface{} {
	db := make(map[string]interface{})
	db["username"] = d.Get("weblogic_server.0.database.0.username")
	db["password"] = d.Get("weblogic_server.0.database.0.password")
	db["name"] = d.Get("weblogic_server.0.database.0.name")

	if v, ok := d.GetOk("weblogic_server.0.pdb_name"); ok {
		db["pdb_name"] = v
	}

	return []interface{}{db}
}

func flattenDomain(d *schema.ResourceData) []interface{} {
	domain := make(map[string]interface{})
	domainConfig := d.Get("weblogic_server.0.domain").([]interface{})

	if len(domainConfig) == 0 {
		return []interface{}{domain}
	}
	if domainConfig[0] != nil {
		attrs := domainConfig[0].(map[string]interface{})

		domain["mode"] = java.ServiceInstanceDomainMode(attrs["mode"].(string))
		if val, ok := attrs["name"].(string); ok && val != "" {
			domain["name"] = val
		}
		if val, ok := attrs["partition_count"].(int); ok {
			domain["partition_count"] = val
		}
		if val, ok := attrs["volume_size"].(string); ok && val != "" {
			domain["volume_size"] = val
		}
	}

	return []interface{}{domain}
}

func flattenManagedServers(d *schema.ResourceData) []interface{} {
	managedServers := make(map[string]interface{})
	managedServerConfig := d.Get("weblogic_server.0.managed_servers").([]interface{})
	if len(managedServerConfig) == 0 {
		return nil
	}
	if managedServerConfig[0] != nil {
		attrs := managedServerConfig[0].(map[string]interface{})
		if val, ok := attrs["server_count"]; ok {
			managedServers["server_count"] = val
		}
		if val, ok := attrs["initial_heap_size"]; ok {
			managedServers["initial_heap_size"] = val
		}
		if val, ok := attrs["max_heap_size"]; ok {
			managedServers["max_heap_size"] = val
		}
		if val, ok := attrs["jvm_args"]; ok && val != "" {
			managedServers["jvm_args"] = val
		}
		if val, ok := attrs["initial_permanent_generation"]; ok {
			managedServers["initial_permanent_generation"] = val
		}
		if val, ok := attrs["max_permanent_generation"]; ok {
			managedServers["max_permanent_generation"] = val
		}
		if val, ok := attrs["overwrite_jvm_args"]; ok {
			managedServers["overwrite_jvm_args"] = val
		}
	}

	return []interface{}{managedServers}
}

func flattenNodeManager(d *schema.ResourceData) []interface{} {
	nodeManager := make(map[string]interface{})
	nodeManagerConfig := d.Get("weblogic_server.0.node_manager").([]interface{})

	if len(nodeManagerConfig) == 0 {
		return nil
	}
	if nodeManagerConfig[0] != nil {
		attrs := nodeManagerConfig[0].(map[string]interface{})
		nodeManager["port"] = attrs["port"].(int)
		if val, ok := attrs["password"].(string); ok && val != "" {
			nodeManager["password"] = val
		}
		if val, ok := attrs["username"].(string); ok && val != "" {
			nodeManager["username"] = val
		}
	}

	return []interface{}{nodeManager}
}

func flattenWLSPorts(d *schema.ResourceData) []interface{} {
	ports := make(map[string]interface{})
	portsConfig := d.Get("weblogic_server.0.ports").([]interface{})

	if len(portsConfig) == 0 {
		return nil
	}
	if portsConfig[0] != nil {
		attrs := portsConfig[0].(map[string]interface{})
		ports["privileged_content_port"] = attrs["privileged_content_port"]
		ports["privileged_secured_content_port"] = attrs["privileged_secured_content_port"]
		ports["deployment_channel_port"] = attrs["deployment_channel_port"]
		ports["content_port"] = attrs["content_port"]
	}

	return []interface{}{ports}
}

func flattenListener(otdAttributes java.OTDAttributes) []interface{} {
	listener := make(map[string]interface{})

	listener["port"] = otdAttributes.ListenerPort
	listener["privileged_port"] = otdAttributes.PrivilgedListenerPort
	listener["privileged_secured_port"] = otdAttributes.PrivilegedSecureListenerPort

	return []interface{}{listener}
}
