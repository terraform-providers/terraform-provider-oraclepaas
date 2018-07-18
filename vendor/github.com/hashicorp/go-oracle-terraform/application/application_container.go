package application

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-oracle-terraform/client"
)

const waitForApplicationContainerRunningPollInterval = time.Second * 10
const waitForApplicationContainerRunningTimeout = time.Second * 600
const waitForApplicationContainerDeletePollInterval = time.Second * 10
const waitForApplicationContainerDeleteTimeout = time.Second * 600

var (
	// ContainerContainerPath is the uri path for containers
	ContainerContainerPath = "/paas/service/apaas/api/v1.1/apps/%s"
	// ContainerResourcePath is the uri path for a specific container resource
	ContainerResourcePath = "/paas/service/apaas/api/v1.1/apps/%s/%s"
)

// ContainerClient is a client for the Application container functions of the Application Container API
type ContainerClient struct {
	ResourceClient
	PollInterval time.Duration
	Timeout      time.Duration
}

// ContainerClient obtains an ApplicationClient which can be used to access to the
// Application Container functions of the Application Container API
func (c *Client) ContainerClient() *ContainerClient {
	return &ContainerClient{
		ResourceClient: ResourceClient{
			Client:           c,
			ContainerPath:    ContainerContainerPath,
			ResourceRootPath: ContainerResourcePath,
		}}
}

/* These values would normally be used as types for their respective attributes but we are unable to because we need to convert a struct to a map[string[string]

type ApplicationRepository string

const (
	ApplicationRepositoryDockerHub ApplicationRepository = "dockerhub"
)

type ApplicationRuntime string

const (
	ApplicationRuntimeJava   ApplicationRuntime = "java"
	ApplicationRuntimeJavaEE ApplicationRuntime = "javaee"
	ApplicationRuntimeNode   ApplicationRuntime = "node"
	ApplicationRuntimePHP    ApplicationRuntime = "php"
	ApplicationRuntimePython ApplicationRuntime = "python"
	ApplicationRuntimeRuby   ApplicationRuntime = "ruby"
	ApplicationRuntimeGo     ApplicationRuntime = "golang"
	ApplicationRuntimeDotnet ApplicationRuntime = "dotnet"
)

type ApplicationContainerAuthType string

const (
	ApplicationContainerAuthTypeBasic ApplicationContainerAuthType = "basic"
	ApplicationContainerAuthTypeOauth ApplicationContainerAuthType = "oauth"
) */

// SubscriptionType are the values available for subscription
type SubscriptionType string

const (
	// SubscriptionTypeHourly specifies that the subscription is metered hourly
	SubscriptionTypeHourly SubscriptionType = "HOURLY"
	// SubscriptionTypeMonthly specifies that the subscription is metered monthly
	SubscriptionTypeMonthly SubscriptionType = "MONTHLY"
)

type containerStatus string

const (
	applicationContainerStatusRunning        containerStatus = "RUNNING"
	applicationContainerStatusNew            containerStatus = "NEW"
	applicationContainerStatusDestroyPending containerStatus = "DESTROY_PENDING"
)

// ManifestType determines whether an application is public or private:
type ManifestType string

const (
	// ManifestTypeWeb specifies a public application, which you can access using a public URL, the public REST API, or the command-line interface.
	ManifestTypeWeb ManifestType = "web"
	// ManifestTypeWorker specifies a worker application, which is private and runs in the background. The isClustered parameter should be set to true in some cases.
	ManifestTypeWorker ManifestType = "worker"
)

// ManifestMode details the optional modes for restarting the application container
type ManifestMode string

const (
	// ManifestModeRolling performs a rolling restart
	ManifestModeRolling ManifestMode = "rolling"
)

// ServiceType details the constants related to what types a service is
type ServiceType string

