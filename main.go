package main

import (
	"github.com/monitobeko/terraform-provider-kibana/kibana"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kibana.Provider,
	})
}
