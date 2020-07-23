package main

import (
	"github.com/red-gate/terraform-provider-ad/ad"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ad.Provider})
}
