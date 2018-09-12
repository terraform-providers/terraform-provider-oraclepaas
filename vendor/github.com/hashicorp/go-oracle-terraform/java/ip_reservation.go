package java

import (
	"fmt"
	"time"
)

// API URI Paths for Container and Root objects
const (
	ipReservationContainerPath = "/paas/api/v1.1/network/%s/services/jaas/ipreservations"
	ipReservationResourcePath  = "/paas/api/v1.1/network/%s/services/jaas/ipreservations/%s"
)

const waitForIPReservationReadyPollInterval = 1 * time.Second
const waitForIPReservationReadyTimeout = 5 * time.Minute

// IPReservationClient IP Reservation API Client
type IPReservationClient struct {
	IPReservationResourceClient
	PollInterval time.Duration
	Timeout      time.Duration
}

// IPReservationStatus IP Reservation Status values
type IPReservationStatus string

const (
	// IPReservationStatusInitializing Initializing
	IPReservationStatusInitializing IPReservationStatus = "INITIALIZING"
	// IPReservationStatusUnused Unused IP Reservation
	IPReservationStatusUnused IPReservationStatus = "UNUSED"
	// IPReservationStatusUsed Used IP Reservation
	IPReservationStatusUsed IPReservationStatus = "USED"
)

// IPReservationClient obtains an new ResourceClient which can be used to access the
// Database Cloud IP Reservation API
func (c *Client) IPReservationClient() *IPReservationClient {
	return &IPReservationClient{
		IPReservationResourceClient: IPReservationResourceClient{
			Client:           c,
			ContainerPath:    ipReservationContainerPath,
			ResourceRootPath: ipReservationResourcePath,
		},
	}
}

// CreateIPReservationInput represents the Create IP Reservation API Request body
type CreateIPReservationInput struct {
	// Identity domain ID for the Database Cloud Service account
	// For a Cloud account with Identity Cloud Service: the identity service ID, which has the form idcs-letters-and-numbers.
	// For a traditional cloud account: the name of the identity domain.
	// Required
	IdentityDomainID string
	// Name of the IP reservation to create.
	// Required
	Name string `json:"ipResName"`
	// Indicates whether the IP reservation is for instances attached to IP networks or the shared network
	// set to `IPNetwork` for IP Network, or omit for shared network
	NetworkType string `json:"networkType,omitempty"`
	// Name of the region to create the IP reservation in
	// Required
	Region string `json:"region"`
}

// CreateIPReservationInfo represents the Create IP Reservation API Response
type CreateIPReservationInfo struct {
	// Name of the IP reservation to create.
	Name string `json:"ipResName"`
	// Location (region) of the IP reservation.
	ComputeSiteName string `json:"computeSite"`
	// Create IP Reservation Job ID
	JobID string `json:"jobId"`
}

// DeleteIPReservationInfo represents the Delete IP Reservation API Response
type DeleteIPReservationInfo struct {
	// Name of the IP reservation to create.
	Name string `json:"ipResName"`
	// Location (region) of the IP reservation.
	ComputeSiteName string `json:"computeSite"`
	// Create IP Reservation Job ID
	JobID string `json:"jobId"`
}

// IPReservationInfo represents the Get IP Reservation API Response
type IPReservationInfo struct {
	// Id of the IP reservation.
	ID int `json:"id"`
	// Name of the IP reservation.
	Name string `json:"name"`
	// Location (region) of the IP reservation.
	ComputeSiteName string `json:"computeSiteName"`
	// Name of the compute node using the IP reservation.
	// This parameter is returned only when the IP reservation is in use.
	Hostname string `json:"hostName"`
	// The identity domain ID of the IP reservation.
	IdentityDomain string `json:"identityDomain"`
	// The public IP address for the IP reservation.
	IPAddress string `json:"ipAddress"`
	// Indicates whether the IP reservation is intended for instances attached to IP networks or the shared network.
	NetworkType string `json:"networkType"`
	// The service entitlement Id of Database Cloud Service within the Cloud account.
	ServiceEntitlementID string `json:"serviceEntitlementId"`
	// Name of the Database Cloud Service instance where the named IP reservation is used.
	// This parameter is returned only when the IP reservation is in use.
	ServiceName string `json:"serviceName"`
	// The Service Type the IP Reservation is valid for. `DBaaS`.
	ServiceType string `json:"serviceType"`
	// Status of the IP reservation. Valid values: `INITIALIZING`, `UNUSED`, `USED`.
	Status IPReservationStatus `json:"status"`
}

// IPReservations - used for the GET request that returns all reservations
type IPReservations struct {
	IPReservations []IPReservationInfo `json:"ipReservations"`
}

// CreateIPReservation creates a new IP Reservation.
func (c *IPReservationClient) CreateIPReservation(input *CreateIPReservationInput) (*IPReservationInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForIPReservationReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForIPReservationReadyTimeout
	}

	info, err := c.createResource(input)
	if err != nil {
		return nil, err
	}

	getJobInput := &GetJobInput{
		ID: info.JobID,
	}

	err = c.Client.Jobs().WaitForJobCompletion(getJobInput, c.PollInterval, c.Timeout)
	if err != nil {
		return nil, fmt.Errorf("error creating IP Reservation %q: %+v", input.Name, err)
	}

	ipReservation, err := c.GetIPReservation(input.Name)
	return ipReservation, err
}

// GetIPReservation get the details of an IP Reservation.
func (c *IPReservationClient) GetIPReservation(name string) (*IPReservationInfo, error) {
	info, err := c.getResource(name)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// DeleteIPReservation deletes an IP Reservation.
func (c *IPReservationClient) DeleteIPReservation(name string) error {
	info, err := c.deleteResource(name)
	if err != nil {
		return err
	}

	getJobInput := &GetJobInput{
		ID: info.JobID,
	}

	err = c.Client.Jobs().WaitForJobCompletion(getJobInput, c.PollInterval, c.Timeout)
	if err != nil {
		return fmt.Errorf("error destroying IP Reservation %q: %+v", name, err)
	}
	return nil
}