const (
	// ServiceTypeJAAS for Oracle Java Cloud Service
	ServiceTypeJAAS ServiceType = "JAAS"
	// ServiceTypeDBAAS for Oracle Database Cloud Service
	ServiceTypeDBAAS ServiceType = "DBAAS"
	// ServiceTypeMYSQLCS for MySQL Cloud Service
	ServiceTypeMYSQLCS ServiceType = "MYSQLCS"
	// ServiceTypeOEHCS for an Oracle Event Hub Cloud Service topic
	ServiceTypeOEHCS ServiceType = "OEHCS"
	// ServiceTypeOEHPCS for an Oracle Event Hub Cloud Service cluster
	ServiceTypeOEHPCS ServiceType = "OEHPCS"
	// ServiceTypeDHCS for Oracle Data Hub Cloud Service
	ServiceTypeDHCS ServiceType = "DHCS"
	// ServiceTypeCaching for a cache
	ServiceTypeCaching ServiceType = "caching"
)

// Container container information about the application container
type Container struct {
	// ID of the application
	AppID string `json:"appId"`
	// URL of the created application
	AppURL string `json:"appURL"`
	// Creation time of the application
	CreatedTime string `json:"createdTime"`
	// The activity acting on the application container
	CurrentOnGoingActitvity string `json:"currentOngoingActivity"`
	// Identity Domain of the application
	IdentityDomain string `json:"identityDomain"`
	// Shows details of all instances currently running.
	Instances []Instance `json:"instances"`
	// Modification time of the application
	LastModifiedTime string `json:"lastModifiedTime"`
	// Shows all deployments currently in progress.
	LatestDeployment Deployment `json:"latestDeployment"`
	// Name of the application
	Name string `json:"name"`
	// Shows all deployments currently running.
	RunningDeployment Deployment `json:"runningDeployment"`
	// Status of the application
	Status string `json:"status"`
	// Type of subscription, Hourly or Monthly
	SubscriptionType SubscriptionType `json:"subscriptionType"`
	// Web URL of the application
	WebURL string `json:"webURL"`
}

// Instance specifies individual instance information on the application container
type Instance struct {
	// Instance URL. Use this url to get a description of the application instance.
	InstanceURL string `json:"instanceURL"`
	// Memory of the instance
	Memory string `json:"memory"`
	// Instance Name. Use this name to manage a specific instance.
	Name string `json:"name"`
	// Status of the instance
	Status string `json:"status"`
}

// Deployment specifies individual deployment information on the application container
type Deployment struct {
	// Deployment ID. Use this ID to manage a specific deployment.
	DeploymentID string `json:"deploymentId"`
	// Status of the deployment
	DeploymentStatus string `json:"deploymentStatus"`
	// Deployment URL. Use this URL to get a description of the application deployment.
	DeploymentURL string `json:"deploymentURL"`
}

// CreateApplicationContainerInput specifies the information needed to create an application container
type CreateApplicationContainerInput struct {
	// The additional fields needed for ApplicationContainer
	AdditionalFields CreateApplicationContainerAdditionalFields
	// Name of the optional deployment file, which specifies memory, number of instances, and service bindings
	// Optional
	Deployment string
	// Deployment Attributes
	// Optional
	DeploymentAttributes *DeploymentAttributes
	// Name of the manifest file, required if this file is not packaged with the application
	// Optional
	Manifest string
	// Manifest Attributes
	// Optional
	ManifestAttributes *ManifestAttributes
	// Time to wait between checks on application container status
	PollInterval time.Duration
	// Timeout for creating an application container
	Timeout time.Duration
}

