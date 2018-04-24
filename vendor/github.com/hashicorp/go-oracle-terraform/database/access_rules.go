// Manages Access Rules for a DBaaS Service Instance.
// The only fields that can be Updated for an Access Rule is the desired state
// of the access rule. From Enabled -> Disabled.
// Deleting an Access Rule also requires an Update call, instead of a Delete API request,
// but the Operation body parameter changes from `update` to `delete`.
// All other parameters for the resource, aside from Status should be ForceNew.
// The READ function for the AccessRule resource is tricky, as there is
// There is an API endpoint to view "all" rules, however, which will be used as a
// data source to match on a supplied AccessRule name.
// Timeout only supported for the CREATE method

package database

import (
	"time"

	"github.com/hashicorp/go-oracle-terraform/helper"
)

// API URI Paths for Container and Root objects
const (
	DBAccessContainerPath = "/paas/api/v1.1/instancemgmt/%s/services/dbaas/instances/%s/accessrules"
	DBAccessRootPath      = "/paas/api/v1.1/instancemgmt/%s/services/dbaas/instances/%s/accessrules/%s"
)

// Default Timeout value for Create
const WaitForAccessRuleTimeout = time.Duration(10 * time.Second)

// Default Poll Interval value for Create
const WaitForAccessRulePollInterval = time.Duration(1 * time.Second)

// AccessRules returns a UtilityClient for managing SSH Keys and Access Rules for a DBaaS Service Instance
func (c *DatabaseClient) AccessRules() *UtilityClient {
	return &UtilityClient{
		UtilityResourceClient: UtilityResourceClient{
			DatabaseClient:   c,
			ContainerPath:    DBAccessContainerPath,
			ResourceRootPath: DBAccessRootPath,
		},
	}
}

// Status Constants for an Access Rule
type AccessRuleStatus string

const (
	AccessRuleEnabled  AccessRuleStatus = "enabled"
	AccessRuleDisabled AccessRuleStatus = "disabled"
)

// Operational Constants for either Updating/Deleting an Access Rule
type AccessRuleOperation string

const (
	AccessRuleUpdate AccessRuleOperation = "update"
	AccessRuleDelete AccessRuleOperation = "delete"
)

// Default Destination for an Access Rule
type AccessRuleDestination string

const (
	AccessRuleDefaultDestination AccessRuleDestination = "DB_1"
)

// Used for the GET request, as there's no direct GET request for a single Access Rule
type AccessRules struct {
	Rules []AccessRuleInfo `json:"accessRules"`
}

type AccessRuleType string

const (
	AccessRuleTypeDefault AccessRuleType = "DEFAULT"
	AccessRuleTypeSystem  AccessRuleType = "SYSTEM"
	AccessRuleTypeUser    AccessRuleType = "USER"
)

// AccessRuleInfo holds all of the known information for a single AccessRule
type AccessRuleInfo struct {
	// The Description of the Access Rule
	Description string `json:"description"`
	// The destination of the Access Rule. Should always be "DB".
	Destination AccessRuleDestination `json:"destination"`
	// The ports for the rule.
	Ports string `json:"ports"`
	// The name of the Access Rule
	Name string `json:"ruleName"`
	// The Type of the rule. One of: "DEFAULT", "SYSTEM", or "USER".
	// Computed Value
	RuleType AccessRuleType `json:"ruleType"`
	// The IP Addresses and subnets from which traffic is allowed
	Source string `json:"source"`
	// The current status of the Access Rule
	Status AccessRuleStatus `json:"status"`
}

// CreateAccessRuleInput defines the input parameters needed to create an Access Rule for a DBaaS Service Instance.
type CreateAccessRuleInput struct {
	// Name of the DBaaS service instance.
	// Required
	ServiceInstanceID string `json:"-"`
	// Description of the Access Rule.
	// Required
	Description string `json:"description"`
	// Destination to which traffic is allowed. Specify the value "DB".
	// Required
	Destination AccessRuleDestination `json:"destination"`
	// The network port or ports to allow traffic on. Specified as a single port or a range.
	// Required
	Ports string `json:"ports"`
	// The name of the Access Rule
	// Required
	Name string `json:"ruleName"`
	// The IP addresses and subnets from which traffic is allowed.
	// Valid values are:
	//   - "DB" for any other cloud service instance in the service instances `ora_db` security list
	//   - "PUBLIC-INTERNET" for any host on the internet.
	//   - A single IP address or comma-separated list of subnets (in CIDR format) or IPv4 addresses.
	// Required
	Source string `json:"source"`
	// Desired Status of the rule. Either "disabled" or "enabled".
	// Required
	Status AccessRuleStatus `json:"status"`
	// Time to wait between polling for access rule to be ready
	PollInterval time.Duration `json:"-"`
	// Time to wait for an access rule to be ready
	Timeout time.Duration `json:"-"`
}

