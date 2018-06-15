package java

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

const waitForServiceInstanceReadyPollInterval = 60 * time.Second
const waitForServiceInstanceReadyTimeout = 3600 * time.Second
const waitForServiceInstanceDeletePollInterval = 60 * time.Second
const waitForServiceInstanceDeleteTimeout = 3600 * time.Second
const deleteMaxRetries = 5

var (
	serviceInstanceContainerPath    = "/paas/api/v1.1/instancemgmt/%s/services/jaas/instances"
	serviceInstanceResourcePath     = "/paas/api/v1.1/instancemgmt/%s/services/jaas/instances/%s"
	serviceInstanceScaleUpDownPath  = "/hosts/scale"
	serviceInstanceDesiredStatePath = "/hosts/%s"
)

// ServiceInstanceClient is a client for the Service functions of the Java API.
type ServiceInstanceClient struct {
	ResourceClient
	PollInterval time.Duration
	Timeout      time.Duration
}

// ServiceInstanceClient obtains an ServiceInstanceClient which can be used to access to the
// Service Instance functions of the Java Cloud API
func (c *Client) ServiceInstanceClient() *ServiceInstanceClient {
	return &ServiceInstanceClient{
		ResourceClient: ResourceClient{
			Client:           c,
			ContainerPath:    serviceInstanceContainerPath,
			ResourceRootPath: serviceInstanceResourcePath,
		}}
}

// ServiceInstanceLevel specifies the level type for the service instance
type ServiceInstanceLevel string

const (
	// ServiceInstanceLevelPAAS - PAAS: Production-level service. This is the default. Supports Oracle Java Cloud Service instance creation
	// and monitoring, backup and restoration, patching, and scaling. Use PAAS if you want to enable domain partitions
	// using WebLogic Server 12.2.1, use AppToCloud artifacts to create a service instance, or create a service instance
	// for an Oracle Fusion Middleware product.
	ServiceInstanceLevelPAAS ServiceInstanceLevel = "PAAS"
	// ServiceInstanceLevelBasic - BASIC: Development-level service. Supports Oracle Java Cloud Service instance creation and monitoring
	// but does not support backup and restoration, patching, or scaling.
	ServiceInstanceLevelBasic ServiceInstanceLevel = "BASIC"
)

// ServiceInstanceBackupDestination specifies the backup destination type
type ServiceInstanceBackupDestination string

const (
	// ServiceInstanceBackupDestinationBoth - BOTH - Enable backups. This is the default. This means automated scheduled backups are enabled,
	// and on-demand backups can be initiated. All backups are stored on disk and the Oracle Storage
	// Cloud Service container that is specified in cloudStorageContainer.
	ServiceInstanceBackupDestinationBoth ServiceInstanceBackupDestination = "BOTH"
	// ServiceInstanceBackupDestinationNone  - NONE - Do not enable backups. This means automated scheduled backups are not enabled,
	// and on-demand backups cannot be initiated. When set to NONE, cloudStorageContainer is not required.
	ServiceInstanceBackupDestinationNone ServiceInstanceBackupDestination = "NONE"
)

// ServiceInstanceTargetDataSourceType specifies the different types for target data sources
type ServiceInstanceTargetDataSourceType string

const (
	// ServiceInstanceTargetDataSourceTypeGeneric - If the specified Database Cloud Service database deployment does not use Oracle RAC, the value must be Generic.
	ServiceInstanceTargetDataSourceTypeGeneric ServiceInstanceTargetDataSourceType = "Generic"
	// ServiceInstanceTargetDataSourceTypeMulti - If the specified Database Cloud Service database deployment uses Oracle RAC and the specified edition
	// (for WebLogic Server software) is EE, the value must be Multi.
	ServiceInstanceTargetDataSourceTypeMulti ServiceInstanceTargetDataSourceType = "Multi"
	// ServiceInstanceTargetDataSourceTypeGridLink - If the specified Database Cloud Service database deployment uses Oracle RAC and the specified edition
	// (for WebLogic Server software) is SUITE, the value can be GridLink or Multi.
	ServiceInstanceTargetDataSourceTypeGridLink ServiceInstanceTargetDataSourceType = "GridLink"
)

// ServiceInstanceDomainMode specifies the differnt domain modes a service instance can be in
type ServiceInstanceDomainMode string

const (
	// ServiceInstanceDomainModeDev - DEVELOPMENT
	ServiceInstanceDomainModeDev ServiceInstanceDomainMode = "DEVELOPMENT"
	// ServiceInstanceDomainModePro - PRODUCTION
	ServiceInstanceDomainModePro ServiceInstanceDomainMode = "PRODUCTION"
)

// ServiceInstanceEdition specifies the different editions a service instance can be
type ServiceInstanceEdition string

const (
	// ServiceInstanceEditionSE  - SE
	ServiceInstanceEditionSE ServiceInstanceEdition = "SE"
	// ServiceInstanceEditionEE - EE
	ServiceInstanceEditionEE ServiceInstanceEdition = "EE"
	// ServiceInstanceEditionSuite - SUITE
	ServiceInstanceEditionSuite ServiceInstanceEdition = "SUITE"
)

// ServiceInstanceLoadBalancingPolicy specifies the different load balancing policies a load balancer can use
type ServiceInstanceLoadBalancingPolicy string

const (
	// ServiceInstanceLoadBalancingPolicyLCC - LEAST_CONNECTION_COUNT
	ServiceInstanceLoadBalancingPolicyLCC ServiceInstanceLoadBalancingPolicy = "LEAST_CONNECTION_COUNT"
	// ServiceInstanceLoadBalancingPolicyLRT - LEAST_RESPONSE_TIME
	ServiceInstanceLoadBalancingPolicyLRT ServiceInstanceLoadBalancingPolicy = "LEAST_RESPONSE_TIME"
	// ServiceInstanceLoadBalancingPolicyRR - ROUND_ROBIN
	ServiceInstanceLoadBalancingPolicyRR ServiceInstanceLoadBalancingPolicy = "ROUND_ROBIN"
)

// ServiceInstanceShape specifies the shapes a service instance can be
type ServiceInstanceShape string

const (
	// Suportted OCI Classic Shapes

	// ServiceInstanceShapeOC3 - oc3: 1 OCPU, 7.5 GB memory
	ServiceInstanceShapeOC3 ServiceInstanceShape = "oc3"
	// ServiceInstanceShapeOC4 - oc4: 2 OCPUs, 15 GB memory
	ServiceInstanceShapeOC4 ServiceInstanceShape = "oc4"
	// ServiceInstanceShapeOC5 - oc5: 4 OCPUs, 30 GB memory
	ServiceInstanceShapeOC5 ServiceInstanceShape = "oc5"
	// ServiceInstanceShapeOC6 - oc6: 8 OCPUs, 60 GB memory
	ServiceInstanceShapeOC6 ServiceInstanceShape = "oc6"
	// ServiceInstanceShapeOC7 - oc7: 16 OCPUS, 120 GB memory
	ServiceInstanceShapeOC7 ServiceInstanceShape = "oc7"
	// ServiceInstanceShapeOC1M - oc1m: 1 OCPU, 15 GB memory
	ServiceInstanceShapeOC1M ServiceInstanceShape = "oc1m"
	// ServiceInstanceShapeOC2M - oc2m: 2 OCPUs, 30 GB memory
	ServiceInstanceShapeOC2M ServiceInstanceShape = "oc2m"
	// ServiceInstanceShapeOC3M - oc3m: 4 OCPUs, 60 GB memory
	ServiceInstanceShapeOC3M ServiceInstanceShape = "oc3m"
	// ServiceInstanceShapeOC4M - oc4m: 8 OCPUs, 120 GB memory
	ServiceInstanceShapeOC4M ServiceInstanceShape = "oc4m"
	// ServiceInstanceShapeOC5M - oc5m: 16 OCPUS, 240 GB memory
	ServiceInstanceShapeOC5M ServiceInstanceShape = "oc5m"

	// Supported OCI VM shapes

	// ServiceInstanceShapeVMStandard1_1 - VM.Standard1.1: 1 OCPU, 7 GB memory
	ServiceInstanceShapeVMStandard1_1 ServiceInstanceShape = "VM.Standard1.1"
	// ServiceInstanceShapeVMStandard1_2 - VM.Standard1.2: 2 OCPU, 14 GB memory
	ServiceInstanceShapeVMStandard1_2 ServiceInstanceShape = "VM.Standard1.2"
	// ServiceInstanceShapeVMStandard1_4 - VM.Standard1.4: 4 OCPU, 28 GB memory
	ServiceInstanceShapeVMStandard1_4 ServiceInstanceShape = "VM.Standard1.4"
	// ServiceInstanceShapeVMStandard1_8 - VM.Standard1.8: 8 OCPU, 56 GB memory
	ServiceInstanceShapeVMStandard1_8 ServiceInstanceShape = "VM.Standard1.8"
	// ServiceInstanceShapeVMStandard1_16 - VM.Standard1.16: 16 OCPU, 112 GB memory
	ServiceInstanceShapeVMStandard1_16 ServiceInstanceShape = "VM.Standard1.16"
	// ServiceInstanceShapeVMStandard2_1 - VM.Standard2.1: 1 OCPU, 15 GB memory
	ServiceInstanceShapeVMStandard2_1 ServiceInstanceShape = "VM.Standard2.1"
	// ServiceInstanceShapeVMStandard2_2 -  VM.Standard2.2: 2 OCPU, 30 GB memory
	ServiceInstanceShapeVMStandard2_2 ServiceInstanceShape = "VM.Standard2.2"
	// ServiceInstanceShapeVMStandard2_4 - VM.Standard2.4: 4 OCPU, 60 GB memory
	ServiceInstanceShapeVMStandard2_4 ServiceInstanceShape = "VM.Standard2.4"
	// ServiceInstanceShapeVMStandard2_8 - VM.Standard2.8: 8 OCPU, 120 GB memory
	ServiceInstanceShapeVMStandard2_8 ServiceInstanceShape = "VM.Standard2.8"
	// ServiceInstanceShapeVMStandard2_16 - VM.Standard2.16: 16 OCPU, 240 GB memory
	ServiceInstanceShapeVMStandard2_16 ServiceInstanceShape = "VM.Standard2.16"
	// ServiceInstanceShapeVMStandard2_24 - VM.Standard2.24: 24 OCPU, 320 GB memory
	ServiceInstanceShapeVMStandard2_24 ServiceInstanceShape = "VM.Standard2.24"

	// Supported OCI Bare Metal shapes

	// ServiceInstanceShapeBMStandard1_36 - BM.Standard1.36: 36 OCPU, 256 GB memory
	ServiceInstanceShapeBMStandard1_36 ServiceInstanceShape = "BM.Standard1.36"
	// ServiceInstanceShapeBMStandard2_52 - BM.Standard2.52: 52 OCPU, 768 GB memory
	ServiceInstanceShapeBMStandard2_52 ServiceInstanceShape = "BM.Standard2.52"
)

// ServiceInstanceType specifies the different types of service instances
type ServiceInstanceType string

const (
	// ServiceInstanceTypeWebLogic - weblogic
	ServiceInstanceTypeWebLogic ServiceInstanceType = "weblogic"
	// ServiceInstanceTypeDataGrid - datagrid
	ServiceInstanceTypeDataGrid ServiceInstanceType = "datagrid"
	// ServiceInstanceTypeOTD - otd
	ServiceInstanceTypeOTD ServiceInstanceType = "otd"
)

// ServiceInstanceUpperStackProductName specifies the different upperstack product names
type ServiceInstanceUpperStackProductName string

const (
	// ServiceInstanceUpperStackProductNameODI  - ODI
	ServiceInstanceUpperStackProductNameODI ServiceInstanceUpperStackProductName = "ODI"
	// ServiceInstanceUpperStackProductNameWCP - WCP
	ServiceInstanceUpperStackProductNameWCP ServiceInstanceUpperStackProductName = "WCP"
)

