---
layout: "oraclepaas"
page_title: "Oracle: oraclepaas_java_access_rule"
sidebar_current: "docs-oraclepaas-resource-access-rule"
description: |-
  Creates and manages an Access Rule for an Java Cloud service instance.

---

# oraclepaas_java_access_rule

The `oraclepaas_java_access_rule` resource creates and manages an Access Rule for an Java Cloud service instance.

## Example Usage

```hcl
resource "oraclepaas_java_service_instance" "default" {
  name        = "java-service-instance-1"
  ...
}

resource "oraclepaas_java_access_rule" "default" {
	name                = "example-access-rule"
	service_instance_id = "${oraclepaas_java_service_instance.default.name}"
	description         = "enable port 8000"
	ports               = "8000"
	source              = "PUBLIC-INTERNET"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Access Rule

* `service_instance_id` - (Required) The name of the java service instance to attach
 the access rule to

* `description` - (Required) The description of the Access Rule

* `ports` - (Required) The port or range of ports to allow traffic on

* `destination` - (Required) Destination to which traffic is allowed. Valid values include `WLS_ADMIN`, `WLS_ADMIN_SERVER`, `OTD_ADMIN_HOST`, `OTD`

* `source` - (Required) The IP addresses and subnets from which traffic is allowed. Valid values include `WLS_ADMIN`, `WLS_ADMIN_SERVER`,
`WLS_MANAGED_SERVER`, `OTD_ADMIN_HOST`, `OTD`, or a single IP address or comma-separated list of subnets (in CIDR format) or IPv4 addresses.

* `enabled` - (Optional) Determines whether the access rule is enabled. Default is `true`.

* `protocol` - (Optional) Specifies the communication protocol. Valid values are `tcp` or `udp`.
Default is `tcp`.
