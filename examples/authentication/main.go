package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/andygrunwald/go-rememberthemilk"
)

func main() {
	apiKey := "<your-api-key>"
	sharedSecret := "<your-shared-secret>"

	// Initialize the client without an authentication token.
	client := rememberthemilk.NewClient(apiKey, sharedSecret, "", nil)

	// Step 1: Get a frob
	frob, _, err := client.Authentication.GetFrob(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Frob: %s\n", frob)

	// Step 2: Get the authentication URL and prompt the user to visit it
	address, err := client.Authentication.GetAuthenticationURLWithFrob(rememberthemilk.PermissionWrite, frob)
	if err != nil {
		panic(err)
	}
	fmt.Println("Please visit the following URL to authenticate the application:")
	fmt.Println(address)

	fmt.Println("Press Enter after you have authenticated...")
	reader := bufio.NewReader(os.Stdin)
	_, _, err = reader.ReadRune()
	if err != nil {
		panic(err)
	}

	// Step 3: Get the authentication token using the frob
	token, _, err := client.Authentication.GetToken(context.Background(), frob)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Authentication successful! Token: %s\n", token.Token)

	// Set the authentication token in the client for future requests
	client.SetAuthenticationToken(token.Token)

	// Now you can use the client with the authenticated token.
}