// ServiceInstanceVersion specifies the different type of versions a service instance can be
type ServiceInstanceVersion string

const (
	// ServiceInstanceVersion1221 - 12.2.1
	ServiceInstanceVersion1221 ServiceInstanceVersion = "12.2.1"
	// ServiceInstanceVersion1213 - 12.1.3
	ServiceInstanceVersion1213 ServiceInstanceVersion = "12.1.3"
	// ServiceInstanceVersion1036 - 10.3.6
	ServiceInstanceVersion1036 ServiceInstanceVersion = "10.3.6"
)

// ServiceInstanceSubscriptionType specifies the different types of subscriptions
type ServiceInstanceSubscriptionType string

const (
	// ServiceInstanceSubscriptionTypeHourly - HOURLY
	ServiceInstanceSubscriptionTypeHourly ServiceInstanceSubscriptionType = "HOURLY"
	// ServiceInstanceSubscriptionTypeMonthly - MONTHLY
	ServiceInstanceSubscriptionTypeMonthly ServiceInstanceSubscriptionType = "MONTHLY"
)

// ServiceInstanceServiceComponentType specifies the different types a component can be
type ServiceInstanceServiceComponentType string

const (
	// ServiceInstanceServiceComponentTypeJDK - JDK
	ServiceInstanceServiceComponentTypeJDK ServiceInstanceServiceComponentType = "JDK"
	// ServiceInstanceServiceComponentTypeOTD - OTD
	ServiceInstanceServiceComponentTypeOTD ServiceInstanceServiceComponentType = "OTD"
	// ServiceInstanceServiceComponentTypeOTDJDK - OTD_JDK
	ServiceInstanceServiceComponentTypeOTDJDK ServiceInstanceServiceComponentType = "OTD_JDK"
	// ServiceInstanceServiceComponentTypeWLS - WLS
	ServiceInstanceServiceComponentTypeWLS ServiceInstanceServiceComponentType = "WLS"
)

// ServiceInstanceServiceComponentVersion specifies the different versions a component can be
type ServiceInstanceServiceComponentVersion string

const (
	// ServiceInstanceServiceComponentVersionWLS - 12.1.3.0.5
	ServiceInstanceServiceComponentVersionWLS ServiceInstanceServiceComponentVersion = "12.1.3.0.5"
	// ServiceInstanceServiceComponentVersionOTD - 11.1.1.9.1
	ServiceInstanceServiceComponentVersionOTD ServiceInstanceServiceComponentVersion = "11.1.1.9.1"
	// ServiceInstanceServiceComponentVersionJDK - 1.7.0_91
	ServiceInstanceServiceComponentVersionJDK ServiceInstanceServiceComponentVersion = "1.7.0_91"
	// ServiceInstanceServiceComponentVersionOTDJDK - 1.7.0_91
	ServiceInstanceServiceComponentVersionOTDJDK ServiceInstanceServiceComponentVersion = "1.7.0_91"
)

// ServiceInstanceShiftStatus specifies the different statuses a shift can be in
type ServiceInstanceShiftStatus string

const (
	// ServiceInstanceShiftStatusReady - readyToShift
	ServiceInstanceShiftStatusReady ServiceInstanceShiftStatus = "readyToShift"
	// ServiceInstanceShiftStatusCompleted - shiftCompleted
	ServiceInstanceShiftStatusCompleted ServiceInstanceShiftStatus = "shiftCompleted"
	// ServiceInstanceShiftStatusFailed - shiftFailed
	ServiceInstanceShiftStatusFailed ServiceInstanceShiftStatus = "shiftFailed"
)

// ServiceInstanceStatus specifies the different status a service instance can be in
type ServiceInstanceStatus string

const (
	// ServiceInstanceStatusNew - NEW
	ServiceInstanceStatusNew ServiceInstanceStatus = "NEW"
	// ServiceInstanceStatusInitializing - INTIALIZING
	ServiceInstanceStatusInitializing ServiceInstanceStatus = "INITIALIZING"
	// ServiceInstanceStatusReady - READY
	ServiceInstanceStatusReady ServiceInstanceStatus = "READY"
	// ServiceInstanceStatusConfiguring - CONFIGURING
	ServiceInstanceStatusConfiguring ServiceInstanceStatus = "CONFIGURING"
	// ServiceInstanceStatusTerminating - TERMINATING
	ServiceInstanceStatusTerminating ServiceInstanceStatus = "TERMINATING"
	// ServiceInstanceStatusStopping - STOPPING
	ServiceInstanceStatusStopping ServiceInstanceStatus = "STOPPING"
	// ServiceInstanceStatusStopped - STOPPED
	ServiceInstanceStatusStopped ServiceInstanceStatus = "STOPPED"
	// ServiceInstanceStatusStarting - STARTING
	ServiceInstanceStatusStarting ServiceInstanceStatus = "STARTING"
	// ServiceInstanceStatusDisabling - DISABLING
	ServiceInstanceStatusDisabling ServiceInstanceStatus = "DISABLING"
	// ServiceInstanceStatusDisabled - DISABLED
	ServiceInstanceStatusDisabled ServiceInstanceStatus = "DISABLED"
	// ServiceInstanceStatusTerminated - TERMINATED
	ServiceInstanceStatusTerminated ServiceInstanceStatus = "TERMINATED"
)

// ServiceInstanceClusterType are the constances around cluster types for a service instance
type ServiceInstanceClusterType string

const (
	// ServiceInstanceClusterTypeApplication - APPLICATION_CLUSTER - Application cluster (default).
	// This is the WebLogic cluster that will run the service applications, which are accessible via the
	// local load balancer resource (OTD) or the Oracle managed load balancer.
	ServiceInstanceClusterTypeApplication ServiceInstanceClusterType = "APPLICATION_CLUSTER"
	// ServiceInstanceClusterTypeCaching  - CACHING_CLUSTER - Caching (data grid) cluster. This is the WebLogic cluster for Coherence storage.
	ServiceInstanceClusterTypeCaching ServiceInstanceClusterType = "CACHING_CLUSTER"
)

// ServiceInstanceCustomPayloadType are the constants for payload type
type ServiceInstanceCustomPayloadType string

const (
	// ServiceInstanceCustomPayloadTypeApp2Cloud - app2cloud
	ServiceInstanceCustomPayloadTypeApp2Cloud ServiceInstanceCustomPayloadType = "app2cloud"
)

// ServiceInstanceActivityStatus are the constants for the different statuses a service instance can be in
type ServiceInstanceActivityStatus string

const (
	// ServiceInstanceActivityStatusNew - NEW
	ServiceInstanceActivityStatusNew ServiceInstanceActivityStatus = "NEW"
	// ServiceInstanceActivityStatusRunning - RUNNING
	ServiceInstanceActivityStatusRunning ServiceInstanceActivityStatus = "RUNNING"
	// ServiceInstanceActivityStatusSucceed - SUCCEED
	ServiceInstanceActivityStatusSucceed ServiceInstanceActivityStatus = "SUCCEED"
	// ServiceInstanceActivityStatusFailed - FAILED
	ServiceInstanceActivityStatusFailed ServiceInstanceActivityStatus = "FAILED"
	// ServiceInstanceActivityStatusInitializing - INITIALIZING
	ServiceInstanceActivityStatusInitializing ServiceInstanceActivityStatus = "INITIALIZING"
	// ServiceInstanceActivityStatusConfiguring - CONFIGURING
	ServiceInstanceActivityStatusConfiguring ServiceInstanceActivityStatus = "CONFIGURING"
	// ServiceInstanceActivityStatusTerminating - TERMINATING
	ServiceInstanceActivityStatusTerminating ServiceInstanceActivityStatus = "TERMINATING"
	// ServiceInstanceActivityStatusStopping - STOPPING
	ServiceInstanceActivityStatusStopping ServiceInstanceActivityStatus = "STOPPING"
	// ServiceInstanceActivityStatusStopped - STOPPED
	ServiceInstanceActivityStatusStopped ServiceInstanceActivityStatus = "STOPPED"
	// ServiceInstanceActivityStatusStarting - STARTING
	ServiceInstanceActivityStatusStarting ServiceInstanceActivityStatus = "STARTING"
	// ServiceInstanceActivityStatusDisabling - DISABLING
	ServiceInstanceActivityStatusDisabling ServiceInstanceActivityStatus = "DISABLING"
	// ServiceInstanceActivityStatusDisabled - DISABLED
	ServiceInstanceActivityStatusDisabled ServiceInstanceActivityStatus = "DISABLED"
	// ServiceInstanceActivityStatusTerminated - TERMINATED
	ServiceInstanceActivityStatusTerminated ServiceInstanceActivityStatus = "TERMINATED"
)

// ServiceInstanceLifecycleState defines the constants for the lifecycle state
type ServiceInstanceLifecycleState string

const (
	// ServiceInstanceLifecycleStateStop - stop: Stops the Database Cloud Service instance or compute node.
	ServiceInstanceLifecycleStateStop ServiceInstanceLifecycleState = "stop"
	// ServiceInstanceLifecycleStateStart - start: Starts the Database Cloud Service instance or compute node.
	ServiceInstanceLifecycleStateStart ServiceInstanceLifecycleState = "start"
	// ServiceInstanceLifecycleStateRestart - restart: Restarts the Database Cloud Service instance or compute node.
	ServiceInstanceLifecycleStateRestart ServiceInstanceLifecycleState = "restart"
)

// ServiceInstance specifies the attributes associated with a service instance
type ServiceInstance struct {
	// Activity logs for the service instance.
	ActivityLogs []ActivityLog `json:"activityLogs"`
	// Host name of the Administration Server for this service instance.
	AdminHostName string `json:"adminHostName"`
	// Information about service instance attributes.
	Attributes Attributes `json:"attributes"`
	// Information about service instance backup operations.
	Backup Backup `json:"backup"`
	// Flag that specifies whether this Oracle Java Cloud Service instance is a clone of an existing service instance.
	Clone bool `json:"clone"`
	// Groups details about the WLS component and the OTD component (if provisioned).
	Components Components `json:"components"`
	// Location where the service instance is provisioned
	ComputeSiteName string `json:"computeSiteName"`
	// Date and time the Oracle Java Cloud Service instance was created.
	CreationDate string `json:"creationDate"`
	// Name of the user account that was used to create the Oracle Java Cloud Service instance.
	Creator string `json:"creator"`
	// Identity domain ID for the Oracle Java Cloud Service account.
	DomainName string `json:"domainName"`
	// The Oracle WebLogic Server software edition that was provisioned on this service instance. For example, SE, EE, or SUITE.
	Edition ServiceInstanceEdition `json:"edition"`
	// Display name of the WebLogic Server software edition. For example, Enterprise Edition.
	EditionDisplayName string `json:"editionDisplayName"`
	// The URL of the Fusion Middleware Control console.
	FMWRoot string `json:"FMW_ROOT"`
	// Components key to the operation of this service instance.
	KeyComponentInstance string `json:"keyComponentInstance"`
	// Current version for the service definition (schema) used by this service instance.
	MetaVersion string `json:"metaVersion"`
	// Metering frequency. For example: HOURLY or MONTHLY
	MeteringFrequency ServiceInstanceSubscriptionType `json:"meteringFrequency"`
	// Display name of metering frequency.
	MeteringFrequencyDisplayName string `json:"meteringFrequencyDisplayName"`
	// The URL of the OTD console
	OTDRoot string `json:"OTD_ROOT"`
	// Patching information related to this service instance
	Patching Patching `json:"patching"`
	// Region where the service instance is provisioned
	Region string `json:"region"`
	// The specific Oracle WebLogic Server software binaries in use. For example: 12.2.1.2.x, 12.1.3.0.x or 10.3.6.0.x.
	// Note: This value is updated when the service instance is patched.
	ReleaseVersion string `json:"releaseVersion"`
	// Free-form text that was provided about this service instance when it was created.
	ServiceDescription string `json:"serviceDescription"`
	// ID of the Oracle Java Cloud Service Instance
	ServiceID int `json:"serviceId"`
	// Service Level
	ServiceLevel ServiceInstanceLevel `json:"serviceLevel"`
	// Display name of the service level
	ServiceLevelDisplayName string `json:"serviceLevelDisplayName"`
	// Name given to this service instance when it was created
	ServiceName string `json:"serviceName"`
	// State of the service instance
	ServiceStateDisplayName string `json:"serviceStateDisplayName"`
	// Type of this service instance
	ServiceType string `json:"serviceType"`
	// The Oracle WebLogic Server software release that was provisioned on this service instance
	ServiceVersion string `json:"serviceVersion"`
	// Current state of this Oracle Java Cloud Service Instance
	State ServiceInstanceStatus `json:"state"`
	// Subscription type
	Subscription ServiceInstanceSubscriptionType `json:"subscription"`
	// SSD Storage
	TotalSSDStorage int `json:"totalSSDStorage"`
	// The URL of the WebLogic Server Administration console
	WLSRoot string `json:"WLS_ROOT"`
}

