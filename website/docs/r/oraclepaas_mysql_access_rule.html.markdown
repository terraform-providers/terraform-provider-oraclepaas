---
layout: "oraclepaas"
page_title: "Oracle: oraclepaas_mysql_access_rule"
sidebar_current: "docs-oraclepaas-resource-access-rule"
description: |-
  Creates and manages a MySQL Access Rule for an Oracle MySQL Cloud service instance.

---

# oraclepaas_mysql_access_rule

## Example Usage


```hcl
resource "oraclepaas_mysql_service_instance" "default" {
	name                 = "SimpleMySQLInstance"
	description          = "This is a simple mysql instance"
	vm_public_key       = "A SSH public key"
	backup_destination   = "NONE"

	mysql_configuration = {
		db_name          = "demo_db"
		db_storage       = 25
		mysql_port       = 3306
	}
}

resource "oraclepass_mysql_access_rule" "myrule" {
	service_instance_id = "SimpleMySQLInstance"
	name                = "My Access Rule"
	description         = "My Simple Access Rule"
	protocol            = "tcp"
	ports               = "8000"
	source              = "0.0.0.0/24"
	destination         = "mysql_MASTER"
	enabled             = true
}
```

## Argument Reference

* `service_instance_id` - (Required) The name of MySQL instance to attach the access rule to. 

* `name` - (Required) Name of the rule.

* `description` - (Optional) Description of the rule.

* `protocol` - (Optional) Communication protocol for the rule. For example, tcp.

* `ports` - (Required) Ports for the rule. This can be a single port or a port range.

* `source` - (Required) The hosts from which traffic is allowed. For example, PUBLIC-INTERNET for any host on the Internet, a single IP address or a comma-separated list of subnets (in CIDR format) or IPv4 addresses.

* `destination` - (Required) The service component to allow traffic to. For example, mysql_MASTER.

* `enabled` - (Optional) Determines whether the access rule is enabled. Valid values are `true` and `false`. The Default is `true`.