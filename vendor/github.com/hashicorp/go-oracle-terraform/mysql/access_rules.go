// Manages the access rules for the MySQL CS Service Instance.
// The only fields that can be updated for an access rule is the desired state: Enabled / Disabled.
// AccessRules are dependent on the existance of ServiceInstance.
package mysql

import (
	"fmt"
	"strings"
	"time"
)

// API URI Paths for the Container and Root objects.
var (
	MySQLAccessRuleContainerPath = "/paas/api/v1.1/instancemgmt/%s/services/MySQLCS/instances/%s/accessrules"
	MySQLAccessRuleRootPath      = "/paas/api/v1.1/instancemgmt/%s/services/MySQLCS/instances/%s/accessrules/%s"
)

// Default Timeout value for Access Rule operations
const WaitForAccessRuleTimeout = time.Duration(30 * time.Second)

// Default polling interval for Access Rule operations
const WaitForAccessRulePollInterval = time.Duration(2 * time.Second)

// AccessRulesClient returns a AccessRulesClient for managing Access Rules for the MySQL CS Service Instance.
func (c *MySQLClient) AccessRules() *AccessRulesClient {
	return &AccessRulesClient{
		AccessRulesResourceClient: AccessRulesResourceClient{
			MySQLClient:      c,
			ContainerPath:    MySQLAccessRuleContainerPath,
			ResourceRootPath: MySQLAccessRuleRootPath,
		},
	}
}

// Status Constants for an Access Rule
type AccessRuleStatus string

const (
	AccessRuleEnabled  AccessRuleStatus = "enabled"
	AccessRuleDisabled AccessRuleStatus = "disabled"
)

// Operations constants for either updating or deleting an access rule.
type AccessRuleOperation string

const (
	AccessRuleUpdate AccessRuleOperation = "update"
	AccessRuleDelete AccessRuleOperation = "delete"
)

// AccessRulesList holds a list of all the AccessRules and a list of Access Rule activities that have
// been performed
type AccessRuleList struct {
	AccessRules []AccessRuleInfo     `json:"accessRules"`
	Activities  []AccessRuleActivity `json:"activities"`
}

// AccessRuleInfo holds the information for a single AccessRule
type AccessRuleInfo struct {
	// Description of the Access Rule
	Description string `json:"description"`
	// The destination of the Access Rule. This is the service object to allow traffic to. e.g. mysql_MASTER
	Destination string `json:"destination,omitempty"`
	// The port(s) to allow traffic to pass through. This can be a single or port range.
	Ports string `json:"ports"`
	// The protocol for the rule. e.g. UDP / TCP
	Protocol string `json:"protocol"`
	// The name of the rule.
	// Required.
	RuleName string `json:"ruleName"`
	// Type of rule. One of "DEFAULT", "SYSTEM" or "USER".
	// Computed value
	RuleType string `json:"ruleType,omitempty"`
	// The hosts which traffic is permitted from. It can be IP address or subnets.
	// Required
	Source string `json:"source"`
	// The status of the rule.
	Status string `json:"status,omitempty"`
}

// AccessRuleActivity describes a single activity operation that have been performed on the access rules.
// Example of activies include created, disabled.
type AccessRuleActivity struct {
	AccessRuleActivityInfo AccessRuleActivityInfo `json:"activity"`
}

// AccessRuleActivityInfo describes an activity/operation that has been performed
type AccessRuleActivityInfo struct {
	// The name of the rule
	RuleName string `json:"ruleName"`
	// The system generated message on the activity that was performed.
	Message string `json:"message"`
	// Error messages that have been reported during the activity.
	Errors string `json:"errors"`
	// The status of the activity
	Status string `json:"status"`
}

// CreateAccessRuleInput defines the input parameters needed to create an access rule.
type CreateAccessRuleInput struct {
	// Name of the MySQL CS Service Instance.
	// Required.
	ServiceInstanceID string `json:"-"`
	// Description of the Access Rule.
	// Required
	Description string `json:"description"`
	// The Destination to which traffic is allowed.
	// Required
	Destination string `json:"destination"`
	// The ports to allow traffic on. Can be a single port or range
	// Required.
	Ports string `json:"ports"`
	// Protocol for the port. Can be tcp or udp
	// Required.
	Protocol string `json:"protocol,omitempty"`
	// Name of the rule.
	// Required.
	RuleName string `json:"ruleName"`
	// The IP Addresses and subnets from which traffic is allowed.
	// Valid values include:
	//   - "PUBLIC-INTERNET" for any host on the internet.
	//   - A single IP address or comma-separated list of subnets (in CIDR format) or IPv4 addresses.
	// Required
	Source string `json:"source"`
	// The desired status of the rule. Either "disabled" or "enabled".
	Status string `json:"status,omitempty"`
	// Time to wait between polling access rule to be ready.
	PollInterval time.Duration `json:"-"`
	// Time to wait for an access rule to be ready.
	Timeout time.Duration `json:"-"`
}