// ActivityLog specifies the acitivty log information around a service instance
type ActivityLog struct {
	// ID of the activity log.
	ActivityLogID int `json:"activityLogId"`
	// Date and time the operation ended.
	EndDate string `json:"endDate"`
	// Name of the identity domain for the Oracle Java Cloud Service account.
	IdentityDomain string `json:"identityDomain"`
	// Name of the user who initiated the operation.
	InitiatedBy string `json:"initiatedBy"`
	// Job ID for the operation.
	// Note: This value may be set to No Job Submitted temporarily, just prior to being submitted for processing.
	JobID int `json:"jobId"`
	// Messages related to the activity.
	Messages []Message `json:"messages"`
	// ID of the operation.
	OperationID int `json:"operationId"`
	// Operation type. For example: RESTORE, BACKUP, START_SERVICE, STOP_SERVICE, and so on
	OperationType string `json:"operationType"`
	// ID of the Oracle Java Cloud Service instance.
	ServiceID int `json:"serviceId"`
	// Name given to this service instance when it was created.
	ServiceName string `json:"serviceName"`
	// Type of this service instance. For Oracle Java Cloud Service instances, the value is JaaS.
	ServiceType string `json:"serviceType"`
	// Date and time the operation started.
	StartDate string `json:"startDate"`
	// Final status of the operation. Example status messages include: NEW, RUNNING, SUCCEED, and FAILED
	Status ServiceInstanceActivityStatus `json:"status"`
	// Summary of the activity.
	SummaryMessage string `json:"summaryMessage"`
}

// Message specifies the message information associated with a service instance
type Message struct {
	// Date and time the activity was logged.
	ActivityDate string `json:"activityDate"`
	// Details of the activity
	Message string `json:"message"`
}

// Attributes specifies the attributes associated with a service instance
type Attributes struct {
	// Service instance backup details.
	BackupDestination AttributeInfo `json:"BACKUP_DESTINATION"`
	// Object storage container details.
	CloudStorageContainer AttributeInfo `json:"cloudStorageContainer"`
	// Application content details
	ContentRoot AttributeInfo `json:"CONTENT_ROOT"`
	// AppToCloud payload
	CustomPayload AttributeInfo `json:"customPayload"`
	// Fusion Middleware Control Console details.
	FMWRoot AttributeInfo `json:"FMW_ROOT"`
	// JDK details
	JDKVersion AttributeInfo `json:"jdkVersion"`
	// Local load balancer details
	OTDRoot AttributeInfo `json:"OTD_ROOT"`
	// Sample application details
	SampleRoot AttributeInfo `json:"SAMPLE_ROOT"`
	// WebLogic Server Console details
	WLSRoot AttributeInfo `json:"WLS_ROOT"`
}

// AttributeInfo specifies the attribute information associated with the service instance
type AttributeInfo struct {
	// Attribute label.
	DisplayName string `json:"displayName"`
	// Attribute display value
	DisplayValue string `json:"displayValue"`
	// Attribute is key binding
	IsKeyBinding bool `json:"isKeyBinding"`
	// Type of the attribute value.
	Type string `json:"type"`
	// Value of Attribute
	Value string `json:"value"`
}

// Backup speicifes the backup information about a service instance
type Backup struct {
	// The date and the time of the last successful backup operation.
	LastBackupDate string `json:"lastBackupDate"`
	// The date and the time of the last failed backup operation.
	LastFailedBackupDate string `json:"lastFailedBackupDate"`
}

// Components specifies the information about the components associated with the service instnace
type Components struct {
	// Details about the OTD component.
	OTD OTD `json:"OTD"`
	// Details about the WLS component.
	WLS WLS `json:"WLS"`
}

// OTD specifies information about the oracle traffic director associated with the service instance
type OTD struct {
	// Host name of the administration server.
	AdminHostName string `json:"adminHostName"`
	// OTD component attribute details
	Attributes OTDAttributes `json:"attributes"`
	// ID of this component of this service instance.
	ComponentID string `json:"componentId"`
	// Type of the component.
	ComponentType string `json:"componentType"`
	// Date and time this component was created.
	CreationDate string `json:"creationDate"`
	// OTD Display Name
	DisplayName string `json:"displayName"`
	// Groups details about hosts that are running in the service instance for the OTD component.
	Hosts Hosts `json:"hosts"`
	// Name of this componenet instance.
	InstanceName string `json:"instanceName"`
	// Role of Instance
	InstanceRole string `json:"instanceRole"`
	// ID of the Oracle Java Cloud Service Instance
	ServiceID int `json:"serviceId"`
	// State of the component
	State ServiceInstanceStatus `json:"state"`
	// Oracle Traffic Director software version
	Version string `json:"version"`
	// Groups details about OTD VM instances by host name. Each VM instance is a JSON object element.
	// This object will be deprecated in the near future. The properties host and userHosts contain the same
	// information as vmInstances, and more.
	VMInstances VMInstances `json:"vmInstances"`
}

// WLS sepcifies the information about the weblogic server associated with the service instance
type WLS struct {
	// Host name of the administration server.
	AdminHostName string `json:"adminHostName"`
	// WLS Component Attribute details
	Attributes WLSAttributes `json:"attributes"`
	// Groups details about clusters in the domain by cluster name. Each cluster is a JSON object element.
	// Clusters have dynamic JSON keys that need to be accounted for
	Clusters map[string]Clusters `json:"clusters"`
	// ID of this component of this service instance.
	ComponentID int `json:"componentId"`
	// Type of the component
	ComponentType string `json:"componentType"`
	// Date and time the component was created
	CreationDate string `json:"creationDate"`
	// Display name
	DisplayName string `json:"displayName"`
	// Groups details about all hosts that are running in the service instance for the WLS component.
	Hosts Hosts `json:"hosts"`
	// Instance name
	InstanceName string `json:"instanceName"`
	// Instance role
	InstanceRole string `json:"instanceRole"`
	// ID of the Oracle Java Cloud Service Instance
	ServiceID int `json:"serviceId"`
	// State of the component
	State ServiceInstanceStatus `json:"state"`
	// Oracle WebLogic Server software version
	Version string `json:"version"`
	// Groups details about WLS VM instances by host name. Each VM instance is a JSON object element.
	VMInstances map[string]HostName `json:"vmInstances"`
}

// Clusters specifies the information about the clusters associated with the service instance
type Clusters struct {
	ClusterID    int                    `json:"clusterId"`
	ClusterName  string                 `json:"clusterName"`
	ClusterType  string                 `json:"clusterType"`
	CreationDate string                 `json:"creationDate"`
	PaaSServers  map[string]PaaSServers `json:"paasServers"`
	// Identifies service specific cluster and server information. Details include:
	// Cluster type - APPLICATION_CLUSTER or CACHING_CLUSTER
	// Shape - Compute shape used by nodes of this cluster
	// Whether this cluster is accessible from the public network (external value is true if accessible)
	Profile string `json:"profile"`
}

// PaaSServers specifies the informaiton about the different paas servers associated with the service instance
type PaaSServers struct {
	// Attribute details of a specific server.
	Attributes PaaSAttributes `json:"attributes"`
}

// PaaSAttributes specifies the platform as a service attributes associated with the service instance
type PaaSAttributes struct {
	// One or more Managed Server JVM arguments separated by a space.
	AdditionalJVMArgs string `json:"additional_jvm_args"`
	// Name of the Coherence cluster (the system-level CoherenceClusterSystemResource) that the WLS cluster
	// (application or caching) is associated with. Default value is DataGridConfig.
	CCSR string `json:"ccsr"`
	// Cluster that this Managed Server is a member of.
	Cluster string `json:"cluster"`
	// Java heap size of the Managed Server JVM.
	HeapSize string `json:"heap_size"`
	// Initial heap size
	HeapStart string `json:"heap_start"`
	// Maximum Permanent Generation (PermGen) space in Java heap memory for a Managed Server JVM.
	MaxPermSize string `json:"max_perm_size"`
	// Initial Permanent Generation (PermGen) space in Java heap memory for a Managed Server JVM.
	PermSize string `json:"perm_size"`
	// Port number
	Port string `json:"port"`
	// Server role
	Role string `json:"role"`
	// Server type
	ServerType string `json:"server_type"`
	// SSL port number
	SSLPort string `json:"ssl_port"`
	// Cluster template for the server. For example: ExampleI_cluster_Template, DataGridServer-Template
	Template string `json:"template"`
}

// WLSAttributes specifies information about the weblogic server associated with the service instance
type WLSAttributes struct {
	// Attributes of the Fusion Middleware (Upper Stack) product to be installed or installed on the service instance.
	UpperStackProductName AttributeInfo `json:"upperStackProductName"`
}

// OTDAttributes specifies information about the Oracle Traffic Director associated with the service instance
type OTDAttributes struct {
	AdminPort                    AttributeInfo `json:"ADMIN_PORT"`
	ListenerPort                 AttributeInfo `json:"LISTENER_PORT"`
	ListenerPortEnabled          AttributeInfo `json:"LISTENER_PORT_ENABLED"`
	PrivilgedListenerPort        AttributeInfo `json:"PRIV_LISTENER_PORT"`
	PrivilegedSecureListenerPort AttributeInfo `json:"PRIV_SECURE_LISTENER_PORT"`
	SecuredListenerPort          AttributeInfo `json:"SECURE_LISTENER_PORT"`
}

// Hosts specifies information about the different hosts on the service instance
type Hosts struct {
	UserHosts UserHosts `json:"userHosts"`
}

// UserHosts specifies information about the user hosts on a service instance
type UserHosts struct {
	// Host names have dynamic JSON keys that need to be accounted for
	HostName map[string]HostName `json:"host-name"`
}

// HostName specifies information about the hostname on the service instances
type HostName struct {
	// Type of component.
	ComponentType string `json:"componentType"`
	// Creation Date
	CreationDate string `json:"creationDate"`
	// DNS host name
	HostName string `json:"hostName"`
	// Id of this host
	ID int `json:"id"`
	// IP address of this host. Not all hosts are accessible from the public Internet.
	IPAddress string `json:"ipAddress"`
	// true if this host contains an administration server, false otherwise.
	IsAdminNode bool `json:"isAdminNode"`
	// Label associated with host
	Label string `json:"label"`
	// Public Accessible IP Address
	PublicIPAddress string `json:"publicIpAddress"`
	// Similar to `usageType`
	Role string `json:"role"`
	// The compute shape of the host
	ShapeID string `json:"shapeId"`
	// State of the host.
	State ServiceInstanceStatus `json:"state"`
	// Total megabytes of block storage used by this host.
	TotalStorage int `json:"totalStorage"`
	// Purpose of this host
	UsageType string `json:"usageType"`
	// Unique identifier fo this host
	UUID string `json:"uuid"`
	// Id of this host
	VMID string `json:"vmId"`
}

