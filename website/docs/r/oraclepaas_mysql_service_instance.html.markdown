---
layout: "oraclepaas"
page_title: "Oracle: mysql_service_instance"
sidebar_current: "docs-oraclepaas-resource-service-instance"
description: |-
  Creates and manages an Oracle MySQL Cloud Service instance on the Oracle Cloud Platform.

---

# oraclepaas_mysql_service_instance
The `oraclepaas_mysql_service_instance` resource creates and manages an Oracle MySQL Cloud Service instance on the Oracle Cloud Platform.

## Example Usage

```hcl
resource "oraclepaas_mysql_service_instance" "default" {
  name                      = "SimpleMySQLInstance"
  description               = "This is a simple mysql instance"
  vm_public_key             = "A SSH public key"
  backup_destination        = "NONE"
  notification_email        = "myemail@mydomain.com"
  shape                     = "oc3"
  ssh_public_key            = "ssh-public-key"

  backups {
    cloud_storage_container = "https://uscom-east-1.storage.oraclecloud.com/v1/MyStorageAccount/MyContainer"
    cloud_storage_username  = "MyCloudStorageAccount"
    cloud_storage_password  = "MyCloudStoragePassword"
    create_if_missing       = "true"
  }

  mysql_configuration {
    db_name                 = "demo_db"
    db_storage              = 25
    mysql_port              = 3306
    mysql_username          = "root"
    mysql_password          = "MySqlPassword_1"

    enterprise_monitor_configuration {
      em_agent_username     = "MyEmAgentUser"
      em_agent_password     = "EmAgentPassw0rd"
      em_username           = "EmAdminUser"
      em_password           = "EmAdminPassw0rd"
      em_port               = 18443
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required). The name of MySQL Cloud Service instance.

* `description` - (Optional). A description of the MySQL Instance

* `ssh_public_key` - (Required). The public key for the secure shell (SSH). This key wil be used for authentication when the user logs on to the instance over SSH.

* `backup_destination` - (Required) The destination where the database backups will be stored.

* `shape` - (Required) The desired compute shape.  A shape defines the number of Oracle Compute Units (OCPUs) and amount of memory (RAM). See [About Shapes](http://www.oracle.com/pls/topic/lookup?ctx=cloud&id=OCSUG210) in _Using Oracle Compute Cloud Service_ for more information about shapes.

* `metering_frequency` - (Optional). The billing frequency of the service instance. Allowed values are `MONTHLY` and `HOURLY`

* `region` - (Optional). Specifies the region where the instance will be provisioned.

* `availability_domain` - (Optional) Name of the availability domain within the region where the Oracle Database Cloud Service instance is to be provisioned. This is applicable only if you wish to provision to an OCI instance.

* `notification_email` - (Optional) The email address to send notifications around successful or unsuccessful completions of the instance-creation operation.

* `ip_network` - (Optional) This attribute is only applicable to accounts where regions are supported. The three-part name of an IP network to which the service instance is added. For example: /Compute-identity_domain/user/object

* `subnet` -(Optional) This attribute is relevant to only Oracle Cloud Infrastructure. Specify the Oracle Cloud Identifier (OCID) of a subnet from a virtual cloud network (VCN) that you had created previously in Oracle Cloud Infrastructure. For the instructions to create a VCN and subnet, see [Prerequisites for Oracle Platform Services on Oracle Cloud Infrastructure](http://www.oracle.com/pls/topic/lookup?ctx=en/cloud/paas/java-cloud&id=oci_general_paasprereqs) in the Oracle Cloud Infrastructure documentation.

* `vm_user` - (Optional) The user name of account to be created in the VM.

* `backups` - (Optional) Provides Cloud Storage information for how to implement service instance backups. Backups is documented below

* `mysql_configuration` - (Required) Specified the detail of how to configure the MySQL database. mysql_configuration is documented below.

`backups` support the following :

* `cloud_storage_container` - (Required). Name of the Oracle Storage Cloud container used for store the backups.

* `cloud_storage_username` - (Required) Username for the Oracle Storage Cloud administrator.

* `cloud_storage_password` - (Required) Password for the Oracle Storage Cloud administrator.

* `create_if_missing` - (Optional) Specifies whether to create the container if it does not exist. Default value is `false`


`mysql_configuration` supports the following :

* `db_name` - (Optional). The name of the database instance. Default value is `mydatabase`

* `db_storage` - (Optional). The storage volume sice for MySQL data. The value must be between 25 to 1024. Defaults to 25 (GB)

* `mysql_charset` - (Optional) MySQL server character set. See [Supported Character Sets and Collation](http://dev.mysql.com/doc/en/charset-charsets.html). Default value is `utf8mb4`

* `mysql_collation` -(Optional) MySQL server collation. See [Supported Character Sets and Collation](http://dev.mysql.com/doc/en/charset-charsets.html) for the permissible collations of each character set.

* `mysql_port` - (Optional) The port number for the MySQL Server. The value must be between 3200-3399. Default value is `3306`

* `mysql_username` - (Optional) The Administration user for connecting to the service via th MySQL protocol. Default value is `root`.

* `mysql_password` - (Optional) The password for the MySQL Administration user.

* `source_service_name` - (Optional) When present, indicates that the service instance should be created as a "snapshot clone" of another service instance. Provide the name of the existing service instance whose snapshot is to be used. `db_name`, `mysql_charset`, `mysql_collation`, `enterpriseMonitor`, and associated MySQL server component parameters do not apply when cloning a service from a snapshot. For those parameters, the clone operation uses the values defined in the snapshot of the source service instance.

* `snapshot_name` - (Optional) The name of the snapshot of the service instance specified by `source_service_name` that is to be used to create a "snapshot clone". This parameter is valid only if `source_service_name` is specified.

* `enterprise_monitor_configuration` - (Optional) Provides the Enterprise Monitor configuration for the MySQL Instance. If this is omitted, there will be no EM created for the MySQL Instance. `enterprise_monitor_configuration` is documented below.

`enterprise_monitor_configuration` supports the following :

* `em_agent_username` - (Optional). Name for the Enterprise Monitor agent user.

* `em_agent_password` - (Optional). Password for MySQL Enterprise Monitor agent.

* `em_username` - (Optional) Name for the Enterprise Monitor Manager user.

* `em_password` - (Optional) Password for MySQL Enterprise Monitor manager.

* `em_port` - (Optional) The port number for the MySQL Enterprise Monitor instance. The default is 18443.
