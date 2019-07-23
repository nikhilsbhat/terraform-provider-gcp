package main

import (
	"terraform-provider-gcp/gcp"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gcp.Provider})
}
