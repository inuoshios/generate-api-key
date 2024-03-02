package main

import (
	"fmt"

	apiKey "github.com/inuoshios/generate-api-key/generate"
)

func main() {
	app := apiKey.Initialize()

	result, err := app.GenerateString(apiKey.StringGenOptions{
		Pool:   "ABCDEFG1234567890",
		Prefix: "pk",
		Length: 20,
		Batch:  1,
	})
	if err != nil {
		fmt.Println("An error occurred")
	}

	fmt.Println(result)
}
