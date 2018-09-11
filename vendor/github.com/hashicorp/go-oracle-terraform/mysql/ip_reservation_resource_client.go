package mysql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// IPReservationResourceClient is a client for the IP Reservation functions of the Java Cloud API.
type IPReservationResourceClient struct {
	*MySQLClient
	ContainerPath    string
	ResourceRootPath string
}

func (c *IPReservationResourceClient) createResource(requestBody interface{}) (*CreateIPReservationInfo, error) {
	resp, err := c.executeRequestWithContentType("POST", c.getContainerPath(c.ContainerPath), requestBody, "application/json")
	if err != nil {
		return nil, err
	}

	var info CreateIPReservationInfo
	if err := c.unmarshalResponseBody(resp, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *IPReservationResourceClient) getResource(name string) (*IPReservationInfo, error) {
	objectPath := c.getContainerPath(c.ContainerPath)

	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return nil, err
	}

	var ipReservations IPReservations
	if err := c.unmarshalResponseBody(resp, &ipReservations); err != nil {
		return nil, err
	}

	// API returns all IP Reservations, iterate to find the one we want
	for _, ipRes := range ipReservations.IPReservations {
		if ipRes.Name == name {
			var ipReservation *IPReservationInfo
			ipReservation = &ipRes
			return ipReservation, nil
		}
	}
	return nil, fmt.Errorf("IP Reservation not found")
}

func (c *IPReservationResourceClient) deleteResource(name string) (*DeleteIPReservationInfo, error) {
	objectPath := c.getObjectPath(c.ResourceRootPath, name)

	resp, err := c.executeRequest("DELETE", objectPath, nil)
	if err != nil {
		return nil, err
	}

	var info DeleteIPReservationInfo
	if err := c.unmarshalResponseBody(resp, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *IPReservationResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
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

func (c *IPReservationResourceClient) getContainerPath(root string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain)
}

func (c *IPReservationResourceClient) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain, name)
}
