package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"./example"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: example.Provider,
	})
}
