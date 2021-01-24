# WIP: go-ticktick

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/kirecek/go-ticktick/Tests?label=tests&logo=github)
[![Go Report Card](https://goreportcard.com/badge/github.com/kirecek/go-ticktick)](https://goreportcard.com/report/github.com/kirecek/go-ticktick)

Golang client for [TickTick open API](https://developer.ticktick.com/docs/index.html#/openapi?id=ticktick-open-api).

> this is unofficial, experimental ticktick Go client library.

> TickTick API is missing a lot of functionality, therefore adding new methods and features to this library is very limited. Please check their RESP API [support ticket](http://help.ticktick.com/forum/topic/393284) for more information ... maybe and also add +1 to increase demand :)

## Get access token

Follow instructions at [ticktick developers portal](https://developer.ticktick.com/api#/openapi?id=authorization).

> Use `http://127.0.0.1:42548` as **OAuth redirect URL** in a application settings.

## Examples

### Auth

```go
package main

import (
	"context"

	"github.com/kirecek/go-ticktick/ticktick"
)

func main() {
	ctx := context.Background()

	authConfig := &ticktick.OAuthConfig{
		Scopes:       []string{ticktick.ScopeReadTask, ticktick.ScopeWriteTask},
		ClientID:     "<client-id>",
		ClientSecret: "<client-secret>",
	}
	client := ticktick.NewOAuthClient(ctx, authConfig)
}
```

### Get Task

```go
task, _, err := client.Tasks.Get(ctx, "inbox", "<task-id>")
if err != nil {
	log.Fatal(err)
}
fmt.Println(task)
```

### Create task

```go
task, _, err := client.Tasks.Create(ctx, &ticktick.Task{Title: "Testing go-ticktick client"})
if err != nil {
	log.Fatal(err)
}
```

### Complete task

```go
_, err := client.Tasks.Complete(ctx, task.ProjectID, task.ID)
if err != nil {
	log.Fatal(err)
}
```
