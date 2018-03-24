---
layout: "oraclepaas"
page_title: "Provider: Oracle Cloud Platform"
sidebar_current: "docs-oraclepaas-index"
description: |-
  The Oracle Cloud Platform (Oracle PaaS) provider is used to interact with resources supported by the Oracle Cloud Platform services. The provider needs to be configured with credentials for the Oracle Cloud Account.
---

# Oracle Cloud Platform Provider

The Oracle Cloud Platform (Oracle PaaS) provider is used to interact with resources supported by the [Oracle Cloud Platform](http://cloud.oracle.com/paas) services. The provider needs to be configured with credentials for the Oracle Cloud Account.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Oracle Cloud Platform provider
provider "oraclepaas" {
  user              = "..."
  password          = "..."
  identity_domain   = "..."
  database_endpoint = "..."
  java_endpoint     = "..."
}

# Create a Database Service Instance
resource "oraclepaas_database_service_instance" "default" {
  name              = "default-service-instance"
  description       = "default-service-instance"
  edition           = "EE"
  shape             = "oc1m"
  subscription_type = "HOURLY"
  version           = "12.2.0.1"
  ssh_public_key    = "ssh key"

  database_configuration {
    admin_password     = "Pa55_Word"
    sid                = "ORCL"
    backup_destination = "NONE"
    usable_storage     = 15
  }
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Optional) The username to use, generally your email address. It can also
  be sourced from the `OPC_USERNAME` environment variable.

* `password` - (Optional) The password associated with the username to use. It can also be sourced from
  the `OPC_PASSWORD` environment variable.

* `identity_domain` - (Optional) The Identity Domain or Service Instance ID of the environment to use. It can also be sourced from the `OPC_IDENTITY_DOMAIN` environment variable.  

* `database_endpoint` - (Optional) The database API endpoint to use, associated with your Oracle Cloud Platform account.
This is known as the `REST Endpoint` within the Oracle portal. It can also be sourced from the
`ORACLEPAAS_DATABASE_ENDPOINT` environment variable.

* `java_endpoint` - (Optional) The java API endpoint to use, associated with your Oracle Cloud Platform Account.
This is known as the `REST Endpoint` within the Oracle portal. It can also be sourced from the
`ORACLEPAAS_JAVA_ENDPOINT` environment variable.

* `max_retries` - (Optional) The maximum number of tries to make for a successful response when operating on
resources within Oracle Cloud Platform. It can also be sourced from the `OPC_MAX_RETRIES` environment variable.
Defaults to 1.

* `insecure` - (Optional) Skips TLS Verification for using self-signed certificates. Should only be used if
absolutely needed. Can also via setting the `OPC_INSECURE` environment variable to `true`.

## Testing

Credentials must be provided via the `OPC_USERNAME`, `OPC_PASSWORD`,
`OPC_IDENTITY_DOMAIN` and `ORACLEPAAS_DATABASE_ENDPOINT` and `ORACLEPAAS_JAVA_ENDPOINT` environment variables in order to run
acceptance tests.
