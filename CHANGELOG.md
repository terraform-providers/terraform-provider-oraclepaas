## 1.5.3 (September 05, 2019)

BUG FIXES

* `oraclepaas_mysql_service_instance` - subnet ([#70](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/70))

## 1.5.2 (May 01, 2019)

NOTES:

This release includes a Terraform SDK upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and there should not be any significant behavioural changes. ([#68](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/68))

## 1.5.1 (April 10, 2019)

BUG FIXES: 

* oraclepaas_java_service_instance - Fixed panic around `load_balancer.0.subnets` ([#67](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/67))

## 1.5.0 (March 27, 2019)

FEATURES: 

* oraclepaas_java_service_instance - `server_count` can now be updated to scale out/in `managed_servers` and `clusters`. ([#65](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/65))

## 1.4.3 (January 28, 2019)

BUG FIXES: 

* oraclepaas_java_service_instance - `weblogic_server.0.connect_string` is properly read from the config file ([#60](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/60))

UPDATES:

* Update travis.yml to go 1.11.x ([#58](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/58))

## 1.4.2 (January 23, 2019)

BUG FIXES: 

* oraclepaas_java_service_instance - `weblogic_server.0.connect_string` can now be used in place of `weblogic_server.0.database.0.name`

## 1.4.1 (December 19, 2018)

BUG FIXES:

* Added timeout support for database service instance and mysql service instance. ([#57](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/57))

## 1.4.0 (December 17, 2018)

FEATURES: 

* oraclepaas_java_service_instance - `load_balancer` support has been added ([#54](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/54))

BUG FIXES:

* oraclepaas_java_service_instance - Timeout support has been fixed ([#56](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/56))

## 1.3.2 (September 26, 2018)

* oraclepaas_java_service_instance -  `bring_your_own_license` is now set correctly ([#51](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/51))

## 1.3.1 (September 12, 2018)

IMPROVEMENTS: 

* oraclepaas_application_container - Additional oci attributes now available ([#49](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/49))

* oraclepaas_database_service_instance - Ability to set volume size for data and backup ([#48](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/48))


BUG FIXES: 

* oraclepaas_database_service_instance - `bring_your_own_license` is now set correctly ([#49](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/49))

* ip_reservations - now sent correctly to the sdk ([#49](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/49))

## 1.3.0 (July 20, 2018)

FEATURES:

* **New Resource:** `oraclepaas_application_container` ([#36](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/36))

* oraclepaas_java_service_instance - Ability to set `desired_state`. ([#33](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/33))

* oraclepaas_java_service_instance - Ability to set `assign_public_ip` ([#28](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/28))

BUG FIXES: 

* oraclepaas_java_service_instance: Oracle Traffic Director will not be provisioned unless an `oracle_traffic_director` block has been specified ([#38](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/38))

## 1.2.1 (June 21, 2018)

BUG FIXES: 

* oraclepaas_mysql_service_instance: Fix `em_agent_username` and `em_username` ([#34](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/34))

## 1.2.0 (June 20, 2018)

FEATURES: 

* **New Resource:** `oraclepaas_mysql_service_instance` ([#27](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/27))
* **New Resource:** `oraclepaas_mysql_access_rule` ([#27](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/27))

IMPROVEMENTS:

* oraclepaas_java_service_instance - Automatically provision otd when `oracle_traffic_director` block is set ([#30](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/30))

* oraclepaas_java_service_instance - Scale up/down of `weblogic_server.0.shape` is now supported ([#29](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/29))

## 1.1.1 (May 25, 2018)

IMPROVEMENTS: 

* oraclepaas_java_service_instance - Updated list of supported service versions ([#23](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/23))

## 1.1.0 (May 25, 2018)

FEATURES:

* oraclepaas_database_service_instance - Scale up and down ([#19](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/19))

* oraclepaas_database_service_instance - Set desired state ([#20](https://github.com/terraform-providers/terraform-provider-oraclepaas/issues/20))

## 1.0.0 (March 23, 2018)

FEATURES:

* **New Resource:** `oraclepaas_database_service_instance`
* **New Resource:** `oraclepaas_java_service_instance`
* **New Resource:** `oraclepaas_database_access_rules`
* **New Resource:** `oraclepaas_java_access_rules`
* **New Datasource:** `oraclepaas_database_service_instance`