// VMInstances specifies information about the vm instances on the service instance
type VMInstances struct {
	// Details about a specific VM instance.
	// TODO HostName will be deprecated in the near future
	// VMOTDs have dyanmic JSON keys that needs to be accounted for
	VMOTD map[string]HostName
}

// Patching specifies information about the patches for the service instance
type Patching struct {
	// Current Operation
	CurrentOperation CurrentOperation `json:"currentOperation"`
	// The total number of patches available for this service instance.
	TotalAvailablePatches int `json:"totalAvailablePatches"`
}

// CurrentOperation specifies the information about the operation currently working on the service instance
type CurrentOperation struct {
	// Details about the current patching operation.
	Operation string `json:"operation"`
}

// IPReservation specifies the information about the ip reservation associated with the service instance
type IPReservation struct {
	// Name of an IP reservation that is assigned to a node on the service instance.
	Name string `json:"name"`
}

// CreateServiceInstanceInput specifies the attributes of the service instance that will be created
type CreateServiceInstanceInput struct {
	// This attribute is only applicable when provisioning an Oracle Java Cloud Service
	// instance in a region on Oracle Cloud Infrastructure Classic, and a custom IP
	// network is specified in ipNetwork. Flag that specifies whether to assign (true)
	// or not assign (false) public IP addresses to the nodes in your service instance.
	// Optional.
	AssignPublicIP bool `json:"assignPublicIP,omitempty"`
	// This attribute is available only on Oracle Cloud Infrastructure. It is required along with region and subnet.
	// Name of a data center location in the Oracle Cloud Infrastructure region that is specified in region.
	// A region is a localized geographic area, composed of one or more availability domains (data centers).
	// The availability domain value format is an account-specific prefix followed by <region>-<ad>.
	// For example, FQCn:US-ASHBURN-AD1 where FQCn is the account-specific prefix.
	// The Oracle Database Cloud Service database deployment on Oracle Cloud Infrastructure must be in the same
	// region and virtual cloud network as the Oracle Java Cloud Service instance you are creating on Oracle Cloud Infrastructure.
	// The service instances do not need to be on the same subnet or availability domain.
	// See Regions and Availability Domains in Oracle Cloud Infrastructure Services.
	// Optional.
	AvailabilityDomain string `json:"availabilityDomain,omitempty"`
	// This attribute is applicable only when serviceLevel is set to PAAS.
	// Specifies whether to enable backups for this Oracle Java Cloud Service instance.
	// Optional.
	BackupDestination ServiceInstanceBackupDestination `json:"backupDestination,omitempty"`
	// URI of the object storage container or bucket for storing Oracle Java Cloud Service instance backups.
	// This attribute is not required if backupDestination is set to NONE. It is also not required when provisioning an
	// Oracle Java Cloud Service service instance with the BASIC service level.
	// Note:
	// Do not use a container or bucket that you use to back up Oracle Java Cloud Service instances for any other purpose.
	// For example, do not also use the same container or bucket to back up Oracle Database Cloud Service database deployments.
	// Using one container or bucket for multiple purposes can result in billing errors.
	// You do not have to specify a container or bucket if you provision the service instance without enabling backups.
	// On Oracle Cloud Infrastructure Classic, the object storage container does not have to be created ahead of provisioning
	// your Oracle Java Cloud Service instance.
	// To specify the container (existing or new), use one of the following formats:
	// Storage-<identitydomainid>/<containername>
	// <storageservicename>-<identitydomainid>/<containername>
	// https://foo.storage.oraclecloud.com/v1/MyService-bar/MyContainer
	// The format to use to specify the container name depends on the URL of your Oracle Cloud Infrastructure Object Storage
	// Classic account. To identify the URL of your account, see Finding the REST Endpoint URL for Your Service Instance in
	// Using Oracle Cloud Infrastructure Object Storage Classic.
	// On Oracle Cloud Infrastructure, the object storage bucket must be created ahead of provisioning your Oracle Java Cloud
	// Service instance. Do not use the same bucket for each service instance. Certain prerequisites must be satisfied when
	// you create the bucket. See Prerequisites for Oracle Platform Services on Oracle Cloud Infrastructure in Oracle Cloud
	// Infrastructure Services. Then use the following URL form to specify the bucket:
	// https://swiftobjectstorage.<region>.oraclecloud.com/v1/<account>/<container>
	// For example:
	// https://swiftobjectstorage.us-phoenix-1.oraclecloud.com/v1/acme/mycontainer
	// Optional.
	CloudStorageContainer string `json:"cloudStorageContainer,omitempty"`
	// On Oracle Cloud Infrastructure Classic, this is the password for the Oracle Cloud Infrastructure Object Storage
	// Classic user who has read and write access to the container that is specified in cloudStorageContainer. This attribute
	// is not required for the BASIC service level.
	// On Oracle Cloud Infrastructure, this is the Swift password to use with the Object Storage service.
	// Optional.
	CloudStoragePassword string `json:"cloudStoragePassword,omitempty"`
	// User name for the object storage user. The user name must be specified if cloudStorageContainer is set.
	// On Oracle Cloud Infrastructure Classic, this is the user name for the Oracle Cloud Infrastructure Object Storage
	// Classic user who has read and write access to the container that is specified in cloudStorageContainer.
	// This attribute is not required for the BASIC service level.
	// On Oracle Cloud Infrastructure, this is the user name for the Object Storage service user.
	// Optional.
	CloudStorageUsername string `json:"cloudStorageUser,omitempty"`
	// Groups properties for the Oracle WebLogic Server component (WLS) and the optional Oracle Traffice Director (OTD) component.
	// Optional
	Components CreateComponents `json:"components,omitempty"`
	// Software edition for Oracle WebLogic Server. Valid values include:
	// SE - Standard edition. See Oracle WebLogic Server Standard Edition. Do not use the Standard edition if you are enabling domain partitions using WebLogic Server 12.2.1, or using upperStackProductName to provision a service instance for an Oracle Fusion Middleware product. Scaling a cluster is also not supported on service instances that are based on the Standard edition.
	// EE - Enterprise Edition. This is the default for both PAAS and BASIC service levels. See Oracle WebLogic Server Enterprise Edition.
	// SUITE - Suite edition. See Oracle WebLogic Suite.
	//Optional
	Edition ServiceInstanceEdition `json:"edition,omitempty"`
	// This attribute is applicable only to accounts where regions are supported.
	// This attribute is not applicable when provisioning Oracle Java Cloud Service instances in Oracle Cloud Infrastructure.
	// The three-part name of a custom IP network to attach this service instance to. For example:
	// /Compute-identity_domain/user/object
	// A region name must be specified in order to use ipNetwork. Only those IP networks already created in the specified
	// Oracle Cloud Infrastructure Compute Classic region can be used.
	// If you specify an IP network, the dbServiceName for this service instance must also be attached to an ipNetwork.
	// If this service instance and the database deployment are attached to different IP networks, the two IP networks
	// must be connected to the same IP network exchange.
	// See Creating an IP Network in Using Oracle Cloud Infrastructure Compute Classic.
	// A consequence of using an IP network is that the auto-assigned IP address could change each time the service
	// instance is started. To assign fixed public IP addresses to a service instance that is attached to an IP network,
	// you can first create reserved IP addresses, then provision the service instance to use those persistent IP addresses.
	// Optional.
	IPNetwork string `json:"ipNetwork,omitempty"`
	// This attribute is not available on Oracle Cloud Infrastructure.
	// This attribute is applicable only when provisioning an Oracle Java Cloud Service instance that uses
	// Oracle Identity Cloud Service to configure user authentication and administer users, groups, and roles.
	// Groups properties for the Oracle managed load balancer.
	// Optional
	LoadBalancer *LoadBalancer `json:"loadbalancer,omitempty"`
	// Metering frequency. Valid values include:
	// HOURLY - Pay only for the number of hours used during your billing period. This is the default.
	// MONTHLY - Pay one price for the full month irrespective of the number of hours used.
	// Optional
	MeteringFrequency ServiceInstanceSubscriptionType `json:"meteringFrequency,omitempty"`
	// The email that will be used to send notifications to.
	// To receive notifications, enableNotification must be set to true.
	// Optional
	NotificationEmail string `json:"notificationEmail,omitempty"`
	// This attribute is not applicable to creating service instances on Oracle Cloud at Customer.
	// This attribute is applicable only when provisioning an Oracle Java Cloud Service instance that uses Oracle Identity Cloud Service to configure user authentication and administer users, groups, and roles. Use it when you want to include additional URL patterns to use to protect JavaEE applications.
	// A comma separated list of context roots that you want protected by Oracle Identity Cloud Service.
	// Each context root must begin with the / character. For example:
	// /store/departments/.*,/store/cart/.*,/marketplace/.*,/application1/.*
	ProtectedRootContext string `json:"protectedRootContext,omitempty"`
	// This attribute is applicable only to accounts where regions are supported, including accounts on Oracle Cloud Infrastructure.
	// Name of the region where the Oracle Java Cloud Service instance is to be provisioned.
	// (Not applicable in Oracle Cloud Infrastructure) A region name must be specified if you intend to use
	// ipReservations or ipNetwork.
	// If you do not specify a region, the service instance is created in the site that has the Database Cloud Service
	// database deployment you specify in dbServiceName.
	// If a region name is specified, note that the dbServiceName for this service instance must be one that is provisioned
	// in the same region.
	// An Oracle Cloud Infrastructure region such as us-phoenix-1 must be specified to provision your service instance on
	// Oracle Cloud Infrastructure host resources.
	// Note the following when provisioning in Oracle Cloud Infrastructure:
	// An availability domain must also be specified using availabilityDomain. See Regions and Availability Domains in Oracle
	// Cloud Infrastructure Services.
	// A subnet must also be specified using subnet. See VCNs and Subnets in Oracle Cloud Infrastructure Services.
	// The Oracle Database Cloud Service database deployment on Oracle Cloud Infrastructure must be in the same region and
	// virtual cloud network as the Oracle Java Cloud Service instance you are creating on Oracle Cloud Infrastructure.
	// The service instances do not need to be on the same subnet or availability domain.
	// Cannot use the WLS component property upperStackProductName.
	// Access rules and IP reservations REST endpoints are not supported.
	// Optional.
	Region string `json:"region,omitempty"`
	// Free-form text that provides additional information about the service instance.
	// Optional
	ServiceDescription string `json:"serviceDescription,omitempty"`
	// Service level.
	// Optional
	ServiceLevel ServiceInstanceLevel `json:"serviceLevel,omitempty"`
	// Name of Oracle Java Cloud Service instance. The service name:
	// Must not exceed 30 characters.
	// Must start with a letter.
	// Must contain only letters, numbers, or hyphens.
	// Must not contain any other special characters.
	// Must be unique within the identity domain.
	// By default, the names of the domain and cluster in the service instance will be
	// generated from the first eight characters of the service instance name (serviceName),
	// using the following formats, respectively:
	// first8charsOfServiceInstanceName_domain
	// first8charsOfServiceInstanceName_cluster
	// Required.
	ServiceName string `json:"serviceName"`
	// Oracle WebLogic Server software version. Valid values are: 12cRelease212 (default), 12cR3 and 11gR1.
	// Do not use 11gR1 if you want to create an instance and configure the caching (data grid) cluster at the same time.
	// Only 12cRelease212 is valid if you are using upperStackProductName to provision a service instance for an Oracle
	// Fusion Middleware product.
	// Optional
	ServiceVersion string `json:"serviceVersion,omitempty"`
	// This attribute is applicable only to provisioning on Oracle Cloud Infrastructure Classic.
	// This attribute and sourceServiceName are required when you provision a clone of an existing Oracle
	// Java Cloud Service instance (the source service instance).
	// Name of the snapshot to clone from.
	// Optional
	SnapshotName string `json:"snapshotName,omitempty"`
	// This attribute is applicable only to provisioning on Oracle Cloud Infrastructure Classic.
	// This attribute and snapshotName are required when you provision a clone of an existing Oracle Java Cloud Service instance.
	// Name of the existing Oracle Java Cloud Service instance that has the snapshot from which you are creating a clone.
	// Note the following when creating a service instance clone:
	// Use snapshotName to specify the name of the snapshot to clone from. See Snapshots REST Endpoints.
	// The following service level attributes of the source service instance cannot be changed in the clone: edition,
	// provisionOTD, region, serviceLevel, serviceVersion and useIdentityService. Of those attributes, the clone operation
	// will always use the same values as found in the snapshot of the source service instance. For example, if provisionOTD
	// or useIdentityService is true in the source service instance, then the clone will be provisioned with OTD or with
	// Oracle Identity Cloud Service enabled.
	// In dbServiceName, specify an Oracle Database Cloud Service database deployment clone; do not use the name of the original
	// Database Cloud Service database deployment as the associated database deployment to host the required Oracle schemas for
	// your Oracle Java Cloud Service instance clone. Similarly, for database deployments that host application schemas (if any),
	// use one or more database deployment clones in the array appDBs.
	// If a caching (data grid) cluster is provisioned in the source service instance, the clone will include a caching cluster
	// that has the exact same configuration for clusterName, shape, serverCount, and serversPerNode. Those attributes cannot
	// be changed in the clone.
	// The following WLS component level attributes of the source service instance cannot be changed in the clone: adminPort,
	// clusterName of application cluster, contentPort, deploymentChannelPort, domainName, domainPartitionCount,
	// managedServerCount, nodeManagerPort, sampleAppDeploymentRequested, securedAdminPort, securedContentPort
	// Optional
	SourceServiceName string `json:"sourceServiceName,omitempty"`
	// This attribute is available only on Oracle Cloud Infrastructure. It is required along with region and availabilityDomain.
	// A subdivision of a cloud network that is set up in the data center as specified in availabilityDomain.
	// A subnet exists in a single availability domain and consists of a contiguous range of IP addresses that do not
	// overlap with other subnets in the cloud network. See VCNs and Subnets in Oracle Cloud Infrastructure Services.
	// The subnet must already be created in the specified availability domain. Certain subnet and policy prerequisites
	// must be satisfied when you create your subnet. See Prerequisites for Oracle Platform Services on Oracle Cloud
	// Infrastructure in Oracle Cloud Infrastructure Services.
	// The Oracle Database Cloud Service database deployment on Oracle Cloud Infrastructure must be in the same region
	// and virtual cloud network as the Oracle Java Cloud Service instance you are creating on Oracle Cloud Infrastructure.
	// The service instances do not need to be on the same subnet or availability domain. In this release, however, if the
	// service instances are on different subnets, your Oracle Java Cloud Service instance backup will fail if the backup is
	// configured to include a backup of the associated Oracle Database Cloud Service database deployment. See Backup fails
	// for an Oracle Java Cloud Service instance on Oracle Cloud Infrastructure in Known Issues.
	// Optional
	Subnet string `json:"subnet,omitempty"`
	// The public key for the secure shell (SSH). This key will be used for authentication when connecting to
	// the Oracle Java Cloud Service instance using an SSH client. You generate an SSH public-private key pair using
	// a standard SSH key generation tool. See Generating a Secure Shell (SSH) Public-Private Key Pair in Administering
	// Oracle Java Cloud Service.
	// Required
	VMPublicKeyText string `json:"vmPublicKeyText"`
	// This attribute is not applicable when provisioning an Oracle Java Cloud Service instance in Oracle Cloud Infrastructure.
	// Flag that specifies whether to create (true) or not create (false) the object storage container if the name specified in
	// cloudStorageContainer does not exist. The default is false.
	CloudStorageContainerAutoGenerate bool `json:"cloudStorageContainerAutoGenerate,omitempty"`
	// This attribute is not relevant when provisioning an Oracle Java Cloud Service instance in Oracle Cloud Infrastructure.
	// Flag that specifies whether to enable (true) or disable (false) the access rules that control external communication to
	// the WebLogic Server Administration Console, Fusion Middleware Control, and Load Balancer Console. The default value is false.
	// If you do not set it to true, after the service instance is created, you have to explicitly enable the rules for the
	// administration consoles before you can gain access to the consoles. See Update an Access Rule.
	// Note: On Oracle Cloud Infrastructure, the security rule that controls access to the WebLogic Server Administration
	// Console and other consoles is enabled by default. You cannot disable it during provisioning.
	EnableAdminConsole bool `json:"enableAdminConsole,omitempty"`
	// Flag that specifies whether to enable (true) or disable (false) notifications by email. If this property
	// is set to true, you must specify a value in notificationEmail.
	// Currently, notifications are sent only when service instance provisioning is successful or not successful.
	// Optional
	EnableNotification bool `json:"enableNotification,omitempty"`
	// This attribute is not available on Oracle Cloud at Customer.
	// Flag that specifies whether to apply an existing on-premises license for Oracle WebLogic Server (true) to the new
	// Oracle Java Cloud Service instance you are provisioning. The default value is false.
	// If this property is set to true, you must have a Universal Credits subscription in order to use your existing license.
	// You are responsible for ensuring that you have the required licenses for BYOL instances in Oracle Java Cloud Service.
	// Optional
	IsBYOL bool `json:"isBYOL,omitempty"`
	// Flag that specifies whether to enable the load balancer.
	// The default value is true when you configure more than one Managed Server for the Oracle
	// Java Cloud Service instance. Otherwise, the default value is false
	// Optional.
	ProvisionOTD bool `json:"provisionOTD,omitempty"`
	// This attribute is not available in Oracle Cloud Infrastructure.
	// This attribute is applicable only to accounts that include Oracle Identity Cloud Service.
	// Flag that specifies whether to use Oracle Identity Cloud Service (true) or the local WebLogic identity store
	// (false) for user authentication and to maintain administrators, application users, groups and roles. The default
	// value is false.
	// If you set the value to true, you do not have the option to configure and manage Oracle Traffic Director (OTD) as
	// a local load balancer that runs within your service instance; instead a load balancer is automatically configured
	// and managed for you by Oracle Cloud. See the service parameter loadbalancer.
	// Note the following restrictions for using Oracle Identity Cloud Service (true):
	// serviceLevel must be PAAS
	// provisionOTD must be false
	// serviceVersion cannot be 11gR1 (that is, you cannot run WebLogic Server 11g on the service instance; instead
	// you must use one of the 12c versions)
	// ipNetwork cannot be used
	// See Using Oracle Identity Cloud Service with Oracle Java Cloud Service in Administering Oracle Java Cloud Service.
	// Optional
	UseIdentityService bool `json:"useIdentityService,omitempty"`
}

