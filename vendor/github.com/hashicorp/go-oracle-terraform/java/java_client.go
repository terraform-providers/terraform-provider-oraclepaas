package java

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const authHeader = "Authorization"
const tenantHeader = "X-ID-TENANT-NAME"

// Client represents an authenticated java client, with compute credentials and an api client.
type Client struct {
	client     *client.Client
	authHeader *string
}

// NewJavaClient returns a new java client
func NewJavaClient(c *opc.Config) (*Client, error) {
	javaClient := &Client{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	javaClient.client = client

	javaClient.authHeader = javaClient.getAuthenticationHeader()

	return javaClient, nil
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
		// Output the request body json
		debugReqString = fmt.Sprintf("%s:\nBody: %+v", debugReqString, string(reqBody))
	}
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)

	// Set the authentiation headers
	req.Header.Add(authHeader, *c.authHeader)
	req.Header.Add(tenantHeader, *c.client.IdentityDomain)

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
