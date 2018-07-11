---
layout: "oraclepaas"
page_title: "Oracle: oraclepaas_database_access_rule"
sidebar_current: "docs-oraclepaas-resource-access-rule"
description: |-
  Creates and manages a Database Access Rule for an Oracle Database Cloud service instance.

---

# oraclepaas_database_access_rule

The `oraclepaas_database_access_rule` resource creates and manages a Database Access Rule for an Oracle Database Cloud service instance.

## Example Usage

```hcl
resource "oraclepaas_database_service_instance" "default" {
  name = "database-service-instance-1"
  ...
}

resource "oraclepaas_database_access_rule" "default" {
	name                = "example-access-rule"
	service_instance_id = "${oraclepaas_database_service_instance.default.name}"
	description         = "enable port 8000"
	ports               = "8000"
	source              = "PUBLIC-INTERNET"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Access Rule

* `service_instance_id` - (Required) The name of the database service instance to attach
 the access rule to

* `description` - (Required) The description of the Access Rule

* `ports` - (Required) The port or range of ports to allow traffic on

* `source` - (Required) The IP addresses and subnets from which traffic is allowed. Valid values are
`DB`, `PUBLIC-INTERNET`, or a single IP address or comma-separated list of subnets (in CIDR format) or IPv4 addresses.

* `enabled` - (Optional)  Determines whether the access rule is enabled. Default is `true`.
