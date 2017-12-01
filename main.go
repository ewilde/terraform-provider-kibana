package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/ewilde/terraform-provider-kibana/kibana"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kibana.Provider})
}