// Creates an AccessRule with the supplied input struct.
// The API call to Create returns a nil body object, and a 202 status code on success.
// Thus, the Create method will return the resulting object from an internal GET call
// during the WaitForReady timeout.
func (c *UtilityClient) CreateAccessRule(input *CreateAccessRuleInput) (*AccessRuleInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	var accessRule AccessRuleInfo
	if err := c.createResource(input, &accessRule); err != nil {
		return nil, err
	}

	pollInterval := input.PollInterval
	if pollInterval == 0 {
		pollInterval = WaitForAccessRulePollInterval
	}

	timeout := input.Timeout
	if timeout == 0 {
		timeout = WaitForAccessRuleTimeout
	}

	getInput := &GetAccessRuleInput{
		Name: input.Name,
	}

	result, err := c.WaitForAccessRuleReady(getInput, pollInterval, timeout)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAccessRuleInput defines the input parameters needed to retrieve information
// on an AccessRule for a DBaas Service Instance.
type GetAccessRuleInput struct {
	// Name of the DBaaS service instance.
	// Required
	ServiceInstanceID string `json:"-"`
	// Name of the Access Rule.
	// Because there is no native "GET" to return a single AccessRuleInfo object, we don't
	// need to marshal a request body for the GET request. Instead the request returns a slice
	// of AccessRuleInfo structs, which we iterate on to interpret the desired AccessRuleInfo struct
	// Required
	Name string `json:"-"`
}

// Get's a slice of every AccessRule, and iterates on the result until
// we find the correctly matching access rule. This is likely an expensive operation depending
// on how many access rules the customer has. However, since there's no direct GET API endpoint
// for a single Access Rule, it's not able to be optimized yet.
func (c *UtilityClient) GetAccessRule(input *GetAccessRuleInput) (*AccessRuleInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	var accessRules AccessRules
	if err := c.getResource("", &accessRules); err != nil {
		return nil, err
	}

	// This is likely not the most optimal path for this, however, the upper bound on
	// performance here is the actual API request, not the iteration.
	for _, rule := range accessRules.Rules {
		if rule.Name == input.Name {
			return &rule, nil
		}
	}

	// Iterated through entire slice, rule was not found.
	// No error occured though, return a nil struct, and allow the Provdier to handle
	// a Nil response case.
	return nil, nil
}

// UpdateAccessRuleInput defines the Update parameters needed to update an AccessRule
// for a DBaaS Service Instance.
type UpdateAccessRuleInput struct {
	// Name of the DBaaS Service Instance.
	// Required
	ServiceInstanceID string `json:"-"`
	// Name of the Access Rule. Used in the request's URI, not as a body parameter.
	// Required
	Name string `json:"-"`
	// Type of Operation being performed. This should never be set in the Provider,
	// as we're explicitly calling an Update function here, so the SDK uses the constant
	// defined for Updating an AccessRule
	// Do not set.
	Operation AccessRuleOperation `json:"operation"`
	// Desired Status of the Access Rule. This is the only attribute that can actually be
	// modified on an access rule.
	// Required
	Status AccessRuleStatus `json:"status"`
}

// Updates an AccessRule with the provided input struct. Returns a fully populated Info struct
// and any errors encountered
func (c *UtilityClient) UpdateAccessRule(input *UpdateAccessRuleInput,
) (*AccessRuleInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	// Since this is strictly an Update call, set the Operation constant
	input.Operation = AccessRuleUpdate
	// Initialize the response struct
	var accessRule AccessRuleInfo
	if err := c.updateResource(input.Name, input, &accessRule); err != nil {
		return nil, err
	}
	return &accessRule, nil
}

// DeleteAccessRuleInput defines the Delete parameters needed to delete an AccessRule
// for a DBaaS Service Instance. There's no dedicated DELETE method on the API, so this
// mimics the same behavior of the Update method, but using the Delete operational constant.
// Instead of implementing, choosing to be verbose here for ease of use in the Provider, and clarity.
type DeleteAccessRuleInput struct {
	// Name of the DBaaS Service Instance.
	// Required
	ServiceInstanceID string `json:"-"`
	// Name of the Access Rule. Used in the request's URI, not as a body parameter.
	// Required
	Name string `json:"-"`
	// Type of Operation being performed. This should never be set in the Provider,
	// as we're explicitly calling an Delete function here, so the SDK uses the constant
	// defined for Deleting an AccessRule
	// Do not set.
	Operation AccessRuleOperation `json:"operation"`
	// Desired Status of the Access Rule. This is the only attribute that can actually be
	// modified on an access rule.
	// Required
	Status AccessRuleStatus `json:"status"`
	// Time to wait between checking access rule state
	PollInterval time.Duration `json:"-"`
	// Time to wait for an access rule to be ready
	Timeout time.Duration `json:"-"`
}

// Deletes an AccessRule with the provided input struct. Returns any errors that occurred.
func (c *UtilityClient) DeleteAccessRule(input *DeleteAccessRuleInput) error {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	// Since this is strictly an Update call, set the Operation constant
	input.Operation = AccessRuleDelete
	// The Update API call with a `DELETE` operation actually returns the same access rule info
	// in a response body. As we are deleting the AccessRule, we don't actually need to parse that.
	// However, the Update API call requires a pointer to parse, or else we throw an error during the
	// json unmarshal
	var result AccessRuleInfo
	if err := c.updateResource(input.Name, input, &result); err != nil {
		return err
	}

	pollInterval := input.PollInterval
	if pollInterval == 0 {
		pollInterval = WaitForAccessRulePollInterval
	}

	timeout := input.Timeout
	if timeout == 0 {
		timeout = WaitForAccessRuleTimeout
	}

	getInput := &GetAccessRuleInput{
		Name: input.Name,
	}

	_, err := c.WaitForAccessRuleDeleted(getInput, pollInterval, timeout)
	if err != nil {
		return err
	}

	return nil
}

func (c *UtilityClient) WaitForAccessRuleReady(input *GetAccessRuleInput, pollInterval, timeout time.Duration) (*AccessRuleInfo, error) {
	var info *AccessRuleInfo
	var getErr error
	err := c.client.WaitFor("access rule to be ready", pollInterval, timeout, func() (bool, error) {
		info, getErr = c.GetAccessRule(input)
		if getErr != nil {
			return false, getErr
		}
		if info != nil {
			// Rule found, return. Desired case
			return true, nil
		}
		// Rule not found, wait
		return false, nil
	})
	return info, err
}

func (c *UtilityClient) WaitForAccessRuleDeleted(input *GetAccessRuleInput, pollInternval, timeout time.Duration) (*AccessRuleInfo, error) {
	var info *AccessRuleInfo
	var getErr error
	err := c.client.WaitFor("access rule to be deleted", pollInternval, timeout, func() (bool, error) {
		info, getErr = c.GetAccessRule(input)
		if getErr != nil {
			return true, nil
		}
		if info != nil {
			// Rule found, continue
			return false, nil
		}
		// Rule not found, return. Desired case
		return true, nil
	})
	return info, err
}

type DefaultAccessRuleInfo struct {
	// Name of the DBaaS service instance.
	ServiceInstanceID string `json:"-"`
	// Enabled for every DB Service Instance
	EnableSSH *bool
	// Single Instance Rules
	EnableHTTP       *bool
	EnableHTTPSSL    *bool
	EnableDBConsole  *bool
	EnableDBExpress  *bool
	EnableDBListener *bool
	// RAC Rules
	EnableEMConsole     *bool
	EnableRACDBListener *bool
	EnableScanListener  *bool
	EnableRACOns        *bool
}

// Default Access Rule prefixes
var DefaultAccessRuleNames = map[string]string{
	"EnableSSH":           "ora_p2_ssh",
	"EnableHTTP":          "ora_p2_http",
	"EnableHTTPSSL":       "ora_p2_httpssl",
	"EnableDBConsole":     "ora_p2_dbconsole",
	"EnableDBExpress":     "ora_p2_dbexpress",
	"EnableDBListener":    "ora_p2_dblistener",
	"EnableEMConsole":     "ora_p2_emconsole",
	"EnableRACDBListener": "ora_p2_db_listener",
	"EnableScanListener":  "ora_p2_scan_listener",
	"EnableRACOns":        "ora_p2_ons",
}

// GetDefaultAccessRuleInput defines the input parameters needed to retrieve information
// on an all the DefaultAccessRule for a DBaas Service Instance.
type GetDefaultAccessRuleInput struct {
	// Name of the DBaaS service instance.
	// Required
	ServiceInstanceID string `json:"-"`
}

// GetDefaultAccessRules retrieves all the default access rules pertaining to Database Service Instance
func (c *UtilityClient) GetDefaultAccessRules(input *GetDefaultAccessRuleInput) (*DefaultAccessRuleInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}
	defaultAccessRules := &DefaultAccessRuleInfo{}
	// Obtain all the access rules since it isn't possible to get a specific one from the api
	var accessRules AccessRules
	if err := c.getResource("", &accessRules); err != nil {
		return nil, err
	}
	for key, ruleName := range DefaultAccessRuleNames {
		// Iterate through AccessRules to get the one we are looking for.
		// Not optimal but it's a limitation on the api.
		var rule *AccessRuleInfo
		for _, accessRule := range accessRules.Rules {
			if ruleName == accessRule.Name {
				rule = &accessRule
				break
			}
		}

		if rule != nil {
			if key == "EnableSSH" {
				defaultAccessRules.EnableSSH = helper.Bool(rule.Status == AccessRuleEnabled)
			}
			if key == "EnableHTTP" {
				defaultAccessRules.EnableHTTP = helper.Bool(rule.Status == AccessRuleEnabled)
			}
			if key == "EnableHTTPSSL" {
				defaultAccessRules.EnableHTTPSSL = helper.Bool(rule.Status == AccessRuleEnabled)
			}
			if key == "EnableDBConsole" {
				defaultAccessRules.EnableDBConsole = helper.Bool(rule.Status == AccessRuleEnabled)
			}
			if key == "EnableDBExpress" {
				defaultAccessRules.EnableDBExpress = helper.Bool(rule.Status == AccessRuleEnabled)
			}
			if key == "EnableDBListener" {
				defaultAccessRules.EnableDBListener = helper.Bool(rule.Status == AccessRuleEnabled)
			}
			if key == "EnableEMConsole" {
				defaultAccessRules.EnableEMConsole = helper.Bool(rule.Status == AccessRuleEnabled)
			}
			if key == "EnableRACDBListener" {
				defaultAccessRules.EnableRACDBListener = helper.Bool(rule.Status == AccessRuleEnabled)
			}
			if key == "EnableScanListener" {
				defaultAccessRules.EnableScanListener = helper.Bool(rule.Status == AccessRuleEnabled)
			}
			if key == "EnableRACOns" {
				defaultAccessRules.EnableRACOns = helper.Bool(rule.Status == AccessRuleEnabled)
			}
		}
		rule = nil
	}
	return defaultAccessRules, nil
}

