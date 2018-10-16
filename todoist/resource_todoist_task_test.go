package todoist

import (
	"fmt"
	"testing"

	todoistRest "github.com/eddiezane/todoist-rest-go"
	// "github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTodoistTask_basic(t *testing.T) {
	var task todoistRest.Task
	content := "test task content"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTodoistTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTodoistTaskConfig_basic(content),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTodoistTaskExists("todoist_task.test", &task),
					testAccCheckTodoistTaskAttributes(&task, content),
					resource.TestCheckResourceAttr("todoist_task.test", "content", content),
				),
			},
		},
	})
}

func TestAccTodoistTask_update(t *testing.T) {
	var task todoistRest.Task
	content := "test task content"
	newContent := content + " new"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTodoistTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTodoistTaskConfig_basic(content),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTodoistTaskExists("todoist_task.test", &task),
					testAccCheckTodoistTaskAttributes(&task, content),
					resource.TestCheckResourceAttr("todoist_task.test", "content", content),
				),
			},
			{
				Config: testAccCheckTodoistTaskConfig_basic(newContent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTodoistTaskExists("todoist_task.test", &task),
					testAccCheckTodoistTaskAttributes(&task, newContent),
					resource.TestCheckResourceAttr("todoist_task.test", "content", newContent),
				),
			},
		},
	})
}

func TestAccTodoistTask_completed(t *testing.T) {
	var task todoistRest.Task
	content := "test task content"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTodoistTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTodoistTaskConfig_completed(content, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTodoistTaskExists("todoist_task.test", &task),
					testAccCheckTodoistTaskAttributes(&task, content),
					resource.TestCheckResourceAttr(
						"todoist_task.test", "content", content),
				),
			},
			{
				Config: testAccCheckTodoistTaskConfig_completed(content, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTodoistTaskExists_completed("todoist_task.test", &task),
					testAccCheckTodoistTaskAttributes(&task, content),
					resource.TestCheckResourceAttr("todoist_task.test", "content", content),
					resource.TestCheckResourceAttr("todoist_task.test", "completed", "true"),
				),
			},
			{
				Config: testAccCheckTodoistTaskConfig_completed(content, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTodoistTaskExists_completed("todoist_task.test", &task),
					testAccCheckTodoistTaskAttributes(&task, content),
					resource.TestCheckResourceAttr("todoist_task.test", "content", content),
					resource.TestCheckResourceAttr("todoist_task.test", "completed", "false"),
				),
			},
		},
	})
}

func testAccCheckTodoistTaskDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*todoistRest.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "todoist_task" {
			continue
		}

		// Try to find the task
		_, err := client.GetTask(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Task still exists")
		}

		_, err = client.GetCompletedTask(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Task still exists but is completed")
		}
	}

	return nil
}

func testAccCheckTodoistTaskAttributes(task *todoistRest.Task, content string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if task.Content != content {
			return fmt.Errorf("Content does not match: %s", task.Content)
		}

		return nil
	}
}

func testAccCheckTodoistTaskExists(n string, task *todoistRest.Task) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*todoistRest.Client)

		foundTask, err := client.GetTask(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundTask.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*task = *foundTask

		return nil
	}
}

func testAccCheckTodoistTaskExists_completed(n string, task *todoistRest.Task) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*todoistRest.Client)

		t, err := client.GetTask(rs.Primary.ID)
		if err != nil {
			completedTask, err2 := client.GetCompletedTask(rs.Primary.ID)

			if err2 != nil {
				return err
			}

			if completedTask.Id != rs.Primary.ID {
				return fmt.Errorf("Record not found")
			}

			*task = todoistRest.Task{
				Id:        completedTask.Id,
				Content:   completedTask.Content,
				Completed: true,
			}
		} else {
			*task = *t
		}

		return nil
	}
}

func testAccCheckTodoistTaskConfig_basic(content string) string {
	return fmt.Sprintf(`
	resource "todoist_task" "test" {
		content = "%s"
	}
	`, content)
}

func testAccCheckTodoistTaskConfig_completed(content string, completed bool) string {
	return fmt.Sprintf(`
	resource "todoist_task" "test" {
		content = "%s"
		completed = %t
	}
	`, content, completed)
}
