---
layout: "oraclepaas"
page_title: "Oracle: oraclepaas_database_service_instance"
sidebar_current: "docs-oraclepaas-datasource-database-service-instance"
description: |-
  Gets information about the configuration of an Oracle Database Cloud Service instance on the Oracle Cloud Platform.
---

# oraclepaas\_database\_service\_instance

Use this data source to access the configuration of a Database Service Instance

## Example Usage

```hcl
data "oraclepaas_database_service_instance" "foo" {
  name = "database-service-instance-1"
}

output "region" {
  value = "${data.oraclepaas_database_service_instance.foo.region}"
}
```

## Argument Reference

* `name` - (Required) The name of the Database Service Instance

## Attributes Reference

* `apex_url` - The URL to use to connect to Oracle Application Express on the service instance.
* `availability_domain` - Name of the availability domain within the region where the Oracle Database Cloud Service instance is provisioned.
* `backup_destination`- The backup configuration of the service instance.
* `character_set` - The database character set of the database.
* `cloud_storage_container` - The Oracle Storage Cloud container for backups.
* `compute_site_name` - The Oracle Cloud location housing the service instance.
* `description` - The description of the service instance.
* `edition` - The software edition of the service instance.
* `enterprise_manager_url` - The URL to use to connect to Enterprise Manager on the service instance.
* `failover_database` - Indicates whether the service instance hosts an Oracle Data Guard configuration.
* `glassfish_url` - The URL to use to connect to the Oracle GlassFish Server Administration Console on the service instance.
* `hybrid_disaster_recovery_ip` - Data Guard Role of the on-premise instance in Oracle Hybrid Disaster Recovery configuration.
* `identity_domain` - The identity domain housing the service instance.
* `ip_network` - The three-part name of an IP network to which the service instance is added.
* `ip_reservations` - Groups one or more IP reservations in use on this service instance.
* `bring_your_own_license` - Indicates whether service instance was provisioned with the 'Bring Your Own License' option.
* `level` - The service level of the service instance.
* `listener_port` - The listener port for Oracle Net Services (SQL*Net) connections.
* `monitor_url` - The URL to use to connect to Oracle DBaaS Monitor on the service instance.
* `national_character_set` - The national character set of the database.
* `pluggable_database_name` - The name of the default PDB (pluggable database) created when the service instance was created.
* `region` - Location where the service instance is provisioned.
* `shape` - The Oracle Compute Cloud shape of the service instance.
* `high_performance_storage` - Indicates whether the service instance was provisioned with high performance storage.
* `uri` - The REST endpoint URI of the service instance.
* `version` - The Oracle Database version on the service instance.