// CreateApplicationContainerAdditionalFields specifies the additional fields needed to create an application container
type CreateApplicationContainerAdditionalFields struct {
	// Location of the application archive file in Oracle Storage Cloud Service, in the format app-name/file-name
	// Optional
	ArchiveURL string `json:"archiveURL"`
	// Uses Oracle Identity Cloud Service to control who can access your Java SE 7 or 8, Node.js, or PHP application.
	// This should be ApplicationContainerAuthType but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Allowed values are 'basic' and 'oauth'.
	// Optional
	AuthType string `json:"authType"`
	// The password of your GitHub repository, required if your repository is private.
	// Optional
	GitPassword string `json:"gitPassword"`
	// URL of your GitHub repository.
	GitRepoURL string `json:"gitRepoUrl"`
	// The user name of your GitHub repository, required if your repository is private.
	GitUsername string `json:"gitUserName"`
	// Name of the application
	// Required
	Name string `json:"name"`
	// Comments on the application deployment
	Notes string `json:"notes"`
	// Email address to which application deployment status updates are sent.
	NotificationEmail string `json:"notifcationEmail"`
	// Repository of the application. The only allowed value is 'dockerhub'.
	// This should be ApplicationRepository but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Optional
	Repository string `json:"repository"`
	// Runtime environment: java (the default), node, php, python, or ruby
	// This should be ApplicationRuntime but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Required
	Runtime string `json:"runtime"`
	// Subscription, either hourly (the default) or monthly
	// This should be ApplicationSubscriptionType but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Optional
	SubscriptionType string `json:"subscription"`
	// List of tags that can be associated with the application.
	Tags []Tag `json:"tags"`
}

// ManifestAttributes details the available attributes in a manifest file
type ManifestAttributes struct {
	// Optional
	Runtime Runtime `json:"runtime,omitempty"`
	// Determines whether an application is public or private
	// Default is `worker`
	// Optional
	Type ManifestType `json:"type,omitempty"`
	// Launch command to execute after the application has been uploaded.
	// Optional
	Command string `json:"command,omitempty"`
	// Release attributes of a specific build.
	// Optional
	Release Release `json:"release,omitempty"`
	// Maximum time in seconds to wait for the application to start. Allowed values are between 10 and 600. The default is 30.
	// If the application doesn’t start in the time specified, the application is deemed to have failed to start and is terminated.
	// For example, if your application takes two minutes to start, set startupTime to at least 120.
	// Optional
	StartupTime string `json:"startupTime,omitempty"`
	// Maximum time in seconds to wait for the application to stop. Allowed values are between 0 and 600.
	// The default is 0. This allows the application to close connections and free up resources gracefully.
	// For example, if your application takes two minutes to shut down, set shutdownTime to at least 120.
	// Optional
	ShutdownTime string `json:"shutdownTime,omitempty"`
	// Comments
	// Optional
	Notes string `json:"notes,omitempty"`
	// Restart mode for application instances when the application is restarted. The only allowed option is rolling for a rolling restart.
	// Omit this parameter to be prompted for a rolling or concurrent restart.
	// Optional
	Mode ManifestMode `json:"mode,omitempty"`
	// Must be set to true for application instances to act as a cluster, with failover capability.
	// Optional
	IsClustered bool `json:"isClustered,omitempty"`
	// Context root of the application. The value of the home parameter is appended to the application URL.
	// Optional
	Home string `json:"home,omitempty"`
	// Allows you to define a URL for your application that the system uses for health checks. The URL must return an HTTP response of 200 OK to indicate that the application is healthy.
	// Optional
	HealthCheck HealthCheck `json:"healthcheck,omitempty"`
}

// DeploymentAttributes details the attributes available for a deployment.json file
type DeploymentAttributes struct {
	// The amount of memory in gigabytes made available to the application. Values range from 1G to 20G.
	// The default is 2G.
	// Optional
	Memory string `json:"memory,omitempty"`
	// Number of application instances. The maximum is 64.
	// The default is 2.
	// Optional
	Instances int `json:"instances,omitempty"`
	// Free-form comment field. It can be used, for example, to describe the deployment plan.
	// Optional
	Notes string `json:"notes,omitempty"`
	// Environment variables used by the application. This element can contain any number of name-value pairs.
	// Optional
	Envrionment map[string]string `json:"environment,omitempty"`
	// List of environment variables marked as secured on the user interface.
	// The environment variables to be secured must be present in the environment property.
	// Optional
	SecureEnvironment []string `json:"secureEnvironment,omitempty"`
	// Java EE system properties used by the application. This element can contain any number of name-value pairs.
	// Optional
	JavaSystemProperties map[string]string `json:"java_system_properties,omitempty"`
	// Service bindings for connections to other Oracle Cloud services. This element can contain more
	// than one block of sub-elements. Each block specifies one service binding.
	// Optional
	Services []Service `json:"services,omitempty"`
}

