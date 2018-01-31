package main

import (
	"github.com/hashicorp/terraform-provider-oraclepaas/oraclepaas"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: oraclepaas.Provider})
}
