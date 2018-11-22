---
layout: "oraclepaas"
page_title: "Oracle: oraclepaas_database_compute_nodes"
sidebar_current: "docs-oraclepaas-datasource-database-compute_nodes"
description: |-
  Gets information about the configuration of the Compute Nodes for a Oracle Database Cloud Service or Exadata Classic Cloud Service instance on the Oracle Cloud Platform.
---

# oraclepaas\_database\_compute\_nodes

Use this data source to access the details of the individual Compute Nodes of a Database Service Instance

## Example Usage

```hcl
data "oraclepaas_database_compute_nodes" "foo" {
  name = "database-service-instance-1"
}

output "compute_nodes" {
  value = "${data.oraclepaas_database_compute_nodes.foo.compute_nodes}"
}
```

## Argument Reference

* `name` - (Required) The name of the Database Service Instance

## Attributes Reference

* `compute_nodes` - List of compute nodes for the Database Cloud Service instance, see [Compute Node attributes](#compute-node-attributes)


### Compute Node attributes


* `connect_descriptor` - The connection descriptor for Oracle Net Services (SQL*Net).

* `connect_descriptor_with_public_ip` - The connection descriptor for Oracle Net Services (SQL*Net) with IP addresses instead of host names.

* `hostname` - The host name of the compute node.

* `initial_primary` - Indicates whether the compute node hosted the primary database of an Oracle Data Guard configuration when the service instance was created.

* `listener_port` - The listener port for Oracle Net Services (SQL*Net) connections.

* `pdb_name` - The name of the default PDB (pluggable database) created when the service instance was created.

* `reserved_ip` - The IP address of the compute node.

* `shape` - The Oracle Compute Cloud shape of the compute node.

* `sid` - The SID of the database on the compute node.

* `status` - The status of the compute node.

* `storage_allocated` - The size in GB of the storage allocated to the compute node. For compute nodes of a service instance hosting an Oracle RAC database, this number does not include the storage shared by the nodes
