## 1.3.0 (Unreleased)
## 1.2.1 (June 21, 2018)

BUF FIXES: 

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
