package application

import (
	"fmt"
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

// Container container information about the application container
type Container struct {
	// ID of the application
	AppID string `json:"appId"`
	// URL of the created application
	AppURL string `json:"appURL"`
	// Creation time of the application
	CreatedTime string `json:"createdTime"`
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
	// Name of the manifest file, required if this file is not packaged with the application
	// Optional
	Manifest string
	// Time to wait between checks on application container status
	PollInterval time.Duration
	// Timeout for creating an application container
	Timeout time.Duration
}

// CreateApplicationContainerAdditionalFields specifies the additional fields needed to create an application container
type CreateApplicationContainerAdditionalFields struct {
	// Location of the application archive file in Oracle Storage Cloud Service, in the format app-name/file-name
	// Optional
	ArchiveURL string
	// Uses Oracle Identity Cloud Service to control who can access your Java SE 7 or 8, Node.js, or PHP application.
	// This should be ApplicationContainerAuthType but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Allowed values are 'basic' and 'oauth'.
	// Optional
	AuthType string // ApplicationContainerAuthType
	// Name of the application
	// Required
	Name string
	// Comments on the application deployment
	Notes string
	// Email address to which application deployment status updates are sent.
	NotificationEmail string
	// Repository of the application. The only allowed value is 'dockerhub'.
	// This should be ApplicationRepository but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Optional
	Repository string // ApplicationRepository
	// Runtime environment: java (the default), node, php, python, or ruby
	// This should be ApplicationRuntime but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Required
	Runtime string // ApplicationRuntime
	// Subscription, either hourly (the default) or monthly
	// This should be ApplicationSubscriptionType but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Optional
	SubscriptionType string //ApplicationSubscriptionType
}

// CreateApplicationContainer creates a new Application Container from an ApplicationClient and an input struct.
// Returns a populated ApplicationContainer struct for the Application, and any errors
func (c *ContainerClient) CreateApplicationContainer(input *CreateApplicationContainerInput) (*Container, error) {

	files := make(map[string]string)
	if input.Deployment != "" {
		files["deployment"] = input.Deployment
	}
	if input.Manifest != "" {
		files["manifest"] = input.Manifest
	}
	additionalFields := structs.Map(input.AdditionalFields)

	var applicationContainer *Container
	if err := c.createResource(files, additionalFields, applicationContainer); err != nil {
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
	// Name of the manifest file, required if this file is not packaged with the application
	// Optional
	Manifest string
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

	files := make(map[string]string)
	if input.Deployment != "" {
		files["deployment"] = input.Deployment
	}
	if input.Manifest != "" {
		files["manifest"] = input.Manifest
	}
	additionalFields := structs.Map(input.AdditionalFields)

	var applicationContainer *Container
	if err := c.updateResource(files, additionalFields, applicationContainer); err != nil {
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
