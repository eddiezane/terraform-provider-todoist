package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Your Todoist API key",
				Required:    true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
	}
}
