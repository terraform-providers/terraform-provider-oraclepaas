package mysql

import (
	"fmt"
	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"net/http"
)

const AUTH_HEADER = "Authorization"
const TENANT_HEADER = "X-ID-TENANT-NAME"
const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_ORA_JSON = "application/vnd.com.oracle.oracloud.provisioning.Service+json"

/** This is the main client that deals with interacting with the OPC MySQL Services. It works with mySQL and mySQL-AccessRules
 */

type MySQLClient struct {
	client            *client.Client
	ServiceInstanceID string
	authHeader        *string
}

func NewMySQLClient(c *opc.Config) (*MySQLClient, error) {
	mysqlClient := &MySQLClient{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	mysqlClient.client = client

	return mysqlClient, nil
}

func (c *MySQLClient) executeRequest(method, path string, body interface{}) (*http.Response, error) {

	resp, err := c.executeRequestWithContentType(method, path, body, CONTENT_TYPE_ORA_JSON)
	return resp, err
}

func (c *MySQLClient) executeRequestWithContentType(method, path string, body interface{}, contentType string) (*http.Response, error) {

	// TODO: Possible bug in content type, especially when create access lists.
	reqBody, err := c.client.MarshallRequestBody(body)
	if err != nil {
		return nil, err
	}

	c.client.DebugLogString(fmt.Sprintf("[DEBUG] : Executing Request %s to %s with contentType : %s", method, path, contentType))

	req, err := c.client.BuildRequestBody(method, path, reqBody)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("HTTP %s Path (%s)", method, path)

	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(fmt.Sprintf("[DEBUG] : RequestString (%+v)", debugReqString))

	// Set the authentication headers
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", CONTENT_TYPE_JSON)
	req.SetBasicAuth(*c.client.UserName, *c.client.Password)
	req.Header.Add(TENANT_HEADER, *c.client.IdentityDomain)

	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *MySQLClient) getContainerPath(root string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain)
}

func (c *MySQLClient) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain, name)

}
