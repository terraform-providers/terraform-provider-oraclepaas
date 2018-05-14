package application

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

// Client represents an authenticated application client, with credentials and an api client.
type Client struct {
	client *client.Client
}

// NewClient returns a new client for the application resources managed by Oracle
func NewClient(c *opc.Config) (*Client, error) {
	appClient := &Client{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	appClient.client = client

	return appClient, nil
}

func (c *Client) executeCreateUpdateRequest(method, path string, files map[string]string, parameters map[string]interface{}) (*http.Response, error) {
	req, err := c.client.BuildMultipartFormRequest(method, path, files, parameters)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("HTTP %s Path (%s)", method, path)
	// req.Header.Set("Content-Type", "multipart/form-data")
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)
	c.client.DebugLogString(fmt.Sprintf("Req (%+v)", req))

	// Set the authentication headers
	req.SetBasicAuth(*c.client.UserName, *c.client.Password)
	req.Header.Add("X-ID-TENANT-NAME", *c.client.IdentityDomain)

	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) executeRequest(method, path string, body interface{}) (*http.Response, error) {
	reqBody, err := c.client.MarshallRequestBody(body)
	if err != nil {
		return nil, err
	}

	req, err := c.client.BuildRequestBody(method, path, reqBody)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("HTTP %s Path (%s)", method, path)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)
	c.client.DebugLogString(fmt.Sprintf("Req (%+v)", req))

	// Set the authentiation headers
	req.SetBasicAuth(*c.client.UserName, *c.client.Password)
	req.Header.Add("X-ID-TENANT-NAME", *c.client.IdentityDomain)

	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) getContainerPath(root string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain)
}

func (c *Client) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain, name)
}
