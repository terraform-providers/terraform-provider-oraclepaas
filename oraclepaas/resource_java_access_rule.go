package oraclepaas

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"fmt"
	"log"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOraclePAASJavaAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOraclePAASJavaAccessRuleCreate,
		Read:   resourceOraclePAASJavaAccessRuleRead,
		Update: resourceOraclePAASJavaAccessRuleUpdate,
		Delete: resourceOraclePAASJavaAccessRuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(java.AccessRuleDestinationWLSAdmin),
					string(java.AccessRuleDestinationWLSAdminServer),
					string(java.AccessRuleDestinationOTD),
					string(java.AccessRuleDestinationOTDAdminHost),
				}, false),
			},
			"ports": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  java.AccessRuleProtocolTCP,
				ValidateFunc: validation.StringInSlice([]string{
					string(java.AccessRuleProtocolTCP),
					string(java.AccessRuleProtocolUDP),
				}, false),
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceOraclePAASJavaAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating database access rule")

	javaClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := javaClient.AccessRules()

	// Status can be enabled or disabled. We'll use `enabled` to determine status which to set
	enabled := d.Get("enabled").(bool)
	var status java.AccessRuleStatus
	if enabled == true {
		status = java.AccessRuleEnabled
	} else {
		status = java.AccessRuleDisabled
	}

	input := java.CreateAccessRuleInput{
		Name:              d.Get("name").(string),
		ServiceInstanceID: d.Get("service_instance_id").(string),
		Description:       d.Get("description").(string),
		Destination:       java.AccessRuleDestination(d.Get("destination").(string)),
		Ports:             d.Get("ports").(string),
		Protocol:          java.AccessRuleProtocol(d.Get("protocol").(string)),
		Source:            d.Get("source").(string),
		Status:            status,
	}

	info, err := client.CreateAccessRule(&input)
	if err != nil {
		return fmt.Errorf("Error creating Access Rule: %+v", err)
	}

	d.SetId(info.Name)

	return resourceOraclePAASJavaAccessRuleRead(d, meta)
}

func resourceOraclePAASJavaAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	javaClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := javaClient.AccessRules()

	log.Printf("[DEBUG] Reading state of access rules %q", d.Id())
	getInput := java.GetAccessRuleInput{
		Name:              d.Id(),
		ServiceInstanceID: d.Get("service_instance_id").(string),
	}

	result, err := client.GetAccessRule(&getInput)
	if err != nil {
		// AccessRule does not exist
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading database access rule %q: %+v", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of database access rule %q: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("service_instance_id", d.Get("service_instance_id"))
	d.Set("description", result.Description)
	d.Set("destination", result.Destination)
	d.Set("ports", result.Ports)
	d.Set("protocol", result.Protocol)
	d.Set("source", result.Source)
	d.Set("enabled", result.Status == java.AccessRuleEnabled)

	return nil
}

func resourceOraclePAASJavaAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Updating database access rule")

	javaClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := javaClient.AccessRules()

	// Status can be enabled or disabled. We'll use `enabled` to determine status which to set
	enabled := d.Get("enabled").(bool)
	var status java.AccessRuleStatus
	if enabled == true {
		status = java.AccessRuleEnabled
	} else {
		status = java.AccessRuleDisabled
	}
	input := java.UpdateAccessRuleInput{
		ServiceInstanceID: d.Get("service_instance_id").(string),
		Name:              d.Get("name").(string),
		Status:            status,
	}

	info, err := client.UpdateAccessRule(&input)
	if err != nil {
		return fmt.Errorf("Error updating Access Rule: %+v", err)
	}

	d.SetId(info.Name)

	return resourceOraclePAASJavaAccessRuleRead(d, meta)
}

func resourceOraclePAASJavaAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Deleting database access rule")

	javaClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := javaClient.AccessRules()

	// Status can be enabled or disabled. We'll use `enabled` to determine status which to set
	enabled := d.Get("enabled").(bool)
	var status java.AccessRuleStatus
	if enabled == true {
		status = java.AccessRuleEnabled
	} else {
		status = java.AccessRuleDisabled
	}

	input := java.DeleteAccessRuleInput{
		ServiceInstanceID: d.Get("service_instance_id").(string),
		Name:              d.Get("name").(string),
		Status:            status,
		Timeout:           d.Timeout(schema.TimeoutDelete),
	}

	err = client.DeleteAccessRule(&input)
	if err != nil {
		return fmt.Errorf("Error deleting Access Rule: %+v", err)
	}

	return nil
}