// CreateAccessRule creates an AccessRule with the supplied input.
// The API returns a http 202 on success.
func (c *AccessRulesClient) CreateAccessRule(input *CreateAccessRuleInput) error {

	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	if err := c.createResource(input, nil); err != nil {
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

	getRuleInput := &GetAccessRuleInput{
		ServiceInstanceID: input.ServiceInstanceID,
		Name:              input.RuleName,
	}

	return c.WaitForAccessRuleReady(getRuleInput, pollInterval, timeout)
}

// GetAccessRuleInput defines the input parameters needed to retrieve information on AccessRules
// for a MySQL CS Service Instance.
type GetAccessRuleInput struct {
	// The name of the MySQL CS Service instance. This is used in forming the URI to retrieve the access rules.
	// Required.
	ServiceInstanceID string `json:"-"`
	// Name of the AccessRule.
	// There is no native "GET" to return a single AccessRuleInfo Object, so what we do is get back
	// the full list and search for the access rule locally.
	Name string `json:"-"`
}

// GetAllAccessRules gets all the access rules from a MySQL CS Service instance.
// We make use of the same GetAccessRuleInput, but we ignore the name attribute.
func (c *AccessRulesClient) GetAllAccessRules(input *GetAccessRuleInput) (*AccessRuleList, error) {

	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	var accessRules AccessRuleList
	if err := c.getResource(&accessRules); err != nil {
		return nil, err
	}

	// Iterated through entire slice, rule was not found.
	// No error occurred though, return a nil struct, and allow the Provider to handle
	// a Nil response case.
	return &accessRules, nil
}

// GetAccessRule gets a single access rule info object from the MySQL CS Service Instance.
// The method gets the full list and iterates locally for the matching rule name.
func (c *AccessRulesClient) GetAccessRule(input *GetAccessRuleInput) (*AccessRuleInfo, error) {

	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	var accessRules AccessRuleList
	if err := c.getResource(&accessRules); err != nil {
		return nil, err
	}

	// This is likely not the most optimal path for this, however, the upper bound on
	// performance here is the actual API request, not the iteration.
	for _, rule := range accessRules.AccessRules {
		if rule.RuleName == input.Name {
			return &rule, nil
		}
	}

	// Iterated through entire slice, rule was not found.
	// No error occurred though, return a nil struct, and allow the Provider to handle
	// a Nil response case.
	return nil, nil
}

// WaitForAccessRuleReady gets into a wait loop for access rule to be created successfully and available.
// The creation typically takes some time before the rule is available, so we get into a wait loop until
// the access rule is ready.
func (c *AccessRulesClient) WaitForAccessRuleReady(input *GetAccessRuleInput, pollInterval time.Duration, timeoutSeconds time.Duration) error {

	err := c.client.WaitFor("access rule to be created.", pollInterval, timeoutSeconds, func() (bool, error) {

		var info AccessRuleList
		if err := c.getResource(&info); err != nil {
			return false, err
		}

		c.client.DebugLogString(fmt.Sprintf("[DEBUG] Checking Activities : %v", info))
		for _, accessRule := range info.AccessRules {
			if accessRule.RuleName == input.Name {
				return true, nil
			}
		}

		for _, activity := range info.Activities {
			if activity.AccessRuleActivityInfo.RuleName == input.Name {
				switch s := strings.ToUpper(activity.AccessRuleActivityInfo.Status); s {
				case "FAILED":
					c.client.DebugLogString(fmt.Sprintf("AccessRule [%s] is FAILED. Status : %s", activity.AccessRuleActivityInfo.RuleName, activity.AccessRuleActivityInfo.Status))
					return false, fmt.Errorf("Error creating Access Rule : %s", activity.AccessRuleActivityInfo.Message)
				case "SUCCESS":
					c.client.DebugLogString(fmt.Sprintf("AccessRule [%s] state is SUCCESS. Status : %s", activity.AccessRuleActivityInfo.RuleName, activity.AccessRuleActivityInfo.Status))
					return true, nil
				case "RUNNING":
					c.client.DebugLogString(fmt.Sprintf("AccessRule [%s] state is RUNNING. Status : %s", activity.AccessRuleActivityInfo.RuleName, activity.AccessRuleActivityInfo.Status))
					return false, nil
				default:
					c.client.DebugLogString(fmt.Sprintf("AccessRule [%s] state is DEFAULT. Status : %s", activity.AccessRuleActivityInfo.RuleName, activity.AccessRuleActivityInfo.Status))
					return false, nil
				}
			}
		}
		return false, nil
	})
	return err
}

// UpdateAccessRuleInput defines the Update parameters needed to update an
// AccessRule for the MySQL CS Service Instance
type UpdateAccessRuleInput struct {
	// Name of the MySQL CS Service Instance. This is used in the request URI.
	// Required.
	ServiceInstanceID string `json:"-"`
	// Name of the Access Rule. This is used in the request URI.
	// Required.
	Name string `json:"-"`
	// The type of operation being performed. The value is implicit and will
	// be set when calling the Update Access rule method.
	// Do not set.
	Operation AccessRuleOperation `json:"operation"`
	// Desired status of the Access Rule. This is the only attribute that can
	// actually be modified on an access rule.
	// Required
	Status AccessRuleStatus `json:"status"`
	// Sets the time to wait before checking the status of the update. This
	// attribute is implicit and is set by the default pollinterval for the
	// provider.
	PollInterval time.Duration `json:"-"`
	// Time to wait for the update to be ready.
	Timeout time.Duration `json:"-"`
}

// UpdateAccessRule updates an AccessRule with the provided input struct. Returns a fully populated Info struct
// and any errors encountered
func (c *AccessRulesClient) UpdateAccessRule(input *UpdateAccessRuleInput) (*AccessRuleInfo, error) {
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

// DeleteAccessRuleInput defines the Delete parameters to delete an access rule on the MySQL CS service instance.
type DeleteAccessRuleInput struct {
	// Name of the MySQL CS Service Instance.
	// Required.
	ServiceInstanceID string `json:"-"`
	// Name of the AccessRule.
	// Required.
	Name string `json:"-"`
	// Type of Operation being performed. This should not be set.
	// The SDK will set the operation to use the constant defined for delete.
	// Do not set.
	Operation AccessRuleOperation `json:"operation"`
	// The Desired staut sof the Access Rule. Because we are calling a delete, this attribute is
	// ignored.
	Status AccessRuleStatus `json:"status"`
	// Time to wait between checking for access rule state.
	PollInterval time.Duration `json:"-"`
	// Time to wait for an access rule to be deleted completely.
	Timeout time.Duration `json:"-"`
}

// DeleteAccessRule deletes an AccessRule with the provided input struct. Returns any errors that occurred.
func (c *AccessRulesClient) DeleteAccessRule(input *DeleteAccessRuleInput) error {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	c.client.DebugLogString(fmt.Sprintf("[DEBUG] Deleting AccessRule : %s", input.Name))

	// Since this is strictly an Update call, set the Operation constant
	input.Operation = AccessRuleDelete
	// The Update API call with a `DELETE` operation actually returns the same access rule info
	// in a response body. As we are deleting the AccessRule, we don't actually need to parse that.
	// However, the Update API call requires a pointer to parse, or else we throw an error during the
	// json unmarshal
	var result AccessRuleInfo
	if err := c.updateResource(input.Name, input, &result); err != nil {
		c.client.DebugLogString(fmt.Sprintf("[DEBUG] Failed to delete access rule : %v", err))
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
		c.client.DebugLogString(fmt.Sprintf("[DEBUG] Failed to delete access rule : %v", err))
		return err
	}

	return nil
}

// WaitForAccessRuleDeleted waits for the access rule to be delete completely. As the operations are asynchronous, we invoke the
// delete an poll the API to check that the AccessRule is completely removed from the access rule list.
func (c *AccessRulesClient) WaitForAccessRuleDeleted(input *GetAccessRuleInput, pollInterval time.Duration, timeout time.Duration) (*AccessRuleInfo, error) {
	var info *AccessRuleInfo
	//var getErr error
	err := c.client.WaitFor("access rule to be deleted", pollInterval, timeout, func() (bool, error) {
		var info AccessRuleList
		if err := c.getResource(&info); err != nil {
			return false, err
		}

		// First level check is to see if the access rule name is still in the access rule list.
		// If its still on the list, the operation has not completed.
		for _, accessRule := range info.AccessRules {
			if accessRule.RuleName == input.Name {
				return false, nil
			}
		}

		// If the access rule is no longer on the access rule list, we check if the activity status shows that the delete operation
		// is still running. If the operation is still running, it blocks some other operations (like delete service instance) from
		// being completed successfully.
		for _, activity := range info.Activities {
			if activity.AccessRuleActivityInfo.RuleName == input.Name {
				switch s := strings.ToUpper(activity.AccessRuleActivityInfo.Status); s {
				case "FAILED":
					c.client.DebugLogString(fmt.Sprintf("AccessRule [%s] state is FAILED. Status : %s", activity.AccessRuleActivityInfo.RuleName, activity.AccessRuleActivityInfo.Status))
					return false, fmt.Errorf("Error deleting Access Rule : %s", activity.AccessRuleActivityInfo.Message)
				case "SUCCESS":
					c.client.DebugLogString(fmt.Sprintf("AccessRule [%s] state is SUCCESS. Status : %s", activity.AccessRuleActivityInfo.RuleName, activity.AccessRuleActivityInfo.Status))
					return true, nil
				case "RUNNING":
					c.client.DebugLogString(fmt.Sprintf("AccessRule [%s] state is RUNNING. Status : %s", activity.AccessRuleActivityInfo.RuleName, activity.AccessRuleActivityInfo.Status))
					return false, nil
				default:
					c.client.DebugLogString(fmt.Sprintf("AccessRule [%s] state is DEFAULT. Status : %s", activity.AccessRuleActivityInfo.RuleName, activity.AccessRuleActivityInfo.Status))
					return false, nil
				}
			}
		}

		// IF it reaches this stage, it means there is that the access rule is not in the returned result, and there are
		// no more activities, so we can safely assume everything is ok.
		return true, nil
	})
	return info, err
}
