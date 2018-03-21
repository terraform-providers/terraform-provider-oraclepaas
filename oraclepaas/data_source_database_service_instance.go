package oraclepaas

import (
	"fmt"
	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOraclePAASDatabaseServiceInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOraclePAASDatabaseServiceInstanceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"apex_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_destination": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"character_set": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_storage_container": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compute_site_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_manager_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failover_database": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"glassfish_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hybrid_disaster_recovery_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_reservations": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bring_your_own_license": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"monitor_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"national_character_set": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pluggable_database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shape": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"high_performance_storage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOraclePAASDatabaseServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()

	// Get required attributes
	name := d.Get("name").(string)

	input := database.GetServiceInstanceInput{
		Name: name,
	}

	result, err := client.GetServiceInstance(&input)
	if err != nil {
		return err
	}

	// Not found, don't error
	if result == nil {
		d.SetId("")
		return nil
	}

	// Populate schema attributes
	d.SetId(fmt.Sprintf("%s", result.Name))
	d.Set("name", result.Name)
	d.Set("apex_url", result.ApexURL)
	d.Set("availability_domain", result.AvailabilityDomain)
	d.Set("backup_destination", result.BackupDestination)
	d.Set("character_set", result.CharSet)
	d.Set("compute_site_name", result.ComputeSiteName)
	d.Set("description", result.Description)
	d.Set("edition", result.Edition)
	d.Set("enterprise_manager_url", result.EMURL)
	d.Set("failover_database", result.FailoverDatabase)
	d.Set("glassfish_url", result.GlassFishURL)
	d.Set("hybrid_disaster_recovery_ip", result.HDGPremIP)
	d.Set("identity_domain", result.IdentityDomain)
	d.Set("ip_network", result.IPNetwork)
	d.Set("ip_reservations", result.IPReservations)
	d.Set("bring_your_own_license", result.IsBYOL)
	d.Set("level", result.Level)
	d.Set("listener_port", result.ListenerPort)
	d.Set("monitor_url", result.DBAASMonitorURL)
	d.Set("national_character_set", result.NCharSet)
	d.Set("pluggable_database_name", result.PDBName)
	d.Set("region", result.Region)
	d.Set("shape", result.Shape)
	d.Set("high_performance_storage", result.UseHighPerformanceStorage)
	d.Set("uri", result.URI)
	d.Set("version", result.Version)

	return nil
}
