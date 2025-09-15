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

	lists, resp, err := client.Lists.GetList(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("resp.StatusCode: %d\n", resp.StatusCode)
	fmt.Printf("lists: %+v\n", lists)

}
