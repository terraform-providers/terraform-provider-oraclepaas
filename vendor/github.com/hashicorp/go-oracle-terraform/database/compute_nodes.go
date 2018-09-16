package database

// API URI Paths for Container and Root objects
const (
	DBComputeNodeContainerPath = "paas/service/dbcs/api/v1.1/instances/%s/%s/servers"
	DBComputeNodeRootPath      = "paas/service/dbcs/api/v1.1/instances/%s/%s/servers/%s"
)

// ComputeNodes returns a UtilityClient for managing Server Nodes for a DBaaS Service Instance
func (c *Client) ComputeNodes() *UtilityClient {
	return &UtilityClient{
		UtilityResourceClient: UtilityResourceClient{
			Client:           c,
			ContainerPath:    DBComputeNodeContainerPath,
			ResourceRootPath: DBComputeNodeRootPath,
		},
	}
}

// ComputeNodeInfo returns the details of a single Compute Node
type ComputeNodeInfo struct {
	// Name of the availability domain within the region where the compute node is provisioned.
	// OCI only.
	AvailabilityDomain string `json:"availabilityDomain"`
	// The connection descriptor for Oracle Net Services (SQL*Net).
	ConnectDescriptor string `json:"connect_descriptor"`
	// The connection descriptor for Oracle Net Services (SQL*Net) with IP addresses instead of host names.
	ConnectDescriptorWithPublicIP string `json:"connect_descriptor_with_public_ip"`
	// The user name of the Oracle Cloud user who created the service instance.
	CreatedBy string `json:"created_by"`
	// The job id of the job that created the service instance.
	CreationJobID string `json:"creation_job_id"`
	// The date-and-time stamp when the service instance was created.
	CreationTime string `json:"creation_time"`
	// The host name of the compute node.
	Hostname string `json:"hostname"`
	// Indicates whether the compute node hosted the primary database of an Oracle Data Guard
	// configuration when the service instance was created.
	InitialPrimary bool `json:"initialPrimary"`
	// The listener port for Oracle Net Services (SQL*Net) connections.
	ListenerPort int `json:"listenerPort"`
	// The size in GB of the memory allocated to the compure node.
	// Exadata only.
	MemoryAllocated int `json:"memoryAllocated"`
	// Number of CPU Cores.
	// Exadata only.
	NumberOfCores int `json:"numberOfCores"`
	// The name of the default PDB (pluggable database) created when the service instance was created.
	PDBName string `json:"pdbName"`
	// The IP address of the compute node.
	ReservedIP string `json:"reservedIP"`
	// The Oracle Compute Cloud shape of the compute node.
	Shape string `json:"shape"`
	// The SID of the database on the compute node.
	SID string `json:"sid"`
	// The status of the compute node.
	Status string `json:"status"`
	// The size in GB of the storage allocated to the compute node.
	// For compute nodes of a service instance hosting an Oracle RAC database, this number does not
	// include the storage shared by the nodes
	StorageAllocated int `json:"storageAllocated"`
	// Name of the subnet within the region where the compute node is provisioned.
	Subnet string `json:"subnet"`
	// Virtual Machine type
	VMType string `json:"vmType"`
}

// ComputeNodesInfo - contains details of a Compute Nodes for a service instance
type ComputeNodesInfo struct {
	// List of Compute Nodes
	Nodes []ComputeNodeInfo `json:"-"`
}

// GetComputeNodesInput Get Request input to query compute nodes for a service instance
type GetComputeNodesInput struct {
	// Name of the DBaaS service instance.
	// Required
	ServiceInstanceID string `json:"-"`
}

// GetComputeNodes gets details of all Compute Nodes for a Service Instance
func (c *UtilityClient) GetComputeNodes(input *GetComputeNodesInput) (*ComputeNodesInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}
	var computeNodes []ComputeNodeInfo
	if err := c.getResource("", &computeNodes); err != nil {
		return nil, err
	}
	return &ComputeNodesInfo{
		Nodes: computeNodes,
	}, nil
}
