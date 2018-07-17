//
package mysql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

// WaitForServiceInstanceReadyPollInterval is the default polling interval value for Creating a service instance and waiting for the instance to be ready
const WaitForServiceInstanceReadyPollInterval = time.Duration(60 * time.Second)

// WaitForServiceInstanceReadyTimeout is the default Timeout value for Creating a service instance and waiting for the instance to be ready
const WaitForServiceInstanceReadyTimeout = time.Duration(3600 * time.Second)

// WaitForServiceInstanceDeletePollInterval is the default polling value for deleting a service instance and waiting for the instance to be completely removed.
const WaitForServiceInstanceDeletePollInterval = time.Duration(60 * time.Second)

// WaitForServiceInstanceDeleteTimeout is the default Timeout value for deleting a service instance and waiting for the instance to be completely removed.
const WaitForServiceInstanceDeleteTimeout = time.Duration(3600 * time.Second)

// API URI Paths for the Container and Root objects.
var (
	ServiceInstanceContainerPath = "/paas/api/v1.1/instancemgmt/%[1]s/services/MySQLCS/instances/"
	ServiceInstanceResourcePath  = "/paas/api/v1.1/instancemgmt/%[1]s/services/MySQLCS/instances/%[2]s"
)

// ServiceInstanceClient is a client for the Service functions of the MySQL API.
type ServiceInstanceClient struct {
	ResourceClient
	Timeout      time.Duration
	PollInterval time.Duration
}

// ServiceInstanceClient obtains an ServiceInstanceClient which can be used to access to the
// Service Instance functions of the Database Cloud API
func (c *MySQLClient) ServiceInstanceClient() *ServiceInstanceClient {
	return &ServiceInstanceClient{
		ResourceClient: ResourceClient{
			MySQLClient:      c,
			ContainerPath:    ServiceInstanceContainerPath,
			ResourceRootPath: ServiceInstanceResourcePath,
		}}
}

// Constants for whether the Enterprise Monitor should be installed
type ServiceInstanceEnterpriseMonitor string

const (
	ServiceInstanceEnterpriseMonitorYes ServiceInstanceEnterpriseMonitor = "Yes"
	ServiceInstanceEnterpriseMonitorNo  ServiceInstanceEnterpriseMonitor = "No"
)

// Constants for the metering frequency for the MySQL CS Service Instance.
type ServiceInstanceMeteringFrequency string

const (
	ServiceInstanceMeteringFrequencyHourly  ServiceInstanceMeteringFrequency = "HOURLY"
	ServiceInstanceMeteringFrequencyMonthly ServiceInstanceMeteringFrequency = "MONTHLY"
)

// Constants for the Backup Destination
type ServiceInstanceBackupDestination string

const (
	ServiceInstanceBackupDestinationBoth ServiceInstanceBackupDestination = "BOTH"
	ServiceInstanceBackupDestinationNone ServiceInstanceBackupDestination = "NONE"
	ServiceInstanceBackupDestinationOSS  ServiceInstanceBackupDestination = "OSS"
)

// Constants for the state of the Service Instance State
type ServiceInstanceState string

const (
	ServiceInstanceReady        ServiceInstanceState = "READY"
	ServiceInstanceInitializing ServiceInstanceState = "INITIALIZING"
	ServiceInstanceStarting     ServiceInstanceState = "STARTING"
	ServiceInstanceStopping     ServiceInstanceState = "STOPPING"
	ServiceInstanceStopped      ServiceInstanceState = "STOPPED"
	ServiceInstanceConfiguring  ServiceInstanceState = "CONFIGURING"
	ServiceInstanceError        ServiceInstanceState = "ERROR"
	ServiceInstanceTerminating  ServiceInstanceState = "TERMINATING"
)

// ActivityLogInfo describes the list of activities that have occurred on the ServiceInstance.
type ActivityLogInfo struct {
	ActivityLogId  string                `json:"activityLogId"`
	AuthDomain     string                `json:"authDomain"`
	AuthUser       string                `json:"authUser"`
	EndDate        string                `json:"endDate"`
	IdentityDomain string                `json:"identityDomain"`
	InitiatedBy    string                `json:"initiatedBy"`
	JobId          string                `json:"jobId"`
	Messages       []ActivityMessageInfo `json:"messages"`
	OperationId    string                `json:"operationId"`
	OperationType  string                `json:"operationType"`
	ServiceId      string                `json:"serviceId"`
	ServiceName    string                `json:"serviceName"`
	StartDate      string                `json:"startDate"`
	Status         string                `json:"status"`
	SummaryMessage string                `json:"summaryMessage"`
	ServiceType    string                `json:"serviceType"` // Not in API
}

