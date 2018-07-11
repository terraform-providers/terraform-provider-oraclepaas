package mysql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

// ResourceClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type ResourceClient struct {
	*MySQLClient
	ResourceDescription string
	ContainerPath       string
	ResourceRootPath    string
	ServiceInstanceID   string
}

// This method calls the MySQL CS Create Service Resource REST API.
// If successful, the API returns a HTTP 202 with a response object container the jobID and a message.
func (c *ResourceClient) createResource(requestBody interface{}, responseBody interface{}) error {

	var objectPath = c.getContainerPath(c.ContainerPath)
	c.client.DebugLogString(fmt.Sprintf("[Debug] : Trying to create ServiceInstance at %s", objectPath))
	_, err := c.executeRequest("POST", objectPath, requestBody)

	if err != nil {
		return err
	}

	return nil
}

func (c *ResourceClient) getResource(instanceName string, responseBody interface{}) error {

	var objectPath = c.getObjectPath(c.ResourceRootPath, instanceName)

	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)

}

// ServiceInstance needs a PUT and a body to be destroyed
func (c *ResourceClient) deleteResource(name string, requestBody interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, name)
	} else {
		objectPath = c.ResourceRootPath
	}
	_, err := c.executeRequest("PUT", objectPath, requestBody)
	if err != nil {
		return err
	}

	// No errors and no response body to write
	return nil
}

func (c *ResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
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