// LoadBalancer specifies the details of the loadbalancer to create
type LoadBalancer struct {
	// Policy to use for routing requests to the origin servers of the Oracle managed load balancer
	// (that is, when useIdentityService is set to true.
	LoadBalancingPolicy ServiceInstanceLoadBalancingPolicy `json:"loadBalancingPolicy,omitempty"`
}

// CreateComponents specifies the details of the components to create
type CreateComponents struct {
	// Properties for the Oracle Traffic Director (OTD) component.
	// Optional
	OTD *CreateOTD `json:"OTD,omitempty"`
	// Properties for the Oracle WebLogic Server (WLS) component.
	// Required.
	WLS *CreateWLS `json:"WLS"`
}

// CreateOTD specifies the atrributes of the oracle traffic director to create
type CreateOTD struct {
	// Password for the Oracle Traffic Director administrator. The password must meet the following requirements:
	// Starts with a letter
	// Is between 8 and 30 characters long
	// Has one or more upper case letters
	// Has one or more lower case letters
	// Has one or more numbers
	// Has one or more of the following special characters: hyphen (-), underscore (_), pound sign (#), dollar sign ($).
	// If Exadata is the database for the service instance, the password cannot contain the dollar sign ($).
	// If an administrator password is not explicitly set, the OTD administrator password defaults to the WebLogic Server
	// (WLS) administrator password.
	// Optional
	AdminPassword string `json:"adminPassword,omitempty"`
	// User name for the Oracle Traffic Director administrator. The name must be between 8 and 128 characters
	// long and cannot contain any of the following characters:
	// Tab
	// Brackets
	// Parentheses
	// The following special characters: left angle bracket (<), right angle bracket (>), ampersand (&),
	// pound sign (#), pipe symbol (|), and question mark (?).
	// If a username is not explicitly set, the OTD user name defaults to the WebLogic Server (WLS) administrator
	// user name.
	// Optional
	AdminUsername string `json:"adminUserName,omitempty"`
	// Additional Properties Allowed:
	// This attribute is not applicable to Oracle Java Cloud Service instances in Oracle Cloud Infrastructure.
	// Reserved or pre-allocated IP addresses can be assigned to local load balancer nodes.
	// A single IP reservation name or two names separated by a comma.
	// The number of names in ipReservations must match the number of load balancer nodes you are provisioning.
	// Note the difference between accounts where regions are supported and not supported.
	// Where regions are supported: A region name must be specified in order to use ipReservations.
	// Only those reserved IPs created in the specified region can be used.
	// See IP Reservations REST Endpoints for information about how to find unused IP reservations and,
	// if needed, create new IP reservations.
	// Where regions are not supported: If you are using an Oracle Database Exadata Cloud Service database deployment
	// with your Oracle Java Cloud Service instance in an account where regions are not enabled, a region name
	// is not required in order to use ipReservations. However, you must first submit a request to get the
	// IP reservations. See the My Oracle Support document titled How to Request Authorized IPs for Provisioning
	// a Java Cloud Service with Database Exadata Cloud Service (MOS Note 2163568.1).
	// Optional.
	IPReservations []string `json:"ipReservations,omitempty"`
	// Policy to use for routing requests to the load balancer. Valid policies include:
	// Optional.
	LoadBalancingPolicy ServiceInstanceLoadBalancingPolicy `json:"loadBalancingPolicy,omitempty"`
	// Desired compute shape. A shape defines the number of Oracle Compute Units (OCPUs)
	// and amount of memory (RAM).
	// Required.
	Shape ServiceInstanceShape `json:"shape"`
	// Port for accessing Oracle Traffic Director using HTTP. The default value is 8989.
	AdminPort int `json:"adminPort,omitempty"`
	// Listener port for the local load balancer for accessing deployed applications using HTTP.
	// The default value is 8080.
	// This value is overridden by privilegedListenerPort unless its value is set to 0.
	// This value has no effect if the local load balancer is disabled.
	// Optional.
	ListenerPort int `json:"listenerPort,omitempty"`
	//Privileged listener port for accessing the deployed applications using HTTP. The default value is 80.
	// This value has no effect if the local load balancer is disabled.
	// To disable the privileged listener port, set the value to 0. In this case, if the local
	// load balancer is provisioned, the listener port defaults to listenerPort, if specified, or 8080.
	// Optional
	PrivilegedListenerPort int `json:"privilegedListenerPort,omitempty"`
	// Privileged listener port for accessing the deployed applications using HTTPS. The default value is 443.
	// This value has no effect if the local load balancer is disabled.
	// To disable the privileged listener port, set the value to 0. In this case, if the local
	// load balancer is provisioned, the listener port defaults to securedListenerPort, if specified, or 8081.
	// Optional.
	PrivilegedSecuredListenerPort int `json:"privilegedSecuredListenerPort,omitempty"`
	// Secured listener port for accessing the deployed applications using HTTPS. The default value is 8081.
	// This value is overridden by privilegedSecuredContentPort unless its value is set to 0.
	// This value has no effect if the local load balancer is disabled.
	SecuredListenerPort int `json:"securedListenerPort,omitempty"`
	// Flag that specifies whether the local load balancer HA is enabled.
	// This value defaults to false (that is, HA is not enabled).
	// Optional
	HAEnabled bool `json:"haEnabled,omitempty"`
	// Flag that specifies whether the non-secure listener port is enabled on the local load balancer.
	// The default value is true.
	// Optional
	ListenerPortEnabled bool `json:"listenerPortEnabled,omitempty"`
}

