package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/arthurh0812/tasks-api/pkg/api/middleware"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arthurh0812/tasks-api/pkg/api"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type TaskHandler struct {
	tasksFilePath string
	taskSep string
	env map[string]string
}

func NewTaskHandler(env map[string]string) *TaskHandler {
	return &TaskHandler{
		env: env,
		taskSep: "T_SPLIT",
		tasksFilePath: filepath.Join(env["PWD"], env["TASKS_DIRECTORY"], "tasks.txt"),
	}
}

func (t *TaskHandler) GetTasks(ctx *context.Context) {
	_, err := middleware.ExtractBearerToken(ctx, t.env["AUTH_API_SERVICE_HOST"])
	if err != nil {
		ctx.StopWithJSON(500, api.Response{
			Status: http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to extract Bearer token: %v", err),
		})
		return
	}
	bs, err := os.ReadFile(t.tasksFilePath)
	if err != nil {
		ctx.StopWithJSON(500, api.Response{
			Status: http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to load the tasks: %v", err),
		})
		return
	}
	tasks := bytes.Split(bs, []byte(t.taskSep))
	if len(tasks) > 0 {
		tasks = tasks[:len(tasks)-1]
	}

	goTasks := make([]*api.Task, 0, len(tasks))
	for _, t := range tasks {
		goTask := &api.Task{}
		err := json.Unmarshal(t, goTask)
		if err != nil {
			ctx.StopWithJSON(500, api.Response{
				Status: http.StatusInternalServerError,
				Message: fmt.Sprintf("Failed to unmarshal task: %v", err),
			})
			return
		}
		goTasks = append(goTasks, goTask)
	}

	ctx.StopWithJSON(http.StatusOK, api.Response{
		Message: "Tasks loaded.",
		Status: http.StatusOK,
		Count: int64(len(goTasks)),
		Data: iris.Map{
			"tasks": goTasks,
		},
	})
}

func (t *TaskHandler) CreateTask(ctx *context.Context) {
	newTask := api.Task{}
	err := ctx.ReadJSON(&newTask)
	if err != nil {
		ctx.StopWithJSON(500, api.Response{
			Status: http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to decode the provided JSON task: %v", err),
		})
		return
	}

	f, err := os.OpenFile(t.tasksFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0750)
	if err != nil {
		ctx.StopWithJSON(500, api.Response{
			Status: http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to open the DB: %v", err),
		})
		return
	}
	enc := json.NewEncoder(f)
	err = enc.Encode(newTask)
	if err != nil {
		ctx.StopWithJSON(500, api.Response{
			Status: http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to encode new Task: %v", err),
		})
		return
	}
	_, err = f.WriteString(t.taskSep) // insert a task separator
	if err != nil {
		ctx.StopWithJSON(500, api.Response{
			Status: http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to write the task seperator: %v", err),
		})
		return
	}
	ctx.StopWithJSON(http.StatusCreated, api.Response{
		Message: "Task stored.",
		Status: http.StatusCreated,
		Count: 1,
		Data: iris.Map{
			"createdTask": &newTask,
		},
	})
}
