package main

import (
	"fmt"

	apiKey "github.com/inuoshios/generate-api-key/generate"
)

func main() {
	app := apiKey.Initialize()

	result, err := app.GenerateAPIKey(apiKey.GenerateKeyOptions{
		Prefix: "pk",
		Batch:  5,
		Method: "base62",
	})

	if err != nil {
		panic(err)
	} else {
		fmt.Println(result)
	}

}