// Release details the attributes for a specific release for the application container.
type Release struct {
	// User-specified value of build.
	// Required
	Build string `json:"build"`
	// User-specified value of commit.
	// Required
	Commit string `json:"commit"`
	// User-specified application version.
	// Required
	Version string `json:"version"`
}

// Runtime details the available runtime attributes for a manifest file
type Runtime struct {
	MajorVersion string `json:"majorVersion"`
}

// HealthCheck specifies the available attributes for a health check
type HealthCheck struct {
	// Defines the URI that is appended to the application URL to create the health check URL
	HTTPEndpoint string `json:"http-endpoint"`
}

// Service details the service bindings for connections to other Oracle Cloud services.
type Service struct {
	// User-specified identifier.
	// Required
	Identifier string `json:"identifier"`
	// Type of the service.
	// Required
	Type ServiceType `json:"type"`
	// Name of the service, the name of an existing Oracle Java Cloud Service instance, Oracle Database Cloud
	// Service database, MySQL Cloud Service database, Oracle Event Hub Cloud Service topic or cluster,
	// Oracle Data Hub Cloud Service instance, or cache name.
	// Required
	Name string `json:"name"`
	// Username used to access the service.
	// Required
	Username string `json:"username"`
	// Password for the username.
	// Required
	Password string `json:"password"`
}

// Tag defines the attributes related to a tag
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateApplicationContainer creates a new Application Container from an ApplicationClient and an input struct.
// Returns a populated ApplicationContainer struct for the Application, and any errors
func (c *ContainerClient) CreateApplicationContainer(input *CreateApplicationContainerInput) (*Container, error) {

	var applicationContainer *Container

	if input.Manifest != "" && input.ManifestAttributes != nil {
		return nil, fmt.Errorf("Cannot specify both a manifest file and manifest attributes")
	}
	if input.Deployment != "" && input.DeploymentAttributes != nil {
		return nil, fmt.Errorf("Cannot specify both a deployment file and deployment attributes")
	}

	files := make(map[string][]byte)
	if input.Deployment != "" {
		fileContents, err := readFileContents(input.Deployment)
		if err != nil {
			return nil, fmt.Errorf("Error reading deployment file: %+v", err)
		}
		files["deployment"] = fileContents
	}
	if input.Manifest != "" {
		fileContents, err := readFileContents(input.Manifest)
		if err != nil {
			return nil, fmt.Errorf("Error reading manifest file: %+v", err)
		}
		files["manifest"] = fileContents
	}

	additionalFields := make(map[string]interface{})
	additionalFieldsStruct := structs.New(input.AdditionalFields)
	additionalFieldsValues := additionalFieldsStruct.Values()
	for i, v := range additionalFieldsValues {
		key := additionalFieldsStruct.Fields()[i]
		if key.Name() == "Tags" {
			modifiedTags := make([]string, 0)
			tags := v.([]Tag)
			for _, tag := range tags {
				modifiedTags = append(modifiedTags, fmt.Sprintf("{'key': %q, 'value': %q}", tag.Key, tag.Value))
			}
			additionalFields[key.Tag("json")] = fmt.Sprintf("[%s]", strings.Join(modifiedTags, ","))
		} else {
			additionalFields[key.Tag("json")] = v
		}
	}

	if input.ManifestAttributes != nil {
		manifestBytes, err := c.client.MarshallRequestBody(input.ManifestAttributes)
		if err != nil {
			return nil, err
		}
		files["manifest"] = manifestBytes
	}

	if input.DeploymentAttributes != nil {
		deploymentBytes, err := c.client.MarshallRequestBody(input.DeploymentAttributes)
		if err != nil {
			return nil, err
		}
		files["deployment"] = deploymentBytes
	}

	if err := c.createApplicationContainerResource("POST", files, additionalFields, applicationContainer); err != nil {
		return nil, err
	}

	getInput := &GetApplicationContainerInput{
		Name: input.AdditionalFields.Name,
	}

	if input.PollInterval == 0 {
		input.PollInterval = waitForApplicationContainerRunningPollInterval
	}
	if input.Timeout == 0 {
		input.Timeout = waitForApplicationContainerRunningTimeout
	}

	// Wait for application container to be ready and return the result
	applicationContainerInfo, err := c.WaitForApplicationContainerRunning(getInput, input.PollInterval, input.Timeout)
	if err != nil {
		return nil, err
	}

	return applicationContainerInfo, nil
}

