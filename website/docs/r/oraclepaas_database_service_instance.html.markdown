---
subcategory: "PaaS"
layout: "oraclepaas"
page_title: "Oracle: oraclepaas_database_service_instance"
sidebar_current: "docs-oraclepaas-resource-service-instance"
description: |-
  Creates and manages an Oracle Database Cloud Service instance on the Oracle Cloud Platform.

---

# oraclepaas\_database\_service\_instance

The `oraclepaas_database_service_instance` resource creates and manages a an Oracle Database Cloud Service instance on the Oracle Cloud Platform.

## Example Usage

```hcl
resource "oraclepaas_database_service_instance" "default" {
  name        = "database-service-instance"
  description = "This is a description for an service instance"

  edition           = "EE"
  shape             = "oc1m"
  subscription_type = "HOURLY"
  version           = "12.2.0.1"
  vm_public_key     = "An ssh public key"

  database_configuration {
      admin_password     = "Pa55_Word"
      sid                = "ORCL"
      backup_destination = "BOTH"
      usable_storage     = 15
  }

  backups {
      cloud_storage_container = "Storage-${var.domain}/database-service-instance-backup"
      auto_generate           = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Instance.

* `edition` - (Required) Database edition for the service instance. Possible values are `SE`, `EE`, `EE_HP`, or `EE_EP`.

* `level` - (Required) Service level for the service instance. Possible values are `BASIC` or `PAAS`.

* `shape` - (Required) Desired compute shape. Possible values are `oc3`, `oc4`, `oc5`, `oc6`, `oc1m`, `oc2m`, `oc3m`, or `oc4m`.

* `subscription_type` - (Required) Billing unit. Possible values are `HOURLY` or `MONTHLY`.

* `version` - (Required) Oracle Database software version; one of: `18.0.0.0`, `12.2.0.1`, `12.1.0.2`, or `11.2.0.4`.

* `vm_public_key` - (Required) Public key for the secure shell (SSH). This key will be used for authentication when connecting to the Database Cloud Service instance using an SSH client.

* `database_configuration` - (Required) Specifies the details on how to configure the database. Database configuration is documented below.

* `default_access_rules` - (Optional) Specifies the details on which default access rules are enable or disabled. Default Access Rules
are configured below.

* `desired_state` - (Optional) Specifies the desired state of the service instance. Allowed values are `start`, `stop`,
and `restart`.

* `instantiate_from_backup` - (Optional) Specify if the service instance's database should, after the instance is created, be replaced by a database
stored in an existing cloud backup that was created using Oracle Database Backup Cloud Service. Instantiate from Backup is documented below.

* `ip_network` - (Optional) This attribute is only applicable to accounts where regions are supported. The three-part name of an IP network to which the service instance is added. For example: /Compute-identity_domain/user/object

* `ip_reservations` - (Optional) Groups one or more IP reservations in use on this service instance. This attribute is only applicable to accounts where regions are supported.

* `backups` - (Optional) Provides Cloud Storage information for how to implement service instance backups. Backups is documented below

* `bring_your_own_license` - (Optional) Specify if you want to use an existing perpetual license to Oracle Database to establish the right to use Oracle Database on the new instance.
Default value is `false`.

* `description` - (Optional) A description of the Service Instance.

* `high_performance_storage` - (Optional) Specifies whether the service instance will be provisioned with high performance storage.
Default value is `false`.

* `hybrid_disastery_recovery` - (Optional) Provides information about an Oracle Hybrid Disaster Recovery configuration. Hybrid Disaster Recovery is documented below.

* `notification_email` - (Optional)  The email address to send notifications around successful or unsuccessful completions of the instance-creation operation.

* `region` - (Optional) Specifies the location where the service instance is provisioned (only for accounts where regions are supported).

* `standby` - (Optional) Specifies the configuration details of the standby database. This is only applicable in Oracle Cloud Infrastructure Regions. `failover_database` and
`disaster_recovery` inside the `database_configuration` block must be set to `true`. Standby is documented below.

* `subnet` - (Optional) Name of the subnet within the region where the Oracle Database Cloud Service instance is to be provisioned.

Database Configuration supports the following:

* `admin_password` - (Required) Password for Oracle Database administrator users sys and system. The password must meet the following requirements: Starts with a letter. Is between 8 and 30 characters long. Contains letters, at least one number, and optionally, any number of these special characters: dollar sign `$`, pound sign `#`, and underscore `_`.

* `usable_storage` - (Optional) Storage size for data (in GB). Minimum value is `15`. Maximum value depends on the backup destination: if `BOTH` is specified, the maximum value is `1200`; if `OSS` or `NONE` is specified, the maximum value is `2048`.

* `availability_domain` - (Optional) Name of the availability domain within the region where the Oracle Database Cloud Service instance is to be provisioned.

