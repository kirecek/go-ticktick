package ticktick

import (
	"context"
	"fmt"
	"net/http"
)

// TasksService handles communication with "task" related methods.
type TasksService service

// Task represents a TickTick task.
//
// API spec: https://developer.ticktick.com/api#/openapi?id=task-1
type Task struct {
	ID          string `json:"id,omitempty"`
	ProjectID   string `json:"projectId,omitempty"`
	AllDay      bool   `json:"allDay,omitempty"`
	Title       string `json:"title,omitempty"`
	Content     string `json:"content,omitempty"`
	Description string `json:"desc,omitempty"`
}

// Get a single TickTick task from the given project.
func (s *TasksService) Get(ctx context.Context, projectID, taskID string) (*Task, *http.Response, error) {
	u := fmt.Sprintf("project/%s/task/%s", projectID, taskID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	task := new(Task)
	resp, err := s.client.Do(ctx, req, task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}

// Create a new task.
func (s *TasksService) Create(ctx context.Context, task *Task) (*Task, *http.Response, error) {
	u := "task"

	req, err := s.client.NewRequest("POST", u, task)
	if err != nil {
		return nil, nil, err
	}

	t := new(Task)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// Update an existing task.
func (s *TasksService) Update(ctx context.Context, taskID string, task *Task) (*Task, *http.Response, error) {
	u := fmt.Sprintf("task/%s", taskID)

	req, err := s.client.NewRequest("POST", u, task)
	if err != nil {
		return nil, nil, err
	}

	t := new(Task)
	resp, err := s.client.Do(ctx, req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, nil
}

// Complete a task in a given project.
func (s *TasksService) Complete(ctx context.Context, projectID, taskID string) (*http.Response, error) {
	u := fmt.Sprintf("project/%s/task/%s/complete", projectID, taskID)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
