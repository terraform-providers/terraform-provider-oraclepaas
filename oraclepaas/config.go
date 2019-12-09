package oraclepaas

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-oracle-terraform/application"
	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/go-oracle-terraform/mysql"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
)

type Config struct {
	User                string
	Password            string
	IdentityDomain      string
	MaxRetries          int
	Insecure            bool
	DatabaseEndpoint    string
	JavaEndpoint        string
	ApplicationEndpoint string
	MySQLEndpoint       string
	UserAgent           string
}

type OPAASClient struct {
	databaseClient    *database.Client
	javaClient        *java.Client
	applicationClient *application.Client
	mysqlClient       *mysql.MySQLClient
}

func (c *Config) Client() (*OPAASClient, error) {

	userAgentString := c.UserAgent

	config := opc.Config{
		IdentityDomain: &c.IdentityDomain,
		Username:       &c.User,
		Password:       &c.Password,
		MaxRetries:     &c.MaxRetries,
		UserAgent:      &userAgentString,
	}

	if logging.IsDebugOrHigher() {
		config.LogLevel = opc.LogDebug
		config.Logger = oraclepaasLogger{}
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

	oraclepaasClient := &OPAASClient{}

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
		oraclepaasClient.databaseClient = databaseClient
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
		oraclepaasClient.javaClient = javaClient
	}

	if c.ApplicationEndpoint != "" {
		applicationEndpoint, err := url.ParseRequestURI(c.ApplicationEndpoint)
		if err != nil {
			return nil, fmt.Errorf("Invalid application endpoint URI: %+v", err)
		}
		config.APIEndpoint = applicationEndpoint
		applicationClient, err := application.NewClient(&config)
		if err != nil {
			return nil, err
		}
		oraclepaasClient.applicationClient = applicationClient
	}

	if c.MySQLEndpoint != "" {
		mysqlEndpoint, err := url.ParseRequestURI(c.MySQLEndpoint)
		if err != nil {
			return nil, fmt.Errorf("Invalid jmysqlava endpoint URI: %+v", err)
		}
		config.APIEndpoint = mysqlEndpoint
		mysqlClient, err := mysql.NewMySQLClient(&config)
		if err != nil {
			return nil, err
		}
		oraclepaasClient.mysqlClient = mysqlClient
	}

	return oraclepaasClient, nil
}

type oraclepaasLogger struct{}

func (l oraclepaasLogger) Log(args ...interface{}) {
	tokens := make([]string, 0, len(args))
	for _, arg := range args {
		if token, ok := arg.(string); ok {
			tokens = append(tokens, token)
		}
	}
	log.SetFlags(0)
	log.Print(fmt.Sprintf("go-oracle-terraform: %s", strings.Join(tokens, " ")))
}
