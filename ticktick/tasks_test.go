package ticktick

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTasksService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/project/abcd/task/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"title":"test task", "projectId": "abcd", "id": "1234"}`)
	})

	task, _, err := client.Tasks.Get(context.Background(), "abcd", "1234")
	if err != nil {
		t.Errorf("Tasks.Get returned error: %v", err)
	}

	want := &Task{
		ID:        "1234",
		Title:     "test task",
		ProjectID: "abcd",
	}
	if !reflect.DeepEqual(task, want) {
		t.Errorf("Tasks.Get returned %+v, want %+v", task, want)
	}
}

func TestTasksService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Task{
		Title: "a updated task",
	}

	mux.HandleFunc("/task/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		v := new(Task)
		json.NewDecoder(r.Body).Decode(v)

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, fmt.Sprintf(`{"title": "%s"}`, input.Title))
	})

	task, _, err := client.Tasks.Update(context.Background(), "1234", input)
	if err != nil {
		t.Errorf("Tasks.Update returned error: %v", err)
	}

	if !reflect.DeepEqual(task, input) {
		t.Errorf("Tasks.Update returned %+v, want %+v", task, input)
	}
}

func TestTasksService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Task{
		Title: "a new test task",
	}

	mux.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		v := new(Task)
		json.NewDecoder(r.Body).Decode(v)

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"title": "a new test task"}`)
	})

	task, _, err := client.Tasks.Create(context.Background(), input)
	if err != nil {
		t.Errorf("Tasks.Create returned error: %v", err)
	}

	if !reflect.DeepEqual(task, input) {
		t.Errorf("Tasks.Create returned %+v, want %+v", task, input)
	}
}

func TestTasksService_Complete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/project/abcd/task/1234/complete", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(200)
	})

	_, err := client.Tasks.Complete(context.Background(), "abcd", "1234")
	if err != nil {
		t.Errorf("Tasks.Complete returned error: %v", err)
	}
}