// GetApplicationContainerInput specifies the information needed to get an application container
type GetApplicationContainerInput struct {
	// Name of the application container
	// Required
	Name string `json:"name"`
}

// GetApplicationContainer retrieves the application container with the given name.
func (c *ContainerClient) GetApplicationContainer(getInput *GetApplicationContainerInput) (*Container, error) {
	var applicationContainer Container
	if err := c.getResource(getInput.Name, &applicationContainer); err != nil {
		return nil, err
	}

	return &applicationContainer, nil
}

// DeleteApplicationContainerInput specifies the information needed to delete an application container
type DeleteApplicationContainerInput struct {
	// Name of the application container
	// Required
	Name string `json:"name"`
	// Time to wait between checks on application container status
	PollInterval time.Duration
	// Timeout to delete an application container
	Timeout time.Duration
}

// DeleteApplicationContainer deletes the application container with the given name.
func (c *ContainerClient) DeleteApplicationContainer(input *DeleteApplicationContainerInput) error {
	// Call to delete the application container
	if err := c.deleteResource(input.Name); err != nil {
		return err
	}

	if input.PollInterval == 0 {
		input.PollInterval = waitForApplicationContainerDeletePollInterval
	}
	if input.Timeout == 0 {
		input.Timeout = waitForApplicationContainerDeleteTimeout
	}

	// Wait for application container to be deleted
	return c.WaitForApplicationContainerDeleted(input, input.PollInterval, input.Timeout)
}

// UpdateApplicationContainerInput specifies the fields needed to update an application container
type UpdateApplicationContainerInput struct {
	// Name of the application container
	Name string
	// The additional fields needed for ApplicationContainer
	AdditionalFields UpdateApplicationContainerAdditionalFields
	// Name of the optional deployment file, which specifies memory, number of instances, and service bindings
	// Optional
	Deployment string
	// Deployment Attributes
	// Optional
	DeploymentAttributes *DeploymentAttributes
	// Name of the manifest file, required if this file is not packaged with the application
	// Optional
	Manifest string
	// Manifest Attributes
	ManifestAttributes *ManifestAttributes
	// Time to wait between checks on application container status
	PollInterval time.Duration
	// Timeout for creating an application container
	Timeout time.Duration
}

// UpdateApplicationContainerAdditionalFields specifies the additional fields needed to update an application container
type UpdateApplicationContainerAdditionalFields struct {
	// Location of the application archive file in Oracle Storage Cloud Service, in the format app-name/file-name
	// Optional
	ArchiveURL string
	// Comments on the application deployment
	Notes string
}

