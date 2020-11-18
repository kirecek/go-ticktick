package ticktick

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", mux)

	server := httptest.NewServer(apiHandler)

	client = NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}

	c2 := NewClient(nil)
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	type testTask struct {
		Name string `json:"name"`
	}

	inURL, outURL := "/task", defaultBaseURL+"task"
	inBody, outBody := &testTask{Name: "Test"}, `{"name":"Test"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("Unexpected User-Agent; got %v, want %v", got, want)
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type testTask struct {
		Name string `json:"name"`
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"name":"test"}`)
	})

	req, err := client.NewRequest("GET", ".", nil)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	body := new(testTask)
	resp, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	t.Log(resp)
	want := &testTask{"test"}

	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	resp, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("Expected not nil error.")
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

func TestCheckResponse_badRequest(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
	}

	err := CheckResponse(res).(*ErrorResponse)
	if err == nil {
		t.Errorf("Expected error response, got nil.")
	}

	want := &ErrorResponse{Response: res}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}