// CreateWLS specifies the attributes of the weblogic server to create
type CreateWLS struct {
	// Password for the WebLogic Server administrator. The password must meet the following requirements:
	// Starts with a letter
	// Is between 8 and 30 characters long
	// Has one or more upper case letters
	// Has one or more lower case letters
	// Has one or more numbers
	// Has one or more of the following special characters: hyphen (-), underscore (_), pound sign (#),
	// dollar sign ($). If you are using Exadata as the database for the service instance, the password
	// cannot contain the dollar sign ($).
	// Required
	AdminPassword string `json:"adminPassword"`
	// User name for the WebLogic Server administrator. The name must be between 8 and 128 characters long and cannot contain any of the following characters:
	// Tab
	// Brackets
	// Parentheses
	// The following special characters: left angle bracket (<), right angle bracket (>), ampersand (&),
	// pound sign (#), pipe symbol (|), and question mark (?).
	// Required
	AdminUsername string `json:"adminUserName"`
	// Groups details of Database Cloud Service database deployments that host application schemas, if used.
	// You can specify up to four application schema database deployments.
	// Optional.
	AppDBs []AppDB `json:"appDBs,omitempty"`
	// Size of the backup volume for the service. The value must be a multiple of GBs. You can specify this
	// value in bytes or GBs. If specified in GBs, use the following format: nG, where n specifies the number of GBs.
	// For example, you can express 10 GBs as bytes or GBs. For example: 100000000000 or 10G.
	// This value defaults to the system configured volume size.
	// Optional
	BackupVolumeSize string `json:"backupVolumeSize,omitempty"`
	// This attribute is ignored if clusters array is used.
	// Name of the WebLogic Server application cluster that contains the Managed Servers for running the service applications.
	// The cluster name:
	// Must not exceed 50 characters.
	// Must start with a letter.
	// Must contain only alphabetical characters, underscores (_), or dashes (-).
	// Must not contain any other special characters.
	// Must be unique within the identity domain.
	// If no value is specified, the name of the cluster will be generated from the first eight characters of the Oracle Java Cloud Service instance name (specified in serviceName), using the following format: first8charsOfServiceInstanceName_cluster
	// Optional.
	ClusterName string `json:"clusterName,omitempty"`
	// Groups properties for one or more clusters.
	// This attribute is optional for the WebLogic Server application cluster.
	// You must, however, use the clusters array if you want to define a caching (data grid) cluster for
	// the service instance.
	// Optional.
	Clusters []CreateCluster `json:"clusters,omitempty"`
	// Connection string for the database. The connection string must be entered using one of the following formats:
	// host:port:SID
	// host:port/serviceName
	// For example, foo.bar.com:1521:orcl or foo.bar.com:1521/mydbservice
	// This attribute is required only when you specify a Virtual Image service level of Database
	// Cloud Service in dbServiceName. It is used to connect to the database deployment on Database Cloud Service - Virtual Image.
	// Optional
	ConnectString string `json:"connectString,omitempty"`
	// User name for the database administrator.
	// For service instances based on Oracle WebLogic Server 11g (10.3.6), this value must be set to a database
	// user with DBA role. You can use the default user SYSTEM or a user that has been granted the DBA role.
	// For service instances based on Oracle WebLogic Server 12c (12.2.1 and 12.1.3), this value must be set to
	// a database user with SYSDBA system privileges. You can use the default user SYS or a user that has been
	// granted the SYSDBA privilege.
	// Required.
	DBAName string `json:"dbaName"`
	// Password for the Database administrator that was specified when the Database Cloud Service database deployment was created.
	// Required.
	DBAPassword string `json:"dbaPassword"`
	// Name of the database deployment on Oracle Database Cloud Service to host the Oracle schemas required for this Oracle Java Cloud Service instance.
	// The specified database deployment must be running. Only an Oracle Java Cloud Service instance based on WebLogic Server
	// 12.2.1 can use a required schema database deployment that is created using the Oracle Database 12.2 version.
	// To ensure that you can restore the database for an Oracle Java Cloud Service instance without risking data loss
	// for other service instances, do not use the same Database Cloud Service database deployment with multiple Oracle Java
	// Cloud Service instances.
	// When provisioning a service instance in a specific region, specify a Database Cloud Service database deployment that is
	// in the same region.
	// When provisioning a production-level Oracle Java Cloud Service instance, you must use a production-level Database Cloud
	// Service. The backup option for that database deployment cannot be NONE.
	// (Not applicable to Oracle Cloud Infrastructure) You can specify a Virtual Image service level of Database Cloud Service only
	// if you are provisioning an Oracle Java Cloud Service - Virtual Image instance (BASIC service level). However, you must configure the Oracle Database Cloud Service - Virtual Image environment before you create this service instance. See Using a Database Cloud Service - Virtual Image Database Deployment in Administering Oracle Java Cloud Service.
	// When you specify a Virtual Image service level of Database Cloud Service, you must also specify its connection string using
	// the connectString attribute.
	// When provisioning an Oracle Java Cloud Service instance on Oracle Cloud Infrastructure, note the following:
	// 	The Oracle Database Cloud Service database deployment on Oracle Cloud Infrastructure must be in the same region and virtual
	// 	cloud network as the Oracle Java Cloud Service instance you are creating on Oracle Cloud Infrastructure, but the service
	// 	instances do not have to be on the same subnet. In this release, however, if the service instances are on different subnets,
	// 	your Oracle Java Cloud Service instance backup will fail if the backup is configured to include a backup of the associated
	// 	Oracle Database Cloud Service database deployment. See Backup fails for an Oracle Java Cloud Service instance on
	// 	Oracle Cloud Infrastructure in Known Issues.
	//
	// 	An Oracle Database Exadata Cloud Service database deployment on Oracle Cloud Infrastructure is not supported as a database
	// 	deployment for your Oracle Java Cloud Service instance on Oracle Cloud Infrastructure.
	//
	// 	An Oracle Database Cloud Service database deployment based on a RAC database is also not supported.
	// Required
	DBServiceName string `json:"dbServiceName"`
	// Mode of the domain. Valid values include: DEVELOPMENT and PRODUCTION. The default value is PRODUCTION.
	// Optional
	DomainMode ServiceInstanceDomainMode `json:"domainMode,omitempty"`
	// Name of the WebLogic domain. By default, the domain name will be generated from the first eight characters
	// of the Oracle Java Cloud Service instance name (serviceName), using the following format:
	// first8charsOfServiceInstanceName_domain
	// By default, the Managed Server names will be generated from the first eight characters of the domain name name
	// (domainName), using the following format: first8charsOfDomainName_server_n, where n starts with 1 and is incremented by
	// 1 for each additional Managed Server to ensure each name is unique.
	// Optional
	DomainName string `json:"domainName,omitempty"`
	// Number of partitions to enable in the domain for WebLogic Server 12.2.1.
	// Valid values include: 0 (no partitions), 1, 2, and 4.
	// Optional
	DomainPartitionCount int `json:"domainPartitionCount,omitempty"`
	// Size of the domain volume for the service. The value must be a multiple of GBs.
	// You can specify this value in bytes or GBs. If specified in GBs, use the following format:
	// nG, where n specifies the number of GBs. For example, you can express 10 GBs as bytes or GBs.
	// For example: 100000000000 or 10G.
	// This value defaults to the system configured volume size.
	// Optional
	DomainVolumeSize string `json:"domainVolumeSize,omitempty"`
	// This attribute is not applicable to Oracle Java Cloud Service instances in Oracle Cloud Infrastructure.
	// Reserved or pre-allocated IP addresses can be assigned to Managed Server nodes in a WebLogic Server application cluster.
	// A single IP reservation name or a list of multiple IP reservation names separated by commas.
	// If using reserved IPs, all nodes in the cluster must be provisioned with pre-allocated IP addresses.
	//  In other words, the number of names in ipReservations must match the number of servers you are provisioning
	// (using managedServerCount or serverCount in clusters array).
	// Note the difference between accounts where regions are supported and not supported.
	// Where regions are supported: A region name must be specified in order to use ipReservations.
	// Only those reserved IPs created in the specified region can be used.
	// See IP Reservations REST Endpoints for information about how to find unused IP reservations and, if needed,
	// create new IP reservations.
	// Where regions are not supported: When using an Oracle Database Exadata Cloud Service database deployment with your
	// Oracle Java Cloud Service instance in an account where regions are not enabled, a region name is not required in
	// order to use ipReservations. However, you must first submit a request to get the IP reservations. See the My Oracle
	// Support document titled How to Request Authorized IPs for Provisioning a Java Cloud Service with Database Exadata
	// Cloud Service (MOS Note 2163568.1).
	// Optional
	IPReservations []string `json:"ipReservations,omitempty"`
	// One or more Managed Server JVM arguments separated by a space.
	// You cannot specify any arguments that are related to JVM heap sizes and PermGen spaces (for example, -Xms, -Xmx,
	// -XX:PermSize, and -XX:MaxPermSize).
	// A typical use case would be to set Java system properties using -Dname=value (for example,
	// -Dmyproject.debugDir=/var/myproject/log).
	// You can overwrite or append the default JVM arguments, which are used to start Managed Server processes.
	// See overwriteMsJvmArgs for information on how to overwrite or append the server start arguments.
	// Optional
	MSJvmArgs string `json:"msJvmArgs,omitempty"`
	// Size of the MW_HOME disk volume for the service (/u01/app/oracle/middleware). The value must be a multiple
	// of GBs. You can specify this value in bytes or GBs. If specified in GBs, use the following format: nG, where n
	//specifies the number of GBs. For example, you can express 10 GBs as bytes or GBs. For example: 100000000000 or 10G.
	// This value defaults to the system configured volume size.
	// Optional
	MWVolumeSize string `json:"mwVolumeSize,omitempty"`
	// Password for Node Manager. This value defaults to the WebLogic administrator password (adminPassword)
	// if no value is supplied.
	// Note that the Node Manager password cannot be changed after the Oracle Java Cloud Service instance is provisioned.
	// Optional
	NodeManagerPassword string `json:"nodeManagerPassword,omitempty"`
	// User name for Node Manager. This value defaults to the WebLogic administrator user name (adminUserName)
	// if no value is supplied.
	// Optional
	NodeManagerUserName string `json:"nodeManagerUserName,omitempty"`
	// Name of the pluggable database for Oracle Database 12c. If not specified, the pluggable database name configured when the database was created will be used.
	// Note: This value does not apply to Oracle Database 11g.
	// Optional
	PDBServiceName string `json:"pdbServiceName,omitempty"`
	//Desired compute shape for the nodes in the cluster. A shape defines the number of Oracle Compute Units
	// (OCPUs) and amount of memory (RAM).
	// Required.
	Shape ServiceInstanceShape `json:"shape"`
	// This attribute is not available on Oracle Cloud Infrastructure.
	// This attribute is required only if you are provisioning an Oracle Java Cloud Service instance for an Oracle Fusion
	// Middleware product.
	// The Oracle Fusion Middleware product installer to add to this Oracle Java Cloud Service instance. Valid values are:
	// ODI - Oracle Data Integrator
	// WCP - Oracle WebCenter Portal
	// To use upperStackProductName, you must specify 12cRelease212 as the WebLogic Server software serviceVersion, EE or
	// SUITE as the edition, and PAAS as the serviceLevel.
	// After the service instance is provisioned, the specified Fusion Middleware product installer is available in
	// /u01/zips/upperstack on the Administration Server virtual machine. To install the product over the provisioned domain,
	// follow the instructions provided by the Oracle product's installation and configuration documentation.
	// Optional
	UpperStackProductName ServiceInstanceUpperStackProductName `json:"upperStackProductName,omitempty"`
	// Port for accessing WebLogic Server using HTTP. The default value is 7001.
	// Note that the adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort
	// values must be unique.
	// Optional
	AdminPort int `json:"adminPort,omitempty"`
	// Port for accessing the deployed applications using HTTP.
	// This value is overridden by privilegedContentPort unless its value is set to 0.
	// If a local load balancer is configured and enabled, this value has no effect.
	// Note that the adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort values must be unique.
	// The default value is 8001.
	// Optional.
	ContentPort int `json:"contentPort,omitempty"`
	// Port for accessing the WebLogic Administration Server using WLST.
	// The default value is 9001.
	// Optional
	DeploymentChannelPort int `json:"deploymentChannelPort,omitempty"`
	// Number of Managed Servers in the WebLogic Server application cluster.
	// This attribute is ignored if clusters array is used.
	// Valid values include: 1, 2, 4, and 8. The default value is 1.
	// Optional
	ManagedServerCount int `json:"managedServerCount,omitempty"`
	// Initial Java heap size (-Xms) for a Managed Server JVM, specified in megabytes. The value must be greater than -1.
	// If you specify this initial value, a value greater than 0 (zero) must also be specified for msMaxHeapMB, msMaxPermMB,
	// and msPermMB. In addition, msInitialHeapMB must be less than msMaxHeapMB, and msPermMB must be less than msMaxPermMB.
	// Optional
	MSInitialHeapMB int `json:"msInitialHeapMB,omitempty"`
	// Maximum Java heap size (-Xmx) for a Managed Server JVM, specified in megabytes. The value must be greater than -1.
	// If you specify this maximum value, a value greater than 0 (zero) must also be specified for msInitialHeapMB,
	// msMaxPermMB, and msPermMB. In addition, msInitialHeapMB must be less than msMaxHeapMB, and msPermMB must be less
	// than msMaxPermMB.
	// Optional
	MSMaxHeapMB int `json:"msMaxHeapMB,omitempty"`
	// Maximum Permanent Generation (PermGen) space in Java heap memory (-XX:MaxPermSize) for a Managed Server JVM,
	// specified in megabytes. The value must be greater than -1.
	// Not applicable for a WebLogic Server 12.2.1 instance, which uses JDK 8.
	// If you specify this maximum value, a value greater than 0 (zero) must also be specified for msInitialHeapMB,
	// msMaxHeapMB, and msPermMB. In addition, msInitialHeapMB must be less than msMaxHeapMB, and msPermMB must be
	// less than msMaxPermMB.
	// Optional
	MSMaxPermMB int `json:"msMaxPermMB,omitempty"`
	// Initial Permanent Generation (PermGen) space in Java heap memory (-XX:PermSize) for a Managed Server JVM,
	// specified in megabytes. The value must be greater than -1.
	// Not applicable for a WebLogic Server 12.2.1 instance which uses JDK 8.
	// If you specify this initial value, a value greater than 0 (zero) must also be specified for msInitialHeapMB,
	// msMaxHeapMB, and msMaxPermMB. In addition, msInitialHeapMB must be less than msMaxHeapMB, and msPermMB must be less
	// than msMaxPermMB.
	// Optional
	MSPermMB int `json:"msPermMB,omitempty"`
	// Port for the Node Manager.
	// Node Manager is a WebLogic Server utility that enables you to start, shut down, and restart Administration Server
	// and Managed Server instances from a remote location.
	// The adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort values must be unique.
	// The default value is 5556.
	// Optional
	NodeManagerPort int `json:"nodeManagerPort,omitempty"`
	// Privileged content port for accessing the deployed applications using HTTP.
	// If a local load balancer is configured and enabled, this value has no effect.
	// To disable the privileged content port, set the value to 0. In this case, if a local load balancer is not
	// provisioned, the content port defaults to contentPort, if specified, or 8001.
	// The default value is 80.
	// Optional
	PrivilegedContentPort int `json:"privilegedContentPort,omitempty"`
	// Privileged content port for accessing the deployed applications using HTTPS.
	// If a local load balancer is configured and enabled, this value has no effect.
	// To disable the privileged listener port, set the value to 0. In this case, if a local load balancer is not
	// provisioned, this value defaults to securedContentPort, if specified, or 8002.
	// The default value is 443.
	// Optional
	PrivilegedSecuredContentPort int `json:"privilegedSecuredContentPort,omitempty"`
	// Port for accessing the WebLogic Administration Server using HTTPS.
	// The adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort values must be unique.
	// The default value is 7002.
	// Optional
	SecuredAdminPort int `json:"securedAdminPort,omitempty"`
	// Port for accessing the WebLogic Administration Server using HTTPS. The default value is 8002.
	// This value is overridden by privilegedSecuredContentPort unless its value is set to 0.
	// If a local load balancer is configured and enabled, this value has no effect.
	// The adminPort, contentPort, securedAdminPort, securedContentPort, and nodeManagerPort values must be unique.
	// Optional
	SecuredContentPort int `json:"securedContentPort,omitempty"`
	// Flag that determines whether the user defined Managed Server JVM arguments specified in msJvmArgs should replace the
	// server start arguments (true), or append the server start arguments (false).
	// The server start arguments are calculated automatically by Oracle Java Cloud Service from site default values.
	// If you append (that is, overwriteMsJvmArgs is false or is not set), the user defined arguments specified in msJvmArgs
	// are added to the end of the server start arguments. If you overwrite (that is, set overwriteMsJvmArgs to true), the
	// calculated server start arguments are replaced.
	// Default is false.
	// Optional
	OverwriteMsJvmArgs bool `json:"overwriteMsJvmArgs,omitempty"`
	// Flag that specifies whether to automatically deploy and start the sample application, sample-app.war,
	// to the default Managed Server in your service instance.
	// The default value is false
	// Optional
	SampleAppDeploymentRequested bool `json:"sampleAppDeploymentRequested,omitempty"`
}

