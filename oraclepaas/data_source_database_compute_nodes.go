package oraclepaas

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOraclePAASDatabaseComputeNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOraclePAASDatabaseComputeNodesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"compute_nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connect_descriptor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connect_descriptor_with_public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"initial_primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pdb_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reserved_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shape": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_allocated": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOraclePAASDatabaseComputeNodesRead(d *schema.ResourceData, meta interface{}) error {
	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ComputeNodes()

	// Get required attributes
	name := d.Get("name").(string)

	input := database.GetComputeNodesInput{
		ServiceInstanceID: name,
	}

	result, err := client.GetComputeNodes(&input)
	if err != nil {
		return err
	}

	// Not found, don't error
	if result == nil {
		d.SetId("")
		return nil
	}

	d.SetId(name)
	d.Set("name", name)

	computeNodes, err := flattenComputeNodes(d, result.Nodes)
	// TODO not working!
	log.Printf(">>> computeNodes: %+v", computeNodes) // TODO remove
	if err != nil {
		return err
	}
	if err := d.Set("compute_nodes", computeNodes); err != nil {
		return fmt.Errorf("Error setting Compute Node info: %+v", err)
	}

	return nil
}

func flattenComputeNodes(d *schema.ResourceData, result []database.ComputeNodeInfo) ([]interface{}, error) {
	flattenedComputeNodes := make([]interface{}, 0)

	for _, info := range result {
		node := make(map[string]interface{})
		node["availability_domain"] = info.AvailabilityDomain
		node["connect_descriptor"] = info.ConnectDescriptor
		node["connect_descriptor_with_public_ip"] = info.ConnectDescriptorWithPublicIP
		node["hostname"] = info.Hostname
		node["initial_primary"] = info.InitialPrimary
		node["listener_port"] = info.ListenerPort
		node["pdb_name"] = info.PDBName
		node["reserved_ip"] = info.ReservedIP
		node["sid"] = info.SID
		node["status"] = info.Status
		node["storage_allocated"] = info.StorageAllocated
		node["subnet"] = info.Subnet
		flattenedComputeNodes = append(flattenedComputeNodes, node)
		log.Printf(">>> flattenedComputeNodes: %+v", flattenedComputeNodes) // TODO remove
	}

	return flattenedComputeNodes, nil
}
