package oraclepaas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPC_USERNAME", nil),
				Description: "The user name for OPAAS API operations.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPC_PASSWORD", nil),
				Description: "The user password for OPAAS API operations.",
			},

			"identity_domain": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPC_IDENTITY_DOMAIN", nil),
				Description: "The OPAAS identity domain for API operations",
			},

			"database_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORACLEPAAS_DATABASE_ENDPOINT", nil),
				Description: "The HTTP endpoint for Oracle Database operations.",
			},

			"java_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORACLEPAAS_JAVA_ENDPOINT", nil),
				Description: "The HTTP endpoint for Oracle Java operations.",
			},

			"application_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORACLEPAAS_APPLICATION_ENDPOINT", nil),
				Description: "The HTTP endpoint for the Oracle Application operations",
			},

			"mysql_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORACLEPAAS_MYSQL_ENDPOINT", nil),
				Description: "The HTTP endpoint for Oracle MySQL operations.",
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPC_MAX_RETRIES", 1),
				Description: "Maximum number retries to wait for a successful response when operating on resources within OPAAS (defaults to 1)",
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPC_INSECURE", false),
				Description: "Skip TLS Verification for self-signed certificates. Should only be used if absolutely required.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"oraclepaas_database_service_instance": dataSourceOraclePAASDatabaseServiceInstance(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"oraclepaas_java_access_rule":          resourceOraclePAASJavaAccessRule(),
			"oraclepaas_database_service_instance": resourceOraclePAASDatabaseServiceInstance(),
			"oraclepaas_java_service_instance":     resourceOraclePAASJavaServiceInstance(),
			"oraclepaas_database_access_rule":      resourceOraclePAASDatabaseAccessRule(),
			"oraclepaas_application_container":     resourceOraclePAASApplicationContainer(),
			"oraclepaas_mysql_service_instance":    resourceOraclePAASMySQLServiceInstance(),
			"oraclepaas_mysql_access_rule":         resourceOraclePAASMySQLAccessRule(),
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		User:                d.Get("user").(string),
		Password:            d.Get("password").(string),
		IdentityDomain:      d.Get("identity_domain").(string),
		ApplicationEndpoint: d.Get("application_endpoint").(string),
		DatabaseEndpoint:    d.Get("database_endpoint").(string),
		JavaEndpoint:        d.Get("java_endpoint").(string),
		MySQLEndpoint:       d.Get("mysql_endpoint").(string),
		MaxRetries:          d.Get("max_retries").(int),
		Insecure:            d.Get("insecure").(bool),
		UserAgent:           fmt.Sprintf("HashiCorp-Terraform-v%s", terraformVersion),
	}

	return config.Client()
}
