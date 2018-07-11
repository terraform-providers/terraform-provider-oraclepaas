package oraclepaas

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-oracle-terraform/application"
	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOraclePAASApplicationContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceOraclePAASApplicationContainerCreate,
		Read:   resourceOraclePAASApplicationContainerRead,
		Delete: resourceOraclePAASApplicationContainerDelete,
		Update: resourceOraclePAASApplicationContainerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manifest_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"manifest_attributes"},
			},
			"manifest_attributes": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"manifest_file"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"runtime": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"major_version": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(application.ManifestTypeWorker),
								string(application.ManifestTypeWeb),
							}, false),
							Default: string(application.ManifestTypeWorker),
						},
						"command": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"release": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"build": {
										Type:     schema.TypeString,
										Required: true,
									},
									"commit": {
										Type:     schema.TypeString,
										Required: true,
									},
									"version": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"startup_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(10, 600),
							Default:      30,
						},
						"shutdown_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 600),
							Default:      0,
						},
						"notes": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(application.ManifestModeRolling),
							}, false),
						},
						"clustered": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"home": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"health_check_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"deployment_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"manifest_attributes"},
			},
			"deployment_attributes": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"deployment_file"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instances": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"notes": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"environment": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"secure_environment": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"java_system_properties": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"services": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"identifier": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(application.ServiceTypeJAAS),
											string(application.ServiceTypeDBAAS),
											string(application.ServiceTypeMYSQLCS),
											string(application.ServiceTypeOEHCS),
											string(application.ServiceTypeOEHPCS),
											string(application.ServiceTypeDHCS),
											string(application.ServiceTypeCaching),
										}, false),
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
								},
							},
						},
					},
				},
			},
			"archive_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"basic",
					"oauth",
				}, false),
			},
			"notification_email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"dockerhub",
				}, false),
			},
			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "java",
				ValidateFunc: validation.StringInSlice([]string{
					"java",
					"node",
					"php",
					"python",
					"ruby",
				}, false),
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(application.SubscriptionTypeHourly),
				ValidateFunc: validation.StringInSlice([]string{
					string(application.SubscriptionTypeHourly),
					string(application.SubscriptionTypeMonthly),
				}, false),
			},
			"app_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOraclePAASApplicationContainerCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating application container")

	aClient, err := getApplicationClient(meta)
	if err != nil {
		return err
	}
	client := aClient.ContainerClient()

	additionalFields := application.CreateApplicationContainerAdditionalFields{
		Name:             d.Get("name").(string),
		SubscriptionType: d.Get("subscription_type").(string),
		Repository:       d.Get("repository").(string),
		Runtime:          d.Get("runtime").(string),
	}

	if v, ok := d.GetOk("archive_url"); ok {
		additionalFields.ArchiveURL = v.(string)
	}

	if v, ok := d.GetOk("notes"); ok {
		additionalFields.Notes = v.(string)
	}

	if v, ok := d.GetOk("notification_email"); ok {
		additionalFields.NotificationEmail = v.(string)
	}

	if v, ok := d.GetOk("auth_type"); ok {
		additionalFields.AuthType = v.(string)
	}

	input := application.CreateApplicationContainerInput{
		AdditionalFields: additionalFields,
	}

	if v, ok := d.GetOk("manifest_file"); ok {
		input.Manifest = v.(string)
	}

	if v, ok := d.GetOk("manifest_attributes"); ok {
		manifestAttrs, err := expandManifestAttributes(v.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		input.ManifestAttributes = manifestAttrs
	}

	if v, ok := d.GetOk("deployment_file"); ok {
		input.Deployment = v.(string)
	}

	if v, ok := d.GetOk("deployment_attributes"); ok {
		deploymentAttr, err := expandDeploymentAttributes(d, v.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		input.DeploymentAttributes = deploymentAttr
	}

	info, err := client.CreateApplicationContainer(&input)
	if err != nil {
		return fmt.Errorf("Error creating Application Container: %+v", err)
	}

	d.SetId(info.Name)
	return resourceOraclePAASApplicationContainerRead(d, meta)
}

func resourceOraclePAASApplicationContainerRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	aClient, err := getApplicationClient(meta)
	if err != nil {
		return err
	}
	client := aClient.ContainerClient()

	log.Printf("[DEBUG] Reading state of application container %s", d.Id())
	getInput := application.GetApplicationContainerInput{
		Name: d.Id(),
	}

	result, err := client.GetApplicationContainer(&getInput)
	if err != nil {
		// Application Container does not exist
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading application container %s: %+v", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of application container %s: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("app_url", result.AppURL)
	d.Set("web_url", result.WebURL)
	d.Set("runtime", d.Get("runtime").(string))
	d.Set("subscription_type", d.Get("subscription_type").(string))

	return nil
}

func resourceOraclePAASApplicationContainerDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	aClient, err := getApplicationClient(meta)
	if err != nil {
		return err
	}
	client := aClient.ContainerClient()
	name := d.Id()

	log.Printf("[DEBUG] Deleting ApplicationClient: %v", name)

	input := application.DeleteApplicationContainerInput{
		Name: name,
	}
	if err := client.DeleteApplicationContainer(&input); err != nil {
		return fmt.Errorf("Error deleting Application Container: %+v", err)
	}
	return nil
}

func resourceOraclePAASApplicationContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating application container")

	aClient, err := getApplicationClient(meta)
	if err != nil {
		return err
	}
	client := aClient.ContainerClient()

	additionalFields := application.UpdateApplicationContainerAdditionalFields{}

	if v, ok := d.GetOk("archive_url"); ok {
		additionalFields.ArchiveURL = v.(string)
	}

	if v, ok := d.GetOk("notes"); ok {
		additionalFields.Notes = v.(string)
	}

	input := application.UpdateApplicationContainerInput{
		Name:             d.Get("name").(string),
		AdditionalFields: additionalFields,
	}

	if v, ok := d.GetOk("manifest_file"); ok {
		input.Manifest = v.(string)
	}

	if v, ok := d.GetOk("manifest_attributes"); ok {
		manifestAttrs, err := expandManifestAttributes(v.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		input.ManifestAttributes = manifestAttrs
	}

	if v, ok := d.GetOk("deployment_file"); ok {
		input.Deployment = v.(string)
	}

	if v, ok := d.GetOk("deployment_attributes"); ok {
		deploymentAttr, err := expandDeploymentAttributes(d, v.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		input.DeploymentAttributes = deploymentAttr
	}

	info, err := client.UpdateApplicationContainer(&input)
	if err != nil {
		return fmt.Errorf("Error updating Application Container: %+v", err)
	}

	d.SetId(info.Name)
	return resourceOraclePAASApplicationContainerRead(d, meta)
}

func expandManifestAttributes(attrs map[string]interface{}) (*application.ManifestAttributes, error) {
	manifestAttributes := &application.ManifestAttributes{
		Type:         application.ManifestType(attrs["type"].(string)),
		StartupTime:  strconv.Itoa(attrs["startup_time"].(int)),
		ShutdownTime: strconv.Itoa(attrs["shutdown_time"].(int)),
	}
	if v := attrs["runtime"]; v != nil {
		runtimeAttrs := application.Runtime{}
		runtimeAttrs.MajorVersion = v.([]interface{})[0].(map[string]interface{})["major_version"].(string)
		manifestAttributes.Runtime = runtimeAttrs
	}
	if v := attrs["command"]; v != nil {
		manifestAttributes.Command = v.(string)
	}
	if v := attrs["release"]; v != nil {
		releaseAttrs := application.Release{}
		releaseConfig := v.([]interface{})[0].(map[string]interface{})
		releaseAttrs.Build = releaseConfig["build"].(string)
		releaseAttrs.Commit = releaseConfig["commit"].(string)
		releaseAttrs.Version = releaseConfig["version"].(string)
		manifestAttributes.Release = releaseAttrs
	}
	if v := attrs["notes"]; v != nil {
		manifestAttributes.Notes = v.(string)
	}
	if v := attrs["mode"]; v != nil {
		manifestAttributes.Mode = application.ManifestMode(v.(string))
	}
	if v := attrs["clustered"]; v != nil {
		manifestAttributes.IsClustered = v.(bool)
	}
	if v := attrs["home"]; v != nil {
		manifestAttributes.Home = v.(string)
	}
	if v := attrs["health_check"]; v != nil {
		manifestAttributes.HealthCheck = application.HealthCheck{HTTPEndpoint: v.(string)}
	}

	return manifestAttributes, nil
}

func expandDeploymentAttributes(d *schema.ResourceData, attrs map[string]interface{}) (*application.DeploymentAttributes, error) {
	deploymentAttributes := &application.DeploymentAttributes{}

	if v := attrs["memory"]; v != nil {
		deploymentAttributes.Memory = v.(string)
	}
	if v := attrs["instances"]; v != nil {
		deploymentAttributes.Instances = v.(int)
	}
	if v := attrs["notes"]; v != nil {
		deploymentAttributes.Notes = v.(string)
	}
	if v := attrs["environment"]; v != nil {
		environment := make(map[string]string)
		for name, value := range v.(map[string]interface{}) {
			environment[name] = value.(string)
		}
		deploymentAttributes.Envrionment = environment
	}
	if v := attrs["secure_environment"]; v != nil {
		deploymentAttributes.SecureEnvironment = getStringList(d, "deployment_attributes.0.secure_environment")
	}
	if v := attrs["java_system_properties"]; v != nil {
		jsp := make(map[string]string)
		for name, value := range v.(map[string]interface{}) {
			jsp[name] = value.(string)
		}
		deploymentAttributes.Envrionment = jsp
	}
	if v := attrs["services"]; v != nil {
		deploymentAttributes.Services = expandServices(v.([]interface{}))
	}
	return deploymentAttributes, nil
}

func expandServices(attrs []interface{}) []application.Service {
	services := make([]application.Service, 0, len(attrs))

	for i, serviceAttr := range attrs {
		serviceConfig := serviceAttr.(map[string]interface{})
		service := application.Service{
			Identifier: serviceConfig["identifier"].(string),
			Type:       application.ServiceType(serviceConfig["type"].(string)),
			Name:       serviceConfig["name"].(string),
			Username:   serviceConfig["username"].(string),
			Password:   serviceConfig["password"].(string),
		}
		services[i] = service
	}
	return services
}
