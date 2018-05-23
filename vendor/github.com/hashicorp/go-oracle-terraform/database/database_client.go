package database

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const authHeader = "Authorization"
const tenantHeader = "X-ID-TENANT-NAME"

// Client - Client represents an authenticated database client, with compute credentials and an api client.
type Client struct {
	client     *client.Client
	authHeader *string
}

// NewDatabaseClient returns a database client
func NewDatabaseClient(c *opc.Config) (*Client, error) {
	databaseClient := &Client{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	databaseClient.client = client

	databaseClient.authHeader = databaseClient.getAuthenticationHeader()

	return databaseClient, nil
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

	debugReqString := fmt.Sprintf("HTTP %s Req (%s)", method, path)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		// Debug the body for database services
		debugReqString = fmt.Sprintf("%s:\nBody: %+v", debugReqString, string(reqBody))
	}
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)

	// Set the authentication headers
	req.Header.Add(authHeader, *c.authHeader)
	req.Header.Add(tenantHeader, *c.client.IdentityDomain)
	resp, err := c.client.ExecuteRequest(req)

	return resp, err
}

func (c *Client) getContainerPath(root string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain)
}

func (c *Client) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain, name)
}