// CreateCluster specifies the attributes of the cluster to create
type CreateCluster struct {
	// Name of the cluster to create.
	// The cluster name:
	// Must not exceed 50 characters.
	// Must start with a letter.
	// Must contain only alphabetical characters, underscores (_), or dashes (-).
	// Must not contain any other special characters.
	// Must be unique within the identity domain.
	// Optional.
	ClusterName string `json:"clusterName,omitempty"`
	// A single path prefix or multiple path prefixes separated by commas. A path prefix must be unique across clusters in the domain.
	// This attribute is applicable only to service instances where Oracle Identity Cloud Service is enabled and a managed load balancer is configured. It is also applicable only to a cluster of type APPLICATION_CLUSTER.
	// When path prefixes are specified, this means the load balancer can route to only those applications that have the context root matching one of the path prefixes.
	// For example, if you specified the following path prefixes:
	// ["/myapp1", "/myapp2"]
	// ...then the load balancer can route to only those applications that have the context root matching these:
	// /myapp1
	// /myapp1/*
	// /myapp2
	// /myapp2/*
	// Optional
	PathPrefixes []string `json:"pathPrefixes,omitempty"`
	// Number of servers to create in this cluster.
	// For APPLICATION_CLUSTER - Valid values include: 1, 2, 4, and 8. The default value is 1.
	// For CACHING_CLUSTER - Use a number from 1 to 32 only. The default value is 1.
	// The serverCount limit is based on the VM (cluster size) limit of four and the serversPerNode limit of eight.
	// Note: The actual server number is rounded up to fill the number of nodes required to create
	// the caching cluster. For example, if serversPerNode is four and serverCount is three, the actual
	// number of servers that will be created is four.
	// Optional
	ServerCount int `json:"serverCount,omitempty"`
	// Number of JVMs to start on each VM (node). This attribute is applicable only to cluster type CACHING_CLUSTER.
	// Use a number from 1 to 8 only. The default value is 1.
	// Optional
	ServerPerNode int `json:"serversPerNode,omitempty"`
	// Desired compute shape for the nodes in this cluster. A shape defines the number of Oracle Compute Units
	// (OCPUs) and amount of memory (RAM). Valid shapes on Oracle Cloud Infrastructure Classic include:
	// On Oracle Cloud Infrastructure, only VM.Standard and BM.Standard shapes are supported.
	// See the Bare Metal Shapes and VM Shapes tables of the topic Overview of the Compute Service in Oracle
	// Cloud Infrastructure Services.
	// Note: This shape attribute is optional. If no shape value is specified here, the shape is inherited from
	// the WLS component level shape.
	Shape ServiceInstanceShape `json:"shape,omitempty"`
	// Type of cluster to create.
	// Optional
	Type ServiceInstanceClusterType `json:"type,omitempty"`
}

// AppDB specifies the configuration of the application databases
type AppDB struct {
	// User name for the database administrator.
	// For service instances based on Oracle WebLogic Server 11g (10.3.6), this value must
	// be set to a database user with DBA role. You can use the default user SYSTEM or a user
	// that has been granted the DBA role.
	// For service instances based on Oracle WebLogic Server 12c (12.2.1 and 12.1.3), this value
	// must be set to a database user with SYSDBA system privileges. You can use the default user
	// SYS or a user that has been granted the SYSDBA privilege.
	// Required.
	DBAName string `json:"dbaName"`
	// Database administrator password that was specified when the database deployment on
	// Database Cloud Service was created.
	// Required.
	DBAPassword string `json:"dbaPassword"`
	// Name of the database deployment on Database Cloud Service to use for an application
	// schema. The specified database deployment must be running.
	// Required.
	DBServiceName string `json:"dbServiceName"`
	// Name of the pluggable database for Oracle Database 12c. If not specified,
	// the pluggable database name configured when the database was created will be used.
	// Note: This value does not apply to Oracle Database 11g.
	// Optional.
	PDBServiceName string `json:"pdbServiceName,omitempty"`
}

// CreateServiceInstance creates a new ServiceInstace.
func (c *ServiceInstanceClient) CreateServiceInstance(input *CreateServiceInstanceInput) (*ServiceInstance, error) {
	if c.PollInterval == 0 {
		c.PollInterval = waitForServiceInstanceReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForServiceInstanceReadyTimeout
	}

	// Since these CloudStorageUsername and CloudStoragePassword are sensitive we'll read them
	// from the environment if they aren't passed in.
	if input.CloudStorageContainer != "" && input.CloudStorageUsername == "" && input.CloudStoragePassword == "" {
		input.CloudStorageUsername = *c.ResourceClient.Client.client.UserName
		input.CloudStoragePassword = *c.ResourceClient.Client.client.Password
	}

	// The JCS API errors if an ssh key has trailing content; we'll trim that here.
	parts := strings.Split(input.VMPublicKeyText, " ")
	if len(parts) > 2 {
		input.VMPublicKeyText = strings.Join(parts[0:2], " ")
	}

	serviceInstance, err := c.startServiceInstance(input.ServiceName, input)
	if err != nil {
		return serviceInstance, fmt.Errorf("unable to create Java Service Instance %q: %+v", input.ServiceName, err)
	}
	return serviceInstance, nil
}

