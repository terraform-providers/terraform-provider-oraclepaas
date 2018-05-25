package java

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

func (c *ResourceClient) createResource(requestBody interface{}, responseBody interface{}) error {
	_, err := c.executeRequest("POST", c.getContainerPath(c.ContainerPath), requestBody)
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

func (c *ResourceClient) updateResource(name, path, method string, requestBody interface{}) error {
	_, err := c.executeRequest(method, fmt.Sprintf("%s%s", c.getObjectPath(c.ResourceRootPath, name), path), requestBody)
	return err
}

// ServiceInstance needs a PUT and a body to be destroyed
func (c *ResourceClient) deleteInstanceResource(name string, requestBody interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, name)
	} else {
		objectPath = c.ResourceRootPath
	}
	_, err := c.executeRequest("PUT", objectPath, requestBody)
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
		return fmt.Errorf("Error decoding: %s\n%+v", err.Error(), resp)
	}

	// Use mapstructure to weakly decode into the resulting interface
	var msdcd *mapstructure.Decoder
	msdcd, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
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
