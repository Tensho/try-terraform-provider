package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	// NOTE: For Go modules I have to specify the full path
	"example.com/terraform-provider-example/example"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: example.Provider,
	})
}
