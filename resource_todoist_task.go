package main

import (
	"github.com/eddiezane/todoist-rest-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTodoistTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTodoistTaskCreate,
		Read:   resourceTodoistTaskRead,
		Update: resourceTodoistTaskUpdate,
		Delete: resourceTodoistTaskDelete,

		Schema: map[string]*schema.Schema{
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"completed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceTodoistTaskCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*todoist.Client)
	content := d.Get("content").(string)

	newTask := &todoist.NewTask{
		Content: content,
	}

	task, err := client.CreateTask(newTask)
	if err != nil {
		return err
	}

	d.SetId(task.Id)
	return nil
}

func resourceTodoistTaskRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*todoist.Client)
	id := d.Id()

	task, err := client.GetTask(id)
	if err != nil {
		return err
	}

	d.Set("content", task.Content)

	return nil
}

func resourceTodoistTaskUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*todoist.Client)
	id := d.Id()
	content := d.Get("content").(string)

	b := d.HasChange("completed")
	if b {
		old, cur := d.GetChange("completed")
		// Close task
		if old == false && cur == true {
			err := client.CloseTask(id)
			if err != nil {
				return err
			}
		}
		// Reopen task
		if old == true && cur == false {
			err := client.ReopenTask(id)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("content") {
		t := &todoist.Task{
			Id:      id,
			Content: content,
		}
		err := client.UpdateTask(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTodoistTaskDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*todoist.Client)
	id := d.Id()

	err := client.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}
