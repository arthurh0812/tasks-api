package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var dirname string
var authAPIService string
var tasksDirectory string
var port string

func init() {
	authAPIService = os.Getenv("AUTH_API_SERVICE_HOST")
	tasksDirectory = os.Getenv("TASKS_DIRECTORY")
	port = os.Getenv("PORT")
	var err error
	dirname, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	tasksFilePath := filepath.Join(dirname, tasksDirectory, "tasks.txt")
	tasksSep := "TASK_SPLIT"

	app := iris.Default()

	app.Use(func(ctx *context.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "POST,GET,OPTION")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		ctx.Next()
	})

	app.Get( "/", func(ctx *context.Context) {
		_, err := extractAndVerifyToken(ctx.Request().Header)
		if err != nil {
			ctx.StopWithJSON(500, APIResponse{
				Status: http.StatusInternalServerError,
				Message: fmt.Sprintf("Failed to extract Bearer token: %v", err),
			})
			return
		}
		bs, err := os.ReadFile(tasksFilePath)
		if err != nil {
			ctx.StopWithJSON(500, APIResponse{
				Status: http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to load the tasks: %v", err),
			})
			return
		}
		tasks := bytes.Split(bs, []byte(tasksSep))
		if len(tasks) > 0 {
			tasks = tasks[:len(tasks)-1]
		}

		goTasks := make([]*Task, 0, len(tasks))
		for _, t := range tasks {
			goTask := &Task{}
			err := json.Unmarshal(t, goTask)
			if err != nil {
				ctx.StopWithJSON(500, APIResponse{
					Status: http.StatusInternalServerError,
					Message: fmt.Sprintf("Failed to unmarshal task: %v", err),
				})
				return
			}
			goTasks = append(goTasks, goTask)
		}

		ctx.StopWithJSON(http.StatusOK, APIResponse{
			Message: "Tasks loaded.",
			Status: http.StatusOK,
			Count: int64(len(goTasks)),
			Data: iris.Map{
				"tasks": goTasks,
			},
		})
	})

	app.Post( "/", func(ctx *context.Context) {
		newTask := Task{}
		err := ctx.ReadJSON(&newTask)
		if err != nil {
			ctx.StopWithJSON(500, APIResponse{
				Status: http.StatusInternalServerError,
				Message: fmt.Sprintf("Failed to decode the provided JSON task: %v", err),
			})
			return
		}

		f, err := os.OpenFile(tasksFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0750)
		if err != nil {
			ctx.StopWithJSON(500, APIResponse{
				Status: http.StatusInternalServerError,
				Message: fmt.Sprintf("Failed to open the DB: %v", err),
			})
			return
		}
		enc := json.NewEncoder(f)
		err = enc.Encode(newTask)
		if err != nil {
			ctx.StopWithJSON(500, APIResponse{
				Status: http.StatusInternalServerError,
				Message: fmt.Sprintf("Failed to encode new Task: %v", err),
			})
			return
		}
		_, err = f.WriteString(tasksSep) // insert a task separator
		if err != nil {
			ctx.StopWithJSON(500, APIResponse{
				Status: http.StatusInternalServerError,
				Message: fmt.Sprintf("Failed to write the task seperator: %v", err),
			})
			return
		}
		ctx.StopWithJSON(http.StatusCreated, APIResponse{
			Message: "Task stored.",
			Status: http.StatusCreated,
			Count: 1,
			Data: iris.Map{
				"createdTask": &newTask,
			},
		})
	})

	err := app.Run(iris.Addr(fmt.Sprintf(":%s", port)))
	if err != nil {
		log.Fatal(err)
	}
}

func extractAndVerifyToken(header http.Header) (string, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("no authorization header provided")
	}
	var token string
	if parts := strings.Split(auth, " "); parts[0] != "Bearer" {
		return "", fmt.Errorf("authorization requires %q key", "Bearer")
	} else if len(parts) < 2 {
		return "", fmt.Errorf("no Bearer token provided")
	} else {
		token = parts[1]
	}
	res, err := http.Get(fmt.Sprintf("http://%s/verify-token/%s", authAPIService, token))
	if err != nil {
		return "", fmt.Errorf("request to authorization service failed: %v", err)
	}
	dec := json.NewDecoder(res.Body)
	defer res.Body.Close()
	apiRes := APIResponse{}
	err = dec.Decode(&apiRes)
	if err != nil {
		return "", fmt.Errorf("failed to decode JSON API response from the authorization service: %v", err)
	}
	if apiRes.Status % 100 > 2 {
		return "", fmt.Errorf("request to authorization service returned error: %s", apiRes.Message)
	}
	if uid := apiRes.Data["uid"]; uid == nil {
		return "", fmt.Errorf("uid not provided")
	} else {
		return uid.(string), nil
	}
}