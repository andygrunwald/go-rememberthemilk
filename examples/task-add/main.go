package main

import (
	"context"
	"fmt"

	"github.com/andygrunwald/go-rememberthemilk"
)

func main() {
	apiKey := "<your-api-key>"
	sharedSecret := "<your-shared-secret>"
	token := "<your-authentication-token>"

	client := rememberthemilk.NewClient(apiKey, sharedSecret, token, nil)

	// Create a new timeline
	timeline, resp, err := client.Timelines.Create(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("resp.StatusCode: %d\n", resp.StatusCode)
	fmt.Printf("timeline: %+v\n", timeline)

	// Add a new task
	taskName := "This is a dummy task ^Today #some-tag !1 https://www.google.com/"
	input := rememberthemilk.TaskInput{
		Timeline: timeline,
		ListID:   "<your-list-id>",
		Name:     taskName,
		Parse:    rememberthemilk.ParseWithSmartAdd,
	}
	task, resp, err := client.Tasks.Add(context.Background(), input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("resp.StatusCode: %d\n", resp.StatusCode)
	fmt.Printf("task: %+v\n", task)
}