// ActivityMessageInfo is the specific message that
type ActivityMessageInfo struct {
	ActivityDate string `json:"activityDate"`
	Messages     string `json:"message"`
}

type AttributeInfo struct {
	DisplayName  string `json:"displayName"`
	Type         string `json:"type"`
	Value        string `json:"value"`
	DisplayValue string `json:"displayValue"`
	IsKeyBinding bool   `json:"isKeyBinding"`
}

// ServiceInstance defines the instance information that is returned from the Get method
// when quering the instance.
type ServiceInstance struct {
	ServiceId                    string                   `json:"serviceId"`
	ServiceUuid                  string                   `json:"serviceUuid"` // Not in API
	ServiceLogicalUuid           string                   `json:"serviceLogicalUuid"`
	ServiceName                  string                   `json:"serviceName"`
	ServiceType                  string                   `json:"serviceType"`
	DomainName                   string                   `json:"domainName"`
	ServiceVersion               string                   `json:"serviceVersion"`
	ReleaseVersion               string                   `json:"releaseVersion"`
	BaseReleaseVersion           string                   `json:"baseReleaseVersion"` // Not in API
	MetaVersion                  string                   `json:"metaVersion"`
	ServiceDescription           string                   `json:"serviceDescription"` // Not in API
	ServiceLevel                 string                   `json:"serviceLevel"`
	Subscription                 string                   `json:"subscription"`
	MeteringFrequency            string                   `json:"meteringFrequency"`
	Edition                      string                   `json:"edition"`
	TotalSSDStorage              int                      `json:"totalSSDStorage"`
	Status                       ServiceInstanceState     `json:"state"`
	ServiceStateDisplayName      string                   `json:"serviceStateDisplayName"`
	Clone                        bool                     `json:"clone"`
	Creator                      string                   `json:"creator"`
	CreationDate                 string                   `json:"creationDate"`
	IsBYOL                       bool                     `json:"isBYOL"`
	IsManaged                    bool                     `json:"isManaged"`
	IaasProvider                 string                   `json:"iaasProvider"`
	Attributes                   map[string]AttributeInfo `json:"attributes"`
	Components                   ComponentInfo            `json:"components"`
	ActivityLogs                 []ActivityLogInfo        `json:"activityLogs"`
	LayeringMode                 string                   `json:"layeringMode"`
	ServiceLevelDisplayName      string                   `json:"serviceLevelDisplayName"`
	EditionDisplayName           string                   `json:"editionDisplayName"`
	MeteringFrequencyDisplayName string                   `json:"meteringFrequencyDisplayName"`
	BackupFilePath               string                   `json:"BACKUP_FILE_PATH"`
	DataVolumeSize               string                   `json:"DATA_VOLUME_SIZE"`
	UseSSD                       string                   `json:"USE_SSD"`
	ProvisionEngine              string                   `json:"provisionEngine"`
	MysqlPort                    string                   `json:"MYSQL_PORT"`
	CloudStorageContainer        string                   `json:"CLOUD_STORAGE_CONTAINER"`
	BackupDestination            string                   `json:"BACKUP_DESTINATION"`
	TotalSharedStorage           int                      `json:"totalSharedStorage"`
	ComputeSiteName              string                   `json:"computeSiteName"`
	Patching                     PatchingInfo             `json:"patching"`

	// The reason for the instance going to error state, if available.
	ErrorReason string `json:"error_reason"`
}