// UpdateApplicationContainer updates an application container from an ApplicationClient and an input struct.
// Returns a populated ApplicationContainer struct for the Application, and any errors
func (c *ContainerClient) UpdateApplicationContainer(input *UpdateApplicationContainerInput) (*Container, error) {

	// To prevent potentially overriding attributes when passing a file and attributes, we only allow one or the other.
	if input.Manifest != "" && input.ManifestAttributes != nil {
		return nil, fmt.Errorf("Cannot specify both a manifest file and manifest attributes")
	}
	if input.Deployment != "" && input.DeploymentAttributes != nil {
		return nil, fmt.Errorf("Cannot specify both a deployment file and deployment attributes")
	}

	additionalFields := make(map[string]interface{})
	additionalFieldsStruct := structs.New(input.AdditionalFields)
	additionalFieldsValues := additionalFieldsStruct.Values()
	for i, v := range additionalFieldsValues {
		additionalFields[additionalFieldsStruct.Fields()[i].Tag("json")] = v
	}

	files := make(map[string][]byte)
	if input.Deployment != "" {
		fileContents, err := readFileContents(input.Deployment)
		if err != nil {
			return nil, fmt.Errorf("Error reading deployment file: %+v", err)
		}
		files["deployment"] = fileContents
	}
	if input.Manifest != "" {
		fileContents, err := readFileContents(input.Manifest)
		if err != nil {
			return nil, fmt.Errorf("Error reading manifest file: %+v", err)
		}
		files["manifest"] = fileContents
	}

	if input.ManifestAttributes != nil {
		manifestBytes, err := c.client.MarshallRequestBody(input.ManifestAttributes)
		if err != nil {
			return nil, err
		}
		files["manifest"] = manifestBytes
	}

	if input.DeploymentAttributes != nil {
		deploymentBytes, err := c.client.MarshallRequestBody(input.DeploymentAttributes)
		if err != nil {
			return nil, err
		}
		files["deployment"] = deploymentBytes
	}

	var applicationContainer *Container
	if err := c.updateApplicationContainerResource(input.Name, "PUT", files, additionalFields, applicationContainer); err != nil {
		return nil, err
	}

	getInput := &GetApplicationContainerInput{
		Name: input.Name,
	}

	if input.PollInterval == 0 {
		input.PollInterval = waitForApplicationContainerRunningPollInterval
	}
	if input.Timeout == 0 {
		input.Timeout = waitForApplicationContainerRunningTimeout
	}

	// Wait for application container to be ready and return the result
	applicationContainerInfo, err := c.WaitForApplicationContainerRunning(getInput, input.PollInterval, input.Timeout)
	if err != nil {
		return nil, err
	}

	return applicationContainerInfo, nil
}

// WaitForApplicationContainerRunning waits for an application container to be completely initialized and ready.
func (c *ContainerClient) WaitForApplicationContainerRunning(input *GetApplicationContainerInput, pollInterval, timeoutSeconds time.Duration) (*Container, error) {
	var info *Container
	err := c.client.WaitFor("Waiting for Application container to be ready", pollInterval, timeoutSeconds, func() (bool, error) {
		var getErr error
		info, getErr = c.GetApplicationContainer(input)
		if getErr != nil {
			return false, getErr
		}
		c.client.DebugLogString(fmt.Sprintf("Application Container name is %v, Application Contiainer info is %+v", info.Name, info))
		switch s := info.Status; s {
		case string(applicationContainerStatusRunning): // Target State
			if info.CurrentOnGoingActitvity != "" {
				return false, nil
			}
			c.client.DebugLogString("Application Container Running")
			return true, nil
		case string(applicationContainerStatusNew):
			c.client.DebugLogString("Application Container New")
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown application container state: %s, waiting", s))
			return false, nil
		}
	})
	return info, err
}

// WaitForApplicationContainerDeleted waits for an application container to be fully deleted.
func (c *ContainerClient) WaitForApplicationContainerDeleted(input *DeleteApplicationContainerInput, pollInterval, timeout time.Duration) error {
	return c.client.WaitFor("application container to be deleted", pollInterval, timeout, func() (bool, error) {
		var (
			info *Container
			err  error
		)
		getApplicationContainerInput := &GetApplicationContainerInput{
			Name: input.Name,
		}
		if info, err = c.GetApplicationContainer(getApplicationContainerInput); err != nil {
			if client.WasNotFoundError(err) {
				// Application Container could not be found, thus deleted
				return true, nil
			}
			// Currently the API returns a 400 with a message containing 404 and we check that here
			if strings.Contains(err.Error(), "404") {
				return true, nil
			}
			// Some other error occurred trying to get application container, exit
			return false, err
		}
		switch s := info.Status; s {
		case string(applicationContainerStatusRunning):
			return false, nil
		case string(applicationContainerStatusDestroyPending):
			c.client.DebugLogString(fmt.Sprintf("Destroy pending for application container state: %s, waiting", s))
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown application container state: %s, waiting", s))
			return false, fmt.Errorf("Unknown application container state: %s", s)
		}
	})
}

func readFileContents(filePath string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		fileErr := file.Close()
		if fileErr != nil {
			err = fileErr
		}
	}()

	// Read the file contents
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileContents, nil
}
