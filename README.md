# go-ticktick

Golang client for [TickTick open API](https://developer.ticktick.com/docs/index.html#/openapi?id=ticktick-open-api).

**WIP** This client is under development and at the same time, it's blocked by incomplete ticktick API.

## Get access token

Follow instructions at [ticktick developers portal](https://developer.ticktick.com/api#/openapi?id=authorization).

> Use `http://127.0.0.1:42548` as **OAuth redirect URL** in a application settings.

## Example

```go
package main

import (
	"context"
	"fmt"
	"log"

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


	task, _, err := client.Tasks.Create(ctx, &ticktick.Task{Title: "Finish go-ticktick client"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(task)
}
```