* `backup_destination` - (Optional) Backup Destination. Possible values are `BOTH`, `OSS`, `NONE`.This defaults to `NONE`.

* `backup_storage_volume_size` - (Optional) The size (in GB) for the backup storage volume. 

* `character_set` - (Optional) Character Set for the Database Cloud Service Instance. Default value is `AL32UTF8`. Supported values are:

  - `AL32UTF8`, `AR8ADOS710`, `AR8ADOS720`, `AR8APTEC715`, `AR8ARABICMACS`, `AR8ASMO8X`, `AR8ISO8859P6`, `AR8MSWIN1256`, `AR8MUSSAD768`, `AR8NAFITHA711`, `AR8NAFITHA721`, `AR8SAKHR706`, `AR8SAKHR707`, `AZ8ISO8859P9E`, `BG8MSWIN`, `BG8PC437S`, `BLT8CP921`, `BLT8ISO8859P13`, `BLT8MSWIN1257`, `BLT8PC775`, `BN8BSCII`, `CDN8PC863`, `CEL8ISO8859P14`, `CL8ISO8859P5`, `CL8ISOIR111`, `CL8KOI8R`, `CL8KOI8U`, `CL8MACCYRILLICS`, `CL8MSWIN1251`, `EE8ISO8859P2`, `EE8MACCES`, `EE8MACCROATIANS`, `EE8MSWIN1250`, `EE8PC852`, `EL8DEC`, `EL8ISO8859P7`, `EL8MACGREEKS`, `EL8MSWIN1253`, `EL8PC437S`, `EL8PC851`, `EL8PC869`, `ET8MSWIN923`, `HU8ABMOD`, `HU8CWI2`, `IN8ISCII`, `IS8PC861`, `IW8ISO8859P8`, `IW8MACHEBREWS`, `IW8MSWIN1255`, `IW8PC1507`, `JA16EUC`, `JA16EUCTILDE`, `JA16SJIS`, `JA16SJISTILDE`, `JA16VMS`, `KO16KSC5601`, `KO16KSCCS`, `KO16MSWIN949`, `LA8ISO6937`, `LA8PASSPORT`, `LT8MSWIN921`, `LT8PC772`, `LT8PC774`, `LV8PC1117`, `LV8PC8LR`, `LV8RST104090`, `N8PC865`, `NE8ISO8859P10`, `NEE8ISO8859P4`, `RU8BESTA`, `RU8PC855`, `RU8PC866`, `SE8ISO8859P3`, `TH8MACTHAIS`, `TH8TISASCII`, `TR8DEC`, `TR8MACTURKISHS`, `TR8MSWIN1254`, `TR8PC857`, `US7ASCII`, `US8PC437`, `UTF8`, `VN8MSWIN1258`, `VN8VN3`, `WE8DEC`, `WE8DG`, `WE8ISO8859P1`, `WE8ISO8859P15`, `WE8ISO8859P9`, `WE8MACROMAN8S`, `WE8MSWIN1252`, `WE8NCR4970`, `WE8NEXTSTEP`, `WE8PC850`, `WE8PC858`, `WE8PC860`, `WE8ROMAN8`, `ZHS16CGB231280`, `ZHS16GBK`, `ZHT16BIG5`, `ZHT16CCDC`, `ZHT16DBT`, `ZHT16HKSCS`, `ZHT16MSWIN950`, `ZHT32EUC`, `ZHT32SOPS`, `ZHT32TRIS`.

* `data_storage_volume_size` - (Optional) The size (in GB) for the data storage volume.

* `disaster_recovery` - (Optional) Specify if an Oracle Data Guard configuration is created using the Disaster Recovery option or the High Availability option.
Default value is `false`.

* `failover_database` - (Optional) Specify if an Oracle Data Guard configuration comprising a primary database and a standby database is created.
Default value is `false`.

* `golden_gate` - (Optional) Specify if the database should be configured for use as the replication database of an Oracle GoldenGate Cloud Service instance.
You cannot set `goldenGate` to `true` if either `is_rac` or `failoverDatabase` is set to `true`. Default value is `false`.

* `is_rac` - (Optional) Specify if a cluster database using Oracle Real Application Clusters should be configured.
Default value is `false`.

* `national_character_set` - (Optional) National Character Set for the Database Cloud Service instance. Valid values are `AL16UTF16` and `UTF8`.

* `pdb_name` - (Optional) This attribute is valid when Database Cloud Service instance is configured with version 12c. Pluggable Database Name for the Database Cloud Service instance. Default value is `pdb1`.

* `sid` - (Optional) Database Name for the Database Cloud Service instance. Default value is `ORCL`.

* `source_service_name` - (Optional) Indicates that the service instance should be created as a "snapshot clone" of another service instance. Provide the name of the existing service instance whose snapshot is to be used.

