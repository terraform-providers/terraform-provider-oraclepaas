package opaas

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
)

type Config struct {
	User             string
	Password         string
	IdentityDomain   string
	MaxRetries       int
	Insecure         bool
	DatabaseEndpoint string
	JavaEndpoint     string
}

type OPAASClient struct {
	databaseClient *database.DatabaseClient
	javaClient     *java.JavaClient
}

func (c *Config) Client() (*OPAASClient, error) {

	userAgentString := fmt.Sprintf("HashiCorp-Terraform-v%s", terraform.VersionString())

	config := opc.Config{
		IdentityDomain: &c.IdentityDomain,
		Username:       &c.User,
		Password:       &c.Password,
		MaxRetries:     &c.MaxRetries,
		UserAgent:      &userAgentString,
	}

	if logging.IsDebugOrHigher() {
		config.LogLevel = opc.LogDebug
		config.Logger = opaasLogger{}
	}

	// Setup HTTP Client based on insecure
	httpClient := cleanhttp.DefaultClient()
	if c.Insecure {
		transport := cleanhttp.DefaultTransport()
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		httpClient.Transport = transport
	}

	config.HTTPClient = httpClient

	opaasClient := &OPAASClient{}

	if c.DatabaseEndpoint != "" {
		databaseEndpoint, err := url.ParseRequestURI(c.DatabaseEndpoint)
		if err != nil {
			return nil, fmt.Errorf("Invalid database endpoint URI: %+v", err)
		}
		config.APIEndpoint = databaseEndpoint
		databaseClient, err := database.NewDatabaseClient(&config)
		if err != nil {
			return nil, err
		}
		opaasClient.databaseClient = databaseClient
	}

	if c.JavaEndpoint != "" {
		javaEndpoint, err := url.ParseRequestURI(c.JavaEndpoint)
		if err != nil {
			return nil, fmt.Errorf("Invalid java endpoint URI: %+v", err)
		}
		config.APIEndpoint = javaEndpoint
		javaClient, err := java.NewJavaClient(&config)
		if err != nil {
			return nil, err
		}
		opaasClient.javaClient = javaClient
	}

	return opaasClient, nil
}

type opaasLogger struct{}

func (l opaasLogger) Log(args ...interface{}) {
	tokens := make([]string, 0, len(args))
	for _, arg := range args {
		if token, ok := arg.(string); ok {
			tokens = append(tokens, token)
		}
	}
	log.SetFlags(0)
	log.Print(fmt.Sprintf("go-oracle-terraform: %s", strings.Join(tokens, " ")))
}
