package oraclepaas

import (
	"fmt"
	"log"
	"time"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOraclePAASMySQLAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOraclePAASMySQLAccessRuleCreate,
		Read:   resourceOraclePAASMySQLAccessRuleRead,
		Update: resourceOraclePAASMySQLAccessRuleUpdate,
		Delete: resourceOraclePAASMySQLAccessRuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"service_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			},
			// name validation: start with a letter, include letter, number and underscore only
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAccessRuleName,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		}, // end schema declaration
	} // end return
}

func resourceOraclePAASMySQLAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	mySQLClient, err := getMySQLClient(meta)
	if err != nil {
		return err
	}

	client := mySQLClient.AccessRules()

	var name = d.Get("name").(string)

	input := mysql.CreateAccessRuleInput{
		ServiceInstanceID: d.Get("service_instance_id").(string),
		RuleName:          name,
		Destination:       d.Get("destination").(string),
		Ports:             d.Get("ports").(string),
		Source:            d.Get("source").(string),
	}

	if d.Get("enabled").(bool) == true {
		input.Status = string(mysql.AccessRuleEnabled)
	} else {
		input.Status = string(mysql.AccessRuleDisabled)
	}

	if value, ok := d.GetOk("description"); ok {
		input.Description = value.(string)
	}

	if value, ok := d.GetOk("protocol"); ok {
		input.Protocol = value.(string)
	}

	err = client.CreateAccessRule(&input)

	if err != nil {
		return fmt.Errorf("Error creating Access Rule: %+v", err)
	}

	d.SetId(name)

	return resourceOraclePAASMySQLAccessRuleUpdate(d, meta)
}

func resourceOraclePAASMySQLAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	mySQLClient, err := getMySQLClient(meta)
	if err != nil {
		return err
	}

	client := mySQLClient.AccessRules()

	log.Printf("[DEBUG] Reading state of access rules %q", d.Id())

	input := mysql.GetAccessRuleInput{
		Name:              d.Id(),
		ServiceInstanceID: d.Get("service_instance_id").(string),
	}

	result, err := client.GetAccessRule(&input)

	if err != nil {
		// AccessRule does not exist
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading mysql access rule %q: %+v", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of mysql access rule %q: %#v", d.Id(), result)
	d.Set("name", result.RuleName)
	d.Set("service_instance_id", d.Get("service_instance_id"))
	d.Set("description", result.Description)
	d.Set("destination", result.Destination)
	d.Set("ports", result.Ports)
	d.Set("source", result.Source)
	d.Set("type", result.RuleType)
	d.Set("enabled", result.Status == string(mysql.AccessRuleEnabled))

	return nil
}

func resourceOraclePAASMySQLAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	mySQLClient, err := getMySQLClient(meta)
	if err != nil {
		return err
	}

	client := mySQLClient.AccessRules()

	enabled := d.Get("enabled").(bool)
	var status mysql.AccessRuleStatus
	if enabled == true {
		status = mysql.AccessRuleEnabled
	} else {
		status = mysql.AccessRuleDisabled
	}

	input := mysql.UpdateAccessRuleInput{
		ServiceInstanceID: d.Get("service_instance_id").(string),
		Name:              d.Get("name").(string),
		Status:            status,
	}

	info, err := client.UpdateAccessRule(&input)
	if err != nil {
		return fmt.Errorf("Error updating Access Rule: %+v", err)
	}

	d.SetId(info.RuleName)

	return resourceOraclePAASMySQLAccessRuleRead(d, meta)
}

func resourceOraclePAASMySQLAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	mySQLClient, err := getMySQLClient(meta)
	if err != nil {
		return err
	}

	client := mySQLClient.AccessRules()

	input := mysql.DeleteAccessRuleInput{
		ServiceInstanceID: d.Get("service_instance_id").(string),
		Name:              d.Get("name").(string),
		Operation:         mysql.AccessRuleDelete,
	}

	err = client.DeleteAccessRule(&input)
	if err != nil {
		return fmt.Errorf("Error deleting Access Rule: %+v", err)
	}
	return nil
}
