package oraclepaas

import (
	"fmt"
	"log"
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_file": {
				Type:     schema.TypeString,
				Optional: true,
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

	if v, ok := d.GetOk("deployment_file"); ok {
		input.Deployment = v.(string)
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

	if v, ok := d.GetOk("deployment_file"); ok {
		input.Deployment = v.(string)
	}

	info, err := client.UpdateApplicationContainer(&input)
	if err != nil {
		return fmt.Errorf("Error updating Application Container: %+v", err)
	}

	d.SetId(info.Name)
	return resourceOraclePAASApplicationContainerRead(d, meta)
}
