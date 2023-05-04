package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/asiless/interview-accountapi/pkg/account"
)

func main() {

	cfg := account.Config{
		Host:    "localhost",
		Port:    8080,
		Version: "v1",
	}

	client := account.NewAccountClient(&cfg)

	fmt.Println(client.Greet())

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := client.FetchAccount(ctx, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	if err != nil {
		fmt.Println("error: %w", err)
	} else {
		j, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Println("error: %w", err)
		}
		fmt.Printf("accounts => %+v", string(j))
	}
}