func (c *ServiceInstanceClient) startServiceInstance(name string, input *CreateServiceInstanceInput) (*ServiceInstance, error) {
	if err := c.createResource(*input, nil); err != nil {
		return nil, err
	}

	// Call wait for instance ready now, as creating the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: name,
	}

	// Wait for the service instance to be running and return the result
	// Don't have to unqualify any objects, as the GetServiceInstance method will handle that
	serviceInstance, err := c.WaitForServiceInstanceState(getInput, ServiceInstanceLifecycleStateStart, c.PollInterval, c.Timeout)
	// If the service instance is returned as nil if it enters a terminating state.
	if err != nil || serviceInstance == nil {
		return nil, fmt.Errorf("error creating service instance %q: %+v", name, err)
	}
	return serviceInstance, nil
}

// WaitForServiceInstanceState waits for a service instance to be in the desired state
func (c *ServiceInstanceClient) WaitForServiceInstanceState(input *GetServiceInstanceInput, desiredState ServiceInstanceLifecycleState, pollInterval, timeoutSeconds time.Duration) (*ServiceInstance, error) {
	var info *ServiceInstance
	var getErr error
	err := c.client.WaitFor("service instance to be ready", pollInterval, timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetServiceInstance(input)
		if getErr != nil {
			return false, getErr
		}
		c.client.DebugLogString(fmt.Sprintf("Service instance name is %v, Service instance info is %+v", info.ServiceName, info))
		switch s := info.State; s {
		case ServiceInstanceStatusReady: // Target State
			c.client.DebugLogString("Service Instance Ready")
			if desiredState == ServiceInstanceLifecycleStateStart || desiredState == ServiceInstanceLifecycleStateRestart {
				return true, nil
			}
			return false, nil
		case ServiceInstanceStatusConfiguring:
			c.client.DebugLogString("Service Instance is being created")
			return false, nil
		case ServiceInstanceStatusInitializing:
			c.client.DebugLogString("Service Instance is being initialized")
			return false, nil
		case ServiceInstanceStatusStopping:
			c.client.DebugLogString("ServiceInstance is stopping")
			return false, nil
		case ServiceInstanceStatusStopped:
			c.client.DebugLogString("ServiceInstance is stopped")
			if desiredState == ServiceInstanceLifecycleStateStop {
				return true, nil
			}
			return false, nil
		case ServiceInstanceStatusTerminating:
			c.client.DebugLogString("Service Instance creation failed, terminating")
			// The Service Instance creation failed. Wait for the instance to be deleted.
			return false, c.waitForServiceInstanceDeleted(input, pollInterval, timeoutSeconds)
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown instance state: %s, waiting", s))
			return false, nil
		}
	})
	return info, err
}

// GetServiceInstanceInput specifies which service instance to retrieve
type GetServiceInstanceInput struct {
	// Name of the Java Cloud Service instance.
	// Required.
	Name string `json:"serviceId"`
}

// GetServiceInstance retrieves the SeriveInstance with the given name.
func (c *ServiceInstanceClient) GetServiceInstance(getInput *GetServiceInstanceInput) (*ServiceInstance, error) {
	var serviceInstance ServiceInstance
	if err := c.getResource(getInput.Name, &serviceInstance); err != nil {
		return nil, err
	}

	return &serviceInstance, nil
}

// DeleteServiceInstanceInput specifies which service instance to delete
type DeleteServiceInstanceInput struct {
	// Name of the Java Cloud Service instance.
	// Required.
	Name string `json:"-"`
	// User name for the database administrator.
	// Required.
	DBAUsername string `json:"dbaName"`
	// The database administrator password that was specified when the Database Cloud Service database deployment
	// was created or the password for the database administrator.
	// Required.
	DBAPassword string `json:"dbaPassword"`
	// Flag that specifies whether you want to force the removal of the service instance even if the database
	// instance cannot be reached to delete the database schemas. If set to true, you may need to delete the associated
	// database schemas manually on the database instance if they are not deleted as part of the service instance
	// delete operation.
	// The default value is false.
	// Optional.
	ForceDelete bool `json:"force,omitempty"`
	// Flag that specifies whether you want to back up the service instance or skip backing up the instance before deleting it.
	// The default value is true (that is, skip backing up).
	// Optional.
	SkipBackupOnTerminate bool `json:"skipBackupOnTerminate,omitempty"`
}

// DeleteServiceInstance deletes the specified service instance
func (c *ServiceInstanceClient) DeleteServiceInstance(deleteInput *DeleteServiceInstanceInput) error {
	if c.PollInterval == 0 {
		c.PollInterval = waitForServiceInstanceDeletePollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForServiceInstanceDeleteTimeout
	}

	// There are times when the service instance isn't in a state to be deleted even though the api returns a ready
	// instance. We'll wait a set amount of time for it to be ready to delete before erroring out.
	var deleteErr error
	for i := 0; i < deleteMaxRetries; i++ {
		c.client.DebugLogString(fmt.Sprintf("(Iteration: %d of %d) Deleting instance with name %s", i, *c.Client.client.MaxRetries, deleteInput.Name))

		deleteErr = c.deleteInstanceResource(deleteInput.Name, deleteInput)
		if deleteErr == nil {
			c.client.DebugLogString(fmt.Sprintf("(Iteration: %d of %d) Finished deleting instance with name %s", i, *c.Client.client.MaxRetries, deleteInput.Name))
			break
		}
		time.Sleep(1 * time.Minute)
	}
	if deleteErr != nil {
		return fmt.Errorf("error submitting delete request for java service instance %q", deleteInput.Name)
	}

	// Call wait for instance deleted now, as deleting the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: deleteInput.Name,
	}

	// Wait for instance to be deleted
	return c.waitForServiceInstanceDeleted(getInput, c.PollInterval, c.Timeout)
}

// WaitForServiceInstanceDeleted waits for a service instance to be fully deleted.
func (c *ServiceInstanceClient) waitForServiceInstanceDeleted(input *GetServiceInstanceInput, pollInterval, timeoutSeconds time.Duration) error {
	return c.client.WaitFor("service instance to be deleted", pollInterval, timeoutSeconds, func() (bool, error) {
		info, err := c.GetServiceInstance(input)
		if err != nil {
			if client.WasNotFoundError(err) {
				// Service Instance could not be found, thus deleted
				return true, nil
			}
			// Some other error occurred trying to get instance, exit
			return false, err
		}
		switch s := info.State; s {
		case ServiceInstanceStatusTerminating:
			c.client.DebugLogString("Service Instance terminating")
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown instance state: %s, waiting", s))
			return false, nil
		}
	})
}

// ScaleUpDownServiceInstanceInput defines the attributes for how to scale up or down the java service instance.
type ScaleUpDownServiceInstanceInput struct {
	// Groups properties for the Oracle WebLogic Server component (WLS).
	// Required
	Components ScaleUpDownComponent `json:"components"`
	// Name of the Java Cloud Service instance.
	// Required.
	Name string `json:"-"`
}

// ScaleUpDownComponent defines the attributes for the WebLogic Server components when scaling up and down
// the service instance.
type ScaleUpDownComponent struct {
	// Properties for the Oracle WebLogic Server (WLS) component.
	// Required
	WLS ScaleUpDownWLS `json:"WLS"`
}

// ScaleUpDownWLS defines the properties for the Oracle WebLogic Server (WLS) component.
type ScaleUpDownWLS struct {
	// A single host name. Only application cluster hosts can be specified.
	// Required
	Hosts []string `json:"hosts"`
	// Flag that indicates whether to ignore Managed Server heap validation (true) or perform heap
	// validation (false) before a scale down request is accepted. Default is false.
	// When the flag is not set or is false, heap validation is performed before scaling.
	// If a validation error is not generated, the Managed Server JVM is restarted with the new shape
	// after scaling down.
	// When the flag is true, heap validation is not performed. Before you set the flag to true,
	// make sure the -Xms value is low enough for the Managed Server JVM to restart on the new shape
	// after scaling down. The -Xms value should be lower than one-fourth the size of the memory
	// associated with the shape. Use the WebLogic Server Administration Console to edit the value in
	// the server start arguments, if necessary.
	// Optional
	IgnoreManagedServerHeapError bool `json:"ignoreManagedServerHeapError,omitempty"`
	// Desired compute shape for the target host.
	// Required
	Shape ServiceInstanceShape `json:"shape"`
}

// ScaleUpDownServiceInstance scales the service instance up or down depending on the shape passed in.
func (c *ServiceInstanceClient) ScaleUpDownServiceInstance(input *ScaleUpDownServiceInstanceInput) error {
	if c.PollInterval == 0 {
		c.PollInterval = waitForServiceInstanceReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForServiceInstanceReadyTimeout
	}

	if err := c.updateResource(input.Name, serviceInstanceScaleUpDownPath, "POST", input); err != nil {
		return fmt.Errorf("unable to update Java Service Instance %q: %+v", input.Name, err)
	}

	// Call wait for instance ready now, as updating the instance is an eventually consistent operation.
	getInput := &GetServiceInstanceInput{
		Name: input.Name,
	}

	// Wait for the service instance to be running and return the result
	// Don't have to unqualify any objects, as the GetServiceInstance method will handle that.
	serviceInstance, err := c.WaitForServiceInstanceState(getInput, ServiceInstanceLifecycleStateStart, c.PollInterval, c.Timeout)
	// The service instance is returned as nil if it enters a terminating state.
	if err != nil || serviceInstance == nil {
		return fmt.Errorf("error creating service instance %q: %+v", input.Name, err)
	}

	return nil
}

// DesiredStateInput defines the attributes for how to set the desired state of a java service instance.
type DesiredStateInput struct {
	// Flag that specifies whether to control the entire service instance.
	// This attribute is not applicable to the restart command.
	// Optional
	AllServiceHosts bool `json:"allServiceHosts,omitemtpy"`
	// Groups properties for the Oracle WebLogic Server component (WLS).
	// Optional
	Components *DesiredStateComponent `json:"components,omitempty"`
	// Name of the Java Cloud Service instance.
	// Required.
	Name string `json:"-"`
	// Type of the request.
	// Required
	LifecycleState ServiceInstanceLifecycleState `json:"-"`
}

// DesiredStateComponent groups properties for the Oracle WebLogic Server component (WLS) or the Oracle
// Traffice Director (OTD) component.
type DesiredStateComponent struct {
	// Properties for the Oracle Traffic Director (OTD) component.
	// Optional
	OTD *DesiredStateHost `json:"OTD,omitempty"`
	// Properties for the Oracle WebLogic Server (WLS) component.
	// Optional
	WLS *DesiredStateHost `json:"WLS,omitempty"`
}

// DesiredStateHost defines the properties of the hosts
type DesiredStateHost struct {
	// A single host name. Only application cluster hosts can be specified.
	// Required
	Hosts []string `json:"hosts"`
}

// UpdateDesiredState updates the specified desired state of a service instance
func (c *ServiceInstanceClient) UpdateDesiredState(input *DesiredStateInput) error {
	if c.PollInterval == 0 {
		c.PollInterval = waitForServiceInstanceReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForServiceInstanceReadyTimeout
	}

	if err := c.updateResource(input.Name, fmt.Sprintf(serviceInstanceDesiredStatePath, input.LifecycleState), "POST", input); err != nil {
		return err
	}

	// Call wait for instance running now, as updating the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: input.Name,
	}

	// Wait for the service instance to be running and return the result
	// Don't have to unqualify any objects, as the GetServiceInstance method will handle that
	_, err := c.WaitForServiceInstanceState(getInput, input.LifecycleState, c.PollInterval, c.Timeout)
	if err != nil {
		return fmt.Errorf("Error updating Service Instance %q: %+v", input.Name, err)
	}
	return nil
}
