package mysql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// ResourceClient is an AuthenticatedClient with some additional information about the resources to be addressed.

type AccessRulesClient struct {
	AccessRulesResourceClient
}

type AccessRulesResourceClient struct {
	*MySQLClient
	ResourceDescription string
	ContainerPath       string
	ResourceRootPath    string
	ServiceInstanceID   string
}

func (c *AccessRulesResourceClient) createResource(requestBody interface{}, responseBody interface{}) error {

	var objectPath = c.getContainerPath(c.ContainerPath)

	_, err := c.executeRequestWithContentType("POST", objectPath, requestBody, "application/json")

	if err != nil {
		return err
	}

	return nil
}

func (c *AccessRulesResourceClient) getResource(responseBody interface{}) error {

	var objectPath = c.getContainerPath(c.ContainerPath)

	resp, err := c.executeRequestWithContentType("GET", objectPath, nil, "application/json")
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *AccessRulesResourceClient) updateResource(name string, requestBody interface{}, responseBody interface{}) error {

	resp, err := c.executeRequestWithContentType("PUT", c.getObjectPath(c.ResourceRootPath, name), requestBody, "application/json")
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *AccessRulesResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	c.client.DebugLogString(fmt.Sprintf("[Debug] : HTTP Resp (%d): %v", resp.StatusCode, buf))
	// JSON decode response into interface
	var tmp interface{}
	dcd := json.NewDecoder(buf)
	if err := dcd.Decode(&tmp); err != nil {
		return fmt.Errorf("%+v", resp)
	}

	// Use mapstructure to weakly decode into the resulting interface
	msdcd, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           iface,
		TagName:          "json",
	})
	if err != nil {
		return err
	}

	if err := msdcd.Decode(tmp); err != nil {
		return err
	}
	return nil
}

func (c *AccessRulesResourceClient) getContainerPath(root string) string {
	// /paas/api/v1.1/instancemgmt/{identityDomainId}/services/MySQLCS/instances/{serviceId}/accessrules
	c.client.DebugLogString(fmt.Sprintf("[DEBUG] getAccessRuleObjectPath : %s / %s", *c.client.IdentityDomain, c.ServiceInstanceID))
	return fmt.Sprintf(root, *c.client.IdentityDomain, c.ServiceInstanceID)
}

func (c *AccessRulesResourceClient) getObjectPath(root, name string) string {
	// /paas/api/v1.1/instancemgmt/{identityDomainId}/services/MySQLCS/instances/{serviceId}/accessrules/{ruleName}
	c.client.DebugLogString(fmt.Sprintf("[DEBUG] getAccessRuleObjectPath : %v / %s / %s", c.client.IdentityDomain, c.ServiceInstanceID, name))
	return fmt.Sprintf(root, *c.client.IdentityDomain, c.ServiceInstanceID, name)

}
