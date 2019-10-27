package example

import (
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"example_box": resourceTemplate(),
		},
		ConfigureFunc: func(*schema.ResourceData) (interface{}, error) {
			client, err := NewClient()
			if err != nil {
				return nil, err
			}

			return client, nil
		},
	}
}
