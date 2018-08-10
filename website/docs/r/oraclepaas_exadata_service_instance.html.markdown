---
layout: "oraclepaas"
page_title: "Oracle: oraclepaas_exadata_service_instance"
sidebar_current: "docs-oraclepaas-resource-service-instance"
description: |-
  Creates and manages an Oracle Exadata Classic Cloud Service instance on the Oracle Cloud Platform.

---

# oraclepaas\_exadata\_service\_instance

The `oraclepaas_exadata_service_instance` resource creates and manages an Oracle Exadata Classic Cloud Service instance on the Oracle Cloud Platform.

## Example Usage

```hcl
resource "oraclepaas_exadata_service_instance" "default" {
  name        = "EXATEST01"
  description = "This is a description for an service instance"

  exadata_system_name = "exad-sys1"
  cluster_name        = "exad-sys1-000"
  version             = "12.2.0.1"

  database_configuration {
      sid                = "MYORCL1"
      pdb_name           = "MYPDB01"
      admin_password     = "Pa55_Word"
      backup_destination = "OSS"
  }

  backups {
      cloud_storage_container = "Storage-${var.domain}/database-service-instance-backup"
      auto_generate = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Instance.

* `version` - (Required) Oracle Database software version; one of: `12.2.0.1`, `12.1.0.2`, or `11.2.0.4`.

* `database_configuration` - (Required) Specifies the details on how to configure the database. Database configuration is documented below.

* `backups` - (Optional) Provides Cloud Storage information for how to implement service instance backups. Backups is documented below

* `bring_your_own_license` - (Optional) Specify if you want to use an existing perpetual license to Oracle Database to establish the right to use Oracle Database on the new instance.
Default value is `false`.

* `description` - (Optional) A description of the Service Instance.

* `desired_state` - (Optional) Specifies the desired state of the service instance. Allowed values are `start`, `stop`,
and `restart`.

* `edition` - (Optional) Database edition for the service instance. For Exadata only `EE_EP` is supported.

* `instantiate_from_backup` - (Optional) Specify if the service instance's database should, after the instance is created, be replaced by a database
stored in an existing cloud backup that was created using Oracle Database Backup Cloud Service. Instantiate from Backup is documented below.

* `level` - (Optional) Service level for the service instance. For Exadata only `PAASEXADATA` is supported.

* `notification_email` - (Optional)  The email address to send notifications around successful or unsuccessful completions of the instance-creation operation.

* `subscription_type` - (Optional) Billing unit. For Exadata only `MONTHLY` is supported.


Database Configuration supports the following:

* `admin_password` - (Required) Password for Oracle Database administrator users sys and system. The password must meet the following requirements: Starts with a letter. Is between 8 and 30 characters long. Contains letters, at least one number, and optionally, any number of these special characters: dollar sign `$`, pound sign `#`, and underscore `_`.

* `backup_destination` - (Optional) Backup Destination. Possible values are `BOTH`, `OSS`, `NONE`.This defaults to `NONE`.