type MysqlInfo struct {
	ServiceId                 string                    `json:"serviceId"`
	ComponentId               string                    `json:"componentId"`
	State                     string                    `json:"state"`
	ComponentStateDisplayName string                    `json:"componentStateDisplayName"`
	Version                   string                    `json:"version"`
	ComponentType             string                    `json:"componentType"` // Not in API
	CreationDate              string                    `json:"creationDate"`
	InstanceName              string                    `json:"instanceName"`
	InstanceRole              string                    `json:"instanceRole"`
	IsKeyComponent            bool                      `json:"isKeyComponent"` // Not in API
	Attributes                map[string]AttributeInfo  `json:"attributes"`
	VMInstances               map[string]VMInstanceInfo `json:"vmInstances"`
	AdminHostName             string                    `json:"adminHostName"`
	Hosts                     map[string]HostInfo       `json:"hosts"`       // Not in API
	DisplayName               string                    `json:"displayName"` // Not in API
	// hosts
	// paasServers
}

type VMInstanceInfo struct {
	VmId               string `json:"vmId"`
	Id                 int    `json:"id"`
	Uuid               string `json:"uuid"`
	HostName           string `json:"hostName"`
	Label              string `json:"label"`
	IPAddress          string `json:"ipAddress"`
	PublicIPAddress    string `json:"publicIpAddress"`
	UsageType          string `json:"usageType"`
	Role               string `json:"role"`
	ComponentType      string `json:"componentType"`
	State              string `json:"state"`
	VmStateDisplayName string `json:"vmStateDisplayName"`
	ShapeId            string `json:"shapeId"`
	TotalStorage       int    `json:"totalStorage"`
	CreationDate       string `json:"creationDate"`
	IsAdminNode        bool   `json:"isAdminNode"`
}

type HostInfo struct {
	Vmid               int                   `json:"vmId"`
	Id                 int                   `json:"id"`
	Uuid               string                `json:"uuid"`
	HostName           string                `json:"hostName"`
	Label              string                `json:"label"`
	UsageType          string                `json:"usageType"`
	Role               string                `json:"role"`
	ComponentType      string                `json:"componentType"`
	State              string                `json:"state"`
	VMStateDisplayName string                `json:"vmStateDisplayName"`
	ShapeId            string                `json:"shapeId"`
	TotalStorage       int                   `json:"totalStorage"`
	CreationDate       string                `json:"creationDate"`
	IsAdminNode        bool                  `json:"isAdminNode"`
	Servers            map[string]ServerInfo `json:"servers"`
}

type ServerInfo struct {
	ServerId               string `json:"serverId"`
	ServerName             string `json:"serverName"`
	ServerType             string `json:"serverType"`
	ServerRole             string `json:"serverRole"`
	State                  string `json:"state"`
	ServerStateDisplayName string `json:"serverStateDisplayName"`
	CreationDate           string `json:"creationDate"`
}

type StorageVolumeInfo struct {
	Name       string `json:"name"`
	Partitions string `json:"partitions"`
	Size       string `json:"size"`
}

type ComponentInfo struct {
	Mysql MysqlInfo `json:"mysql"`
}

type PatchingInfo struct {
	CurrentOperation      map[string]string `json:"currentOperation"`
	TotalAvailablePatches string            `json:"totalAvailablePatches"`
}

// Details for the CreateServiceInstance to create the MySQL Service Instance
type CreateServiceInstanceInput struct {
	// MySQL Component parameters for they MySQL Service.
	ComponentParameters ComponentParameters `json:"componentParameters"`
	// Service parameters for the MySQL Service
	ServiceParameters ServiceParameters `json:"serviceParameters"`
}

// ComponentParameters used for creating the MySQL Service Instance. This wraps the MySQLParamters.
type ComponentParameters struct {
	Mysql MySQLParameters `json:"mysql"`
}

