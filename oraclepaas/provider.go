package oraclepaas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
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
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		User:                d.Get("user").(string),
		Password:            d.Get("password").(string),
		IdentityDomain:      d.Get("identity_domain").(string),
		ApplicationEndpoint: d.Get("application_endpoint").(string),
		DatabaseEndpoint:    d.Get("database_endpoint").(string),
		JavaEndpoint:        d.Get("java_endpoint").(string),
		MaxRetries:          d.Get("max_retries").(int),
		Insecure:            d.Get("insecure").(bool),
	}

	return config.Client()
}