* `snapshot_name` - (Optional) The name of the snapshot of the service instance specified by sourceServiceName that is to be used to create a "snapshot clone". This parameter is valid only if source_service_name is specified.

* `timezone` - (Optional) Time Zone for the Database Cloud Service instance. Default value is `UTC`.

* `type` - (Optional) Component type to which the set of parameters applies. Defaults to `db`

* `db_demo` - (Optional) Indicates whether to include the Demos PDB.

Default Access Rules supports the following:

* `enable_ssh` - (Optional) Indicates whether to enable the ssh access rule.

* `enable_http` - (Optional) Indicates whether to enable the http access rule. This is only configurable with a single instance.

* `enable_https` - (Optional) Indiciates whether to enable the http with ssl access rule. This is only configurable with a single instance.

* `enable_db_console` - (Optional) Indicates whether to enable the db console access rule. This is only configurable with a single instance.

* `enable_db_express` - (Optional) Indicates whether to enable the db express access rule. This is only configurable with a single instance.

* `enable_db_listener` - (Optional) Indicates whether to enable the db listener access rule. This is only configurable with a single instance

* `enable_em_console` - (Optional) Indicates whether to enable the em console access rule. This is only configurable with a RAC instance.

* `enable_rac_db_listener` - (Optional) Indicates whether to enable the rac db listene access rule. This is only configurable with a RAC instance

* `enable_scan_listener` - (Optional) Indicates whether to enable the scan listener access rule. This is only configurable with a RAC instance

* `enable_rac_ons` - (Optional) Indicates whether to enable the rac ons access rule. This is only configurable with a RAC instance.

Standby supports the following:

* `availability_domain` - (Required) Name of the availability domain within the region where the standby database of the Oracle Database Cloud Service instance is to be provisioned.

* `subnet` - (Required) Name of the subnet within the region where the standby database of the Oracle Database Cloud Service instance is to be provisioned.

Instantiate from Backup supports the following:

* `cloud_storage_container` - (Required) Name of the Oracle Storage Cloud Service container where the existing cloud backup is stored.

* `cloud_storage_username` - (Required) Username of the Oracle Cloud user.

* `cloud_storage_password` - (Required) Password of the Oracle Cloud user specified in `ibkup_cloud_storage_user`.

* `database_id` - (Required) Database id of the database from which the existing cloud backup was created.

* `decryption_key` - (Optional) Password used to create the existing, password-encrypted cloud backup. This password is used to decrypt the backup. Specify either `ibkup_decryption_key` or `ibkup_wallet_file_content` for decrypting the backup.

* `on_premise` - (Optional) Specify if the existing cloud backup being used to replace the database is from an on-premises database or another Database Cloud Service instance.
The default value is false.

* `service_id` - (Optional) Oracle Database Cloud Service instance name from which the database of new Oracle Database Cloud Service instance should be created. This value is required if
`on_premise` is set to true.

* `wallet_file_content` - (Optional) String containing the xsd:base64Binary representation of the cloud backup's wallet file. This wallet is used to decrypt the backup. Specify either `ibkup_decryption_key` or `ibkup_wallet_file_content` for decrypting the backup.

Backups support the following:

* `cloud_storage_container` - (Required) Name of the Oracle Storage Cloud Service container used to provide storage for your service instance backups.
Use the following format to specify the container name: `<storageservicename>-<storageidentitydomain>/<containername>`

* `cloud_storage_username` - (Required) Username for the Oracle Storage Cloud Service administrator.

* `cloud_storage_password` - (Required) Password for the Oracle Storage Cloud Service administrator.

* `create_if_missing` - (Optional) Specify if the given cloud_storage_container is to be created if it does not already exist. Default value is `false`.

Hybrid Disaster Recovery supports the following:

* `cloud_storage_container` - (Required) Name of the Oracle Storage Cloud Service container where the backup from on-premise instance is stored.
Use the following format to specify the container name: `<storageservicename>-<storageidentitydomain>/<containername>`

* `cloud_storage_username` - (Required) Username for the Oracle Storage Cloud Service administrator.

* `cloud_storage_password` - (Required) Password for the Oracle Storage Cloud Service administrator.

In addition to the above, the following values are exported:

* `compute_site_name` - The Oracle Cloud location housing the service instance.

* `dbaas_monitor_url`- The URL to use to connect to Oracle DBaaS Monitor on the service instance.

* `em_url` - The URL to use to connect to Enterprise Manager on the service instance.

* `glassfish_url` - The URL to use to connect to the Oracle GlassFish Server Administration Console on the service instance.

* `identity_domain` - The identity domain housing the service instance.

* `status` - The status of the service instance.

* `uri` - The Uniform Resource Identifier for the Service Instance