// MySQLParameters used for create the MySQL Service Instance.
type MySQLParameters struct {
	// The name of the MySQL Database. This defaults to mydatabase if the value is omitted or blank.
	DBName string `json:"dbName,omitempty"`
	// The Storage Volume size (in GB) for the MySQL Data. The value must be between 25 and 1024. The default value is 25.
	DBStorage string `json:"dbStorage, omitempty"`
	// Indicate whether the MySQL Enterprise Monitor should be configured. Values : [ "Yes, "No"]. The default is "No"
	EnterpriseMonitor string `json:"enterpriseMonitor,omitempty"`
	// Password for the EM Agent. The password must be at least 8 characters long, and have at least one lower case letter, one upper case letter, one number and one special character
	EnterpriseMonitorAgentPassword string `json:"enterpriseMonitorAgentPassword,omitempty"`
	// Username for the EM Agent. The Name must start with a letter, and consist of letters and numbers. and be between 2 and 32 characters.
	EnterpriseMonitorAgentUser string `json:"enterpriseMonitorAgentUser,omitempty"`
	// Password for the EM Manager. The password must be at least 8 characters long, and have at least one lower case letter, one upper case letter, one number and one special character
	EnterpriseMonitorManagerPassword string `json:"enterpriseMonitorManagerPassword,omitempty"`
	// Username for the EM Manager. The Name must start with a letter, and consist of letters and numbers. and be between 2 and 32 characters.
	EnterpriseMonitorManagerUser string `json:"enterpriseMonitorManagerUser,omitempty"`
	// The MySQL Server Character set. Default Value: 'utbmb4'
	MysqlCharset string `json:"mysqlCharset,omitempty"`
	// The MySQL Server collation.
	MysqlCollation string `json:"mysqlCollation,omitempty"`
	// The Port for the MySQL Enterprise Monitor. The default is 18443
	MysqlEMPort string `json:"mysqlEMPort,omitempty"`
	// The port for the MySQL Service. The default is 3306
	MysqlPort string `json:"mysqlPort,omitempty"`

	// The MySQL Server time zone. The default is SYSTEM. The value can be  given as a named time zone, such as "Europe/Paris, or "Asia/Shanghai"
	// Although this is in the API, the REST APIS are throwing an error that the parameter is invalid
	// MysqlTimezone string `json:"mysqlTimezone, omitempty"`
	// MySQL server options and variables. Only comma separated key value pairs with no spaces are permitted (e.g., option1=value,option2=value). MySQL server options that are available as MySQL Server Component Parameters, such as mysqlPort, are not permitted in a mysqlOptions string.
	// Although this is in the API, the REST APIS are throwing an error that the parameter is invalid
	// MysqlOptions string `json:"mysqlOptions,omitempty"`

	// The MySQL Administration user for connecting to the service. The Name must start with a letter, and consist of letters and numbers. and be between 2 and 32 characters. Default Value: root.
	MysqlUserName string `json:"mysqlUserName,omitempty"`
	// The password for the MySQL Username. The password must start with a letter, be between 8 and 30 characters long, and contains letters, at least one number, and any number of special chacters ($#_)
	MysqlUserPassword string `json:"mysqlUserPassword,omitempty"`
	// Desired compute shape. Default: oc3
	Shape string `json:"shape,omitempty"`
	// The name of the snapshot of the service instance specified by sourceServiceName that is to be used to create a "snapshot clone". This parameter is valid only if sourceServiceName is specified.
	SnapshotName string `json:"snapshot,omitempty"`
	//  indicates that the service instance should be created as a "snapshot clone" of another service instance. Provide the name of the existing service instance whose snapshot is to be used. dbName, mysqlCharset, mysqlCollation, mysqlEMPort, enterpriseMonitor, and associated MySQL server component parameters do not apply when cloning a service from a snapshot.
	SourceServiceName string `json:"sourceServiceName,omitempty"`
	// This attribute is relevant to only Oracle Cloud Infrastructure. Specify the Oracle Cloud Identifier (OCID) of a subnet from a virtual cloud network (VCN) that you had created previously in Oracle Cloud Infrastructure. For the instructions to create a VCN and subnet
	Subnet string `json:"subnet,omitempty"`
}

