package main

import (
	"github.com/ewilde/terraform-provider-kibana/kibana"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kibana.Provider})
}