* `character_set` - (Optional) Character Set for the Database Cloud Service Instance. Default value is `AL32UTF8`. Supported values are:

  - `AL32UTF8`, `AR8ADOS710`, `AR8ADOS720`, `AR8APTEC715`, `AR8ARABICMACS`, `AR8ASMO8X`, `AR8ISO8859P6`, `AR8MSWIN1256`, `AR8MUSSAD768`, `AR8NAFITHA711`, `AR8NAFITHA721`, `AR8SAKHR706`, `AR8SAKHR707`, `AZ8ISO8859P9E`, `BG8MSWIN`, `BG8PC437S`, `BLT8CP921`, `BLT8ISO8859P13`, `BLT8MSWIN1257`, `BLT8PC775`, `BN8BSCII`, `CDN8PC863`, `CEL8ISO8859P14`, `CL8ISO8859P5`, `CL8ISOIR111`, `CL8KOI8R`, `CL8KOI8U`, `CL8MACCYRILLICS`, `CL8MSWIN1251`, `EE8ISO8859P2`, `EE8MACCES`, `EE8MACCROATIANS`, `EE8MSWIN1250`, `EE8PC852`, `EL8DEC`, `EL8ISO8859P7`, `EL8MACGREEKS`, `EL8MSWIN1253`, `EL8PC437S`, `EL8PC851`, `EL8PC869`, `ET8MSWIN923`, `HU8ABMOD`, `HU8CWI2`, `IN8ISCII`, `IS8PC861`, `IW8ISO8859P8`, `IW8MACHEBREWS`, `IW8MSWIN1255`, `IW8PC1507`, `JA16EUC`, `JA16EUCTILDE`, `JA16SJIS`, `JA16SJISTILDE`, `JA16VMS`, `KO16KSC5601`, `KO16KSCCS`, `KO16MSWIN949`, `LA8ISO6937`, `LA8PASSPORT`, `LT8MSWIN921`, `LT8PC772`, `LT8PC774`, `LV8PC1117`, `LV8PC8LR`, `LV8RST104090`, `N8PC865`, `NE8ISO8859P10`, `NEE8ISO8859P4`, `RU8BESTA`, `RU8PC855`, `RU8PC866`, `SE8ISO8859P3`, `TH8MACTHAIS`, `TH8TISASCII`, `TR8DEC`, `TR8MACTURKISHS`, `TR8MSWIN1254`, `TR8PC857`, `US7ASCII`, `US8PC437`, `UTF8`, `VN8MSWIN1258`, `VN8VN3`, `WE8DEC`, `WE8DG`, `WE8ISO8859P1`, `WE8ISO8859P15`, `WE8ISO8859P9`, `WE8MACROMAN8S`, `WE8MSWIN1252`, `WE8NCR4970`, `WE8NEXTSTEP`, `WE8PC850`, `WE8PC858`, `WE8PC860`, `WE8ROMAN8`, `ZHS16CGB231280`, `ZHS16GBK`, `ZHT16BIG5`, `ZHT16CCDC`, `ZHT16DBT`, `ZHT16HKSCS`, `ZHT16MSWIN950`, `ZHT32EUC`, `ZHT32SOPS`, `ZHT32TRIS`.

* `failover_database` - (Optional) Specify if an Oracle Data Guard configuration comprising a primary database and a standby database is created.
Default value is `false`.

* `golden_gate` - (Optional) Specify if the database should be configured for use as the replication database of an Oracle GoldenGate Cloud Service instance.
You cannot set `goldenGate` to `true` if either `is_rac` or `failoverDatabase` is set to `true`. Default value is `false`.

* `is_rac` - (Optional) Specify if a cluster database using Oracle Real Application Clusters should be configured.
Default value is `true`.

* `national_character_set` - (Optional) National Character Set for the Database Cloud Service instance. Valid values are `AL16UTF16` and `UTF8`.

* `oracle_home_name` - (Optional) Name for the Oracle Home directory location that you want to use for the database deployment. If you specify the name for an existing Oracle Home directory location, then the database deployment shares the existing Oracle Database binaries at that location. Otherwise a new Oracle Home is created.

* `pdb_name` - (Optional) This attribute is valid when Database Cloud Service instance is configured with version 12c. Pluggable Database Name for the Database Cloud Service instance. Default value is `pdb1`.

* `sid` - (Optional) Database Name for the Database Cloud Service instance. Default value is `ORCL`.

* `source_service_name` - (Optional) Indicates that the service instance should be created as a "snapshot clone" of another service instance. Provide the name of the existing service instance whose snapshot is to be used.

* `snapshot_name` - (Optional) The name of the snapshot of the service instance specified by sourceServiceName that is to be used to create a "snapshot clone". This parameter is valid only if source_service_name is specified.

* `type` - (Optional) Component type to which the set of parameters applies. Defaults to `db`


Standby supports the following:

* `cluster_name` - (Required) Name of the cluster on which to create the standby database for a database deployment that uses Oracle Data Guard.

* `exadata_system_name` - (Required) Name of the Exadata system on which to create the standby database for a database deployment that uses Oracle Data Guard.

* `node_list` - (Optional) Specifies the list of compute nodes that host database instances for the standby database. If `node_list` is not specified the database is deployed across all compute nodes.


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