// ServiceParameters details the service parameters for the create instance operation
type ServiceParameters struct {
	// ONLY FOR OCI. Name of the data center location for the OCI Region. e.g. FQCn:US-ASHBURN-AD1"
	AvailabilityDomain string `json:"availabilityDomain,omitempty"`
	// Backup Destination. Value values are : BOTH, OSS, NONE. Default: NONE
	BackupDestination string `json:"backupDestination,omitempty"`
	// The URI of the object storage container for storing the MySQL CS Instance backups. On OCI-C, the container does NOT need to be created before provisioning. On OCI the container MUST BE created before provisioning.
	// Use one of the following formats: Storage-<IdentityDomainID>/<Container Name>, <StorageServiceName>-<IdentityDomainId>/<ContainerName>, <restEndPointURL>/<ContainerName>
	CloudStorageContainer string `json:"cloudStorageContainer,omitempty"`
	// Specifies whether to creat the storage container if it does not exist. Not applicable to OCI. Only for OCI-C. Default: False
	CloudStorageContainerAutoGenerate bool `json:"cloudStorageContainerAutoGenerate,omitempty"`
	// Password for the object storage user. The password must be specified if cloudStorageContainer is set.
	CloudStoragePassword string `json:"cloudStoragePassword,omitempty"`
	// User name for the object storage user. The user name must be specified if cloudStorageContainer is set.
	CloudStorageUsername string `json:"cloudStorageUser,omitempty"`
	// Flag that specifies whether to enable (true) or disable (false) notifications by email. If this property is set to true, you must specify a value in notificationEmail.
	EnableNotification bool `json:"enableNotification,omitempty"`
	// The three-part name of a custom IP network to attach this service instance to. For example: /Compute-identity_domain/user/object. This attribute is applicable only to accounts where regions are supported. This is not applicable to OCI
	IPNetwork string `json:"ipNetwork"`
	// The billing frequency of the service instance; either MONTHLY or HOURLY. Default: MONTHLY
	MeteringFrequency string `json:"meteringFrequency,omitempty"`
	// The email that will be used to send notifications to.
	NotificationEmail string `json:"notificationEmail,omitempty"`
	// Name of the region where the MySQL Service instance is to be provisioned. This attribute is applicable only to accounts where regions are supported
	Region string `json:"region,omitempty"`
	// Text that provides addition information about the service instance.
	ServiceDescription string `json:"serviceDescription,omitempty"`
	// Name of the Service Instance. The name must be between 1 to 50 characters, start with a letter, contain only letters, numbers or hyphens, must not end with a hyphen, and must be unique within the identity domain.
	ServiceName string `json:"serviceName"`
	// Public key for the secure shell (SSH). This key will be used for authentication when connecting to the MySQL Cloud Service instance using an SSH client.
	VMPublicKeyText string `json:"vmPublicKeyText,omitempty"`
	// VM operating system user that is valid for variations of compute based services. It will default to the username opc when not specified.
	VMUser string `json:"vmUser,omitempty"`
}

// CreateServiceInstance calls the MySQL CS APIs to create the service instance. The method is used internally by the startServiceInstance method.
// The method returns a http 202 on success.
func (c *ServiceInstanceClient) CreateServiceInstance(input *CreateServiceInstanceInput) (*ServiceInstance, error) {
	var (
		serviceInstance      *ServiceInstance
		serviceInstanceError error
	)

	if c.PollInterval == 0 {
		c.PollInterval = WaitForServiceInstanceReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = WaitForServiceInstanceReadyTimeout
	}

	// Since these CloudStorageUsername and CloudStoragePassword are sensitive we'll read them
	// from the client if they haven't specified in the config.
	if input.ServiceParameters.CloudStorageContainer != "" && input.ServiceParameters.CloudStorageUsername == "" && input.ServiceParameters.CloudStoragePassword == "" {
		input.ServiceParameters.CloudStorageUsername = *c.ResourceClient.MySQLClient.client.UserName
		input.ServiceParameters.CloudStoragePassword = *c.ResourceClient.MySQLClient.client.Password
	}

	for i := 0; i < *c.MySQLClient.client.MaxRetries; i++ {
		serviceInstance, serviceInstanceError = c.startServiceInstance(input.ServiceParameters.ServiceName, input)
		if serviceInstanceError == nil {
			return serviceInstance, nil
		}
	}

	return nil, serviceInstanceError
}

// startServiceInstance calls the CreateServiceInstance method to create the MySQL Service Instance, then calls the WaitForServiceInstance to wait util the MySQL Instance is ready and accessible.
func (c *ServiceInstanceClient) startServiceInstance(name string, input *CreateServiceInstanceInput) (*ServiceInstance, error) {
	if err := c.createResource(*input, nil); err != nil {
		return nil, err
	}

	// Call wait for instance ready now, as creating the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: name,
	}

	serviceInstance, serviceInstanceError := c.WaitForServiceInstanceRunning(getInput, c.PollInterval, c.Timeout)

	if serviceInstanceError != nil {
		c.client.DebugLogString(fmt.Sprintf(": Create Failed %s", serviceInstanceError))
		return nil, serviceInstanceError
	}

	return serviceInstance, nil
}

