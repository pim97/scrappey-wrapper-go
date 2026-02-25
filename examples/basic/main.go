package main

import (
	"context"
	"fmt"
	"log"
	"os"

	scrappey "github.com/scrappey/wrapper-go"
)

func main() {
	apiKey := os.Getenv("SCRAPPEY_API_KEY")
	if apiKey == "" {
		log.Fatal("set SCRAPPEY_API_KEY before running this example")
	}

	client, err := scrappey.NewClient(apiKey, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	targetURL := "https://httpbin.org/get"

	response, err := client.Get(context.Background(), scrappey.RequestOptions{
		"url":         targetURL,
	})
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}

	fmt.Printf("Data: %s\n", response.Data)
	fmt.Printf("Status: %d\n", response.SolutionInt("statusCode"))

	currentURL := response.SolutionString("currentUrl")
	if currentURL == "" {
		// request mode often does not return currentUrl; this is expected.
		currentURL = targetURL
	}
	fmt.Printf("Current URL: %s\n", currentURL)
}
