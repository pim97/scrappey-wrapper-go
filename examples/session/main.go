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

	ctx := context.Background()

	sessionResponse, err := client.CreateSession(ctx, scrappey.SessionOptions{
		"proxyCountry": "UnitedStates",
		"premiumProxy": true,
	})
	if err != nil {
		log.Fatalf("create session failed: %v", err)
	}
	sessionID := sessionResponse.Session
	fmt.Printf("Created session: %s\n", sessionID)

	requestResponse, err := client.Get(ctx, scrappey.RequestOptions{
		"url":     "https://httpbin.org/cookies",
		"session": sessionID,
	})
	if err != nil {
		log.Fatalf("session request failed: %v", err)
	}
	fmt.Printf("Session request status: %d\n", requestResponse.SolutionInt("statusCode"))

	if _, err := client.DestroySession(ctx, sessionID); err != nil {
		log.Fatalf("destroy session failed: %v", err)
	}
	fmt.Printf("Destroyed session: %s\n", sessionID)
}
