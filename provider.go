package main

import (
	"github.com/eddiezane/todoist-rest-go"
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
		ResourcesMap: map[string]*schema.Resource{
			"todoist_task": resourceTodoistTask(),
		},
		ConfigureFunc: configureFunc(),
	}
}

func configureFunc() func(*schema.ResourceData) (interface{}, error) {
	return func(rd *schema.ResourceData) (interface{}, error) {
		client := todoist.NewClient(rd.Get("api_key").(string))
		return client, nil
	}
}
