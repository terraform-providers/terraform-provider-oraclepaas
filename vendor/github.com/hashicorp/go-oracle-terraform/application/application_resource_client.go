package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// ResourceClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type ResourceClient struct {
	*Client
	ContainerPath    string
	ResourceRootPath string
}

func (c *ResourceClient) createResource(files map[string]string, additionalParams map[string]interface{}, responseBody interface{}) error {
	_, err := c.executeCreateUpdateRequest("POST", c.getContainerPath(c.ContainerPath), files, additionalParams)

	return err
}

func (c *ResourceClient) updateResource(files map[string]string, additionalParams map[string]interface{}, responseBody interface{}) error {
	_, err := c.executeCreateUpdateRequest("PUT", c.getContainerPath(c.ContainerPath), files, additionalParams)

	return err
}

func (c *ResourceClient) getResource(name string, responseBody interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, name)
	} else {
		objectPath = c.ResourceRootPath
	}
	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) deleteResource(name string) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, name)
	} else {
		objectPath = c.ResourceRootPath
	}
	_, err := c.executeRequest("DELETE", objectPath, nil)

	return err
}

func (c *ResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	c.client.DebugLogString(fmt.Sprintf("HTTP Resp (%d): %s", resp.StatusCode, buf.String()))
	// JSON decode response into interface
	var tmp interface{}
	dcd := json.NewDecoder(buf)
	if err = dcd.Decode(&tmp); err != nil {
		return err
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