// GetServiceInstanceInput defines the parameters needed to retrieve information on ServiceInstance.
type GetServiceInstanceInput struct {
	// Name of the MySQL Cloud Service instance.
	// Required.
	Name string `json:"serviceId"`
}

// GetServiceInstance retrieves the ServiceInstance with the given name.
func (c *ServiceInstanceClient) GetServiceInstance(getInput *GetServiceInstanceInput) (*ServiceInstance, error) {
	var serviceInstance ServiceInstance
	if err := c.getResource(getInput.Name, &serviceInstance); err != nil {
		return nil, err
	}

	return &serviceInstance, nil
}

// WaitForServiceInstanceRunning waits for an instance to be created and completely initialized and available.
func (c *ServiceInstanceClient) WaitForServiceInstanceRunning(input *GetServiceInstanceInput, pollingInterval time.Duration, timeoutSeconds time.Duration) (*ServiceInstance, error) {
	var info *ServiceInstance
	var getErr error

	err := c.client.WaitFor("service instance to be ready", pollingInterval, timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetServiceInstance(input)
		if getErr != nil {
			return false, getErr
		}

		switch s := info.Status; s {

		case ServiceInstanceReady: // Target State
			c.client.DebugLogString(fmt.Sprintf("ServiceInstance [%s] is in ready state. Status : %s", info.ServiceId, info.Status))
			return true, nil
		case ServiceInstanceInitializing:
			c.client.DebugLogString(fmt.Sprintf("ServiceInstance [%s] is in Initializing state. Status : %s", info.ServiceId, info.Status))
			return false, nil
		case ServiceInstanceStarting:
			c.client.DebugLogString(fmt.Sprintf("ServiceInstance [%s] is in Starting state. Status : %s", info.ServiceId, info.Status))
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("ServiceInstance [%s] is in DEFAULT. Status : %s", info.ServiceId, info.Status))
			return false, nil
		}
	})
	return info, err
}

// DeleteServiceInput defines the parameters needed to delete a MySQL Instance.
type DeleteServiceInput struct {
	//Options string `json:"options,omitempty"`
}

// DeleteServiceInstance delete the MySQL instance, then waits for the actual instance to be removed before returning.
func (c *ServiceInstanceClient) DeleteServiceInstance(serviceName string) error {
	if c.Timeout == 0 {
		c.Timeout = WaitForServiceInstanceDeleteTimeout
	}
	if c.PollInterval == 0 {
		c.PollInterval = WaitForServiceInstanceDeletePollInterval
	}

	c.client.DebugLogString(fmt.Sprintf("Deleting Instance : %s", serviceName))

	deleteInput := &DeleteServiceInput{}

	deleteErr := c.deleteResource(serviceName, deleteInput)
	if deleteErr != nil {
		c.client.DebugLogString(fmt.Sprintf(": Delete Failed %s", deleteErr))
		return deleteErr
	}

	// Call wait for instance deleted now, as deleting the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: serviceName,
	}

	// Wait for instance to be deleted
	return c.WaitForServiceInstanceDeleted(getInput, c.PollInterval, c.Timeout)
}

// WaitForServiceInstanceDeleted waits for a service instance to be fully deleted.
func (c *ServiceInstanceClient) WaitForServiceInstanceDeleted(input *GetServiceInstanceInput, pollingInterval time.Duration, timeoutSeconds time.Duration) error {
	return c.client.WaitFor("service instance to be deleted", pollingInterval, timeoutSeconds, func() (bool, error) {

		c.client.DebugLogString(fmt.Sprintf("Waiting to destroy instance : %s", input.Name))

		info, err := c.GetServiceInstance(input)
		if err != nil {
			if client.WasNotFoundError(err) {
				// Service Instance could not be found, thus deleted
				return true, nil
			}
			// Some other error occurred trying to get instance, exit
			return false, err
		}

		c.client.DebugLogString(fmt.Sprintf("ServiceInstance [%s] waiting for deletion . Status : %s", info.ServiceId, info.Status))

		switch s := info.Status; s {
		case ServiceInstanceError:
			return false, fmt.Errorf("Error stopping instance: %s", info.ErrorReason)
		case ServiceInstanceTerminating:
			return false, nil
		default:
			return false, nil
		}
	})
}
