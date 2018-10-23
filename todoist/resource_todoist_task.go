package todoist

import (
	todoistRest "github.com/eddiezane/todoist-rest-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTodoistTask() *schema.Resource {
	return &schema.Resource{
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
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				// Hack because Todoist API doesn't let you change the project_id
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "0"
				},
			},
		},
		Create: resourceTodoistTaskCreate,
		Read:   resourceTodoistTaskRead,
		Update: resourceTodoistTaskUpdate,
		Delete: resourceTodoistTaskDelete,
	}
}

func resourceTodoistTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*todoistRest.Client)
	content := d.Get("content").(string)

	newTask := &todoistRest.NewTask{
		Content: content,
	}

	if project_id, ok := d.GetOk("project_id"); ok {
		newTask.ProjectId = project_id.(int)
	}

	task, err := client.CreateTask(newTask)
	if err != nil {
		return err
	}

	d.SetId(task.Id)
	return resourceTodoistTaskRead(d, meta)
}

func resourceTodoistTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*todoistRest.Client)
	id := d.Id()
	var t todoistRest.Task

	task, err := client.GetTask(id)
	if err == nil {
		t.Content = task.Content
		t.Completed = task.Completed
	} else {
		completedTask, err2 := client.GetCompletedTask(id)
		if err2 != nil {
			return err
		}

		t.Content = completedTask.Content
		t.Completed = true
	}

	d.Set("content", t.Content)
	d.Set("completed", t.Completed)
	d.Set("project_id", t.ProjectId)

	return nil
}

func resourceTodoistTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*todoistRest.Client)
	id := d.Id()
	content := d.Get("content").(string)
	// project_id := d.Get("project_id").(int)

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
		t := &todoistRest.Task{
			Id:      id,
			Content: content,
		}
		err := client.UpdateTask(t)
		if err != nil {
			return err
		}
	}

	// Update project_id unsupported by API
	// if d.HasChange("project_id") {
	// t := &todoistRest.Task{
	// Id:        id,
	// Content:   content,
	// ProjectId: project_id,
	// }
	// err := client.UpdateTask(t)
	// if err != nil {
	// return err
	// }
	// }

	return nil
}

func resourceTodoistTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*todoistRest.Client)
	id := d.Id()

	err := client.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}
