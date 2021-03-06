package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-oraclepaas/oraclepaas"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: oraclepaas.Provider})
}
