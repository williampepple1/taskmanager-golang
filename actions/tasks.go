// tasks.go

package actions

import (
	"taskmanager/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

func TasksList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	tasks := &models.Tasks{}
	if err := tx.All(tasks); err != nil {
		return errors.WithStack(err)
	}
	return c.Render(200, r.JSON(tasks))
}

func TaskShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	task := &models.Task{}
	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(404, err)
	}
	return c.Render(200, r.JSON(task))
}

func TaskCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	task := &models.Task{}
	if err := c.Bind(task); err != nil {
		return errors.WithStack(err)
	}
	if verr, err := tx.ValidateAndCreate(task); err != nil {
		return c.Render(422, r.JSON(verr))
	}
	return c.Render(201, r.JSON(task))
}

func TaskUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	task := &models.Task{}
	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(404, err)
	}
	if err := c.Bind(task); err != nil {
		return errors.WithStack(err)
	}
	if verr, err := tx.ValidateAndUpdate(task); err != nil {
		return c.Render(422, r.JSON(verr))
	}
	return c.Render(200, r.JSON(task))
}

func TaskDestroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	task := &models.Task{}
	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(404, err)
	}
	if err := tx.Destroy(task); err != nil {
		return errors.WithStack(err)
	}
	return c.Render(200, r.JSON(map[string]string{"message": "Task deleted successfully"}))
}