// UpdateDefaultAccessRules Updates all the specified/relevant default access rules for a database service instance
func (c *UtilityClient) UpdateDefaultAccessRules(input *DefaultAccessRuleInfo) (*DefaultAccessRuleInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}
	var accessRules AccessRules
	if err := c.getResource("", &accessRules); err != nil {
		return nil, err
	}
	for key, ruleName := range DefaultAccessRuleNames {
		if key == "EnableSSH" && input.EnableSSH != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableSSH)
			if err != nil {
				return nil, err
			}
		}
		if key == "EnableHTTP" && input.EnableHTTP != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableHTTP)
			if err != nil {
				return nil, err
			}
		}
		if key == "EnableHTTPSSL" && input.EnableHTTPSSL != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableHTTPSSL)
			if err != nil {
				return nil, err
			}
		}
		if key == "EnableDBConsole" && input.EnableDBConsole != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableDBConsole)
			if err != nil {
				return nil, err
			}
		}
		if key == "EnableDBExpress" && input.EnableDBExpress != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableDBExpress)
			if err != nil {
				return nil, err
			}
		}
		if key == "EnableDBListener" && input.EnableDBListener != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableDBListener)
			if err != nil {
				return nil, err
			}
		}
		if key == "EnableEMConsole" && input.EnableEMConsole != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableEMConsole)
			if err != nil {
				return nil, err
			}
		}
		if key == "EnableRACDBListener" && input.EnableRACDBListener != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableRACDBListener)
			if err != nil {
				return nil, err
			}
		}
		if key == "EnableScanListener" && input.EnableScanListener != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableScanListener)
			if err != nil {
				return nil, err
			}
		}
		if key == "EnableRACOns" && input.EnableRACOns != nil {
			err := updateDefaultAccessRule(c, accessRules, ruleName, input.ServiceInstanceID, *input.EnableRACOns)
			if err != nil {
				return nil, err
			}
		}
	}
	getInput := &GetDefaultAccessRuleInput{
		ServiceInstanceID: input.ServiceInstanceID,
	}
	defaultAccessRules, err := c.GetDefaultAccessRules(getInput)
	if err != nil {
		return nil, err
	}
	defaultAccessRules.ServiceInstanceID = input.ServiceInstanceID
	return defaultAccessRules, nil
}

// Updates a specific Default Access Rule if it's status differs from the requested status
func updateDefaultAccessRule(c *UtilityClient, accessRules AccessRules, ruleName, serviceInstanceID string, enabled bool) error {
	var rule *AccessRuleInfo
	for _, accessRule := range accessRules.Rules {
		if ruleName == accessRule.Name {
			rule = &accessRule
			break
		}
	}
	if rule != nil {
		var status AccessRuleStatus
		if enabled {
			status = AccessRuleEnabled
		} else {
			status = AccessRuleDisabled
		}
		if rule.Status != status {
			updateRuleInput := &UpdateAccessRuleInput{
				ServiceInstanceID: serviceInstanceID,
				Name:              rule.Name,
				Status:            status,
			}
			_, err := c.UpdateAccessRule(updateRuleInput)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
