# go-ticktick

Golang client for [TickTick open API](https://developer.ticktick.com/docs/index.html#/openapi?id=ticktick-open-api).

**WIP** This client is under development and at the same time, it's blocked by incomplete ticktick API.

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