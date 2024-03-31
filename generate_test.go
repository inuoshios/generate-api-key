package generateapikey

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("Test string option", func(t *testing.T) {
		gen, _ := Initialize().GenerateAPIKey(GenerateKeyOptions{
			Method: "string",
			Prefix: "pk",
		})
		fmt.Printf("%s", gen)
	})
}
