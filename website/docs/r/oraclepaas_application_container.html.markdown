---
layout: "oraclepaas"
page_title: "Oracle: oraclepaas_application_container"
sidebar_current: "docs-oraclepaas-resource-application-container"
description: |-
  Creates and manages an Appliction Container.

---

# oraclepaas_application_container

The `oraclepaas_application_container` resource creates and manages an Application Container.

## Example Usage

```hcl
resource "oraclepaas_application_container" "example-app" {
  name               = "ExampleWebApp"
  runtime            = "java"
  archive_url        = "my-accs-apps/example-web-app.zip"
  subscription_type  = "HOURLY"

  deployment {
    memory = "1G"
    instances = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Application Container.

* `manifest_file` - (Optional) The json manifest file containing the attributes related to launching an application. Use either `manifest_file` or `manifest_attributes` when specifying 
launch information.

* `manifest` - (Optional) The manifest attributes related to launching an application. Use either `manifest_file` or `manifest` when specifying 
launch information. Manifest attributes is documented below.

* `deployment_file` - (Optional) The json deployment file containing the attributes related to deploying an application. Use either `deployment_file` or `deployment_attributes` when specifying
deployment information. 

* `deployment` - (Optional) The deployment attributes related to deploying an application. Use either `deployment_file` or `deployment` when specifying
deployment information. Deployment attributes is documented below.

* `archive_url` - (Optional) Location of the application archive file in Oracle Storage Cloud Service, in the format app-name/file-name.

* `auth_type` - (Optional) Uses Oracle Identity Cloud Service to control who can access your Java SE 7 or 8, Node.js, or PHP application. Allowed values are `basic` and `oauth`.

* `git_repository` - (Optional) The URL of the git repository to use the application container.

* `git_username` - (Optional) The username of a user with access to the git respository if the repository is private.

* `git_password` - (Optional) The password for the user with access to the git repository if the repository is private.

* `notes` - (Optional) Comments about the application deployment.

* `notification_email` - (Optional) Email address to which application deployment status updates are sent.

* `repository` (Optional) Repository of the application. The only allowed value is 'dockerhub'.

* `runtime` - (Optional) The allowed runtime environment variables. The allowed variables are `java`, `node`, `php`, `python`, `golang`, `dotnet`, or `ruby`. The default is `java`.

* `subscription_type` - (Optional) Whether the subscription type is `hourly` or `monthly`. The default is `hourly`.

* `tags` - (Optional) A map of tags for the application container.

Manifest attributes supports the following: 

* `runtime` - (Optional) Details the availble runtime attributes. Runtime is documented below.

* `type` - (Optional) Determines whether the application is public or private. The default is `worker` (private).

* `command` - (Optional) Launch command to execute after the application has been uploaded.

* `release` - (Optional) Details the release attributes of a specific build. Release is documented below.

* `startup_time` - (Optional) The maximum time in seconds to wait for an application to start.

* `shutdown_time` - (Optional) The maximum time in seconds to wait for an application to stop.

* `notes` - (Optional) Comments about the launch configuration.

* `mode` - (Optional) The restart mode for application instances when the application is restarted. The only allowed value is `rolling`.

* `clustered` - (Optional) Boolean for whether the application instances act as a cluster with failover capability.

* `home` - (Optional) The context root of the application.

* `health_check_endpoint` - (Optional) The URL that the application uses for health checks.

Deployment attributes supports the following: 

* `memory` - (Optional) The amount of memory in gigabytes made available to the application. The default is `2G`. 

* `instances` - (Optional) The number of application instances. The default is `2`.

* `notes` - (Optional) Comments about the deployment.

* `environment` - (Optional) A map of environment variables used by the application.

* `secure_environment` - (Optional) A list of environment variables marked as secured on the user interface.

* `java_system_properties` - (Optional) A map os java system properties used by the application.

* `services` - (Optional) Service bindings for connections to other Oracle Cloud services. Services is documented below.

Runtime supports the following:

* `major_version` - (Required) The major version of the runtime environment.

Release supports the following:

* `build` - (Optional) The value for a specific build.

* `commit` - (Optional) The value for a specific commit.

* `version` - (Optional) The value for a specific version.

Services supports the following:

* `identifier` - (Required) The value for the identifier

* `type` - (Required) The type of service. Allowed values are `JAAS`, `DBAAS`, `MYSQLCS`, `OEHCS`, `OEHPCS`, `DHCS`, `caching`.

* `name` - (Required) The name of the existing service. 

* `username` - (Required) The username to connect to the service.

* `password` - (Required) The password to connect to the service.

In addition to the above, the following values are exported: 

* `app_url` - URL of the created application

* `web_url` - Web URL of the application