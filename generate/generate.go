package generateapikey

import (
	"fmt"
	"math/rand"
)

type Generate struct{}

func Initialize() *Generate {
	return &Generate{}
}

type StringGenOptions struct {
	Length uint32
	Pool   string
	Prefix string
	Batch  uint32
}

func (*Generate) GenerateString(options StringGenOptions) (any, error) {
	if options.Pool == "" {
		options.Pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~+/"
	}
	if options.Length == 0 {
		options.Length = 10
	}
	if options.Prefix == "" {
		options.Prefix = ""
	} else {
		options.Prefix = fmt.Sprintf("%s.", options.Prefix)
	}

	if options.Batch > 1 {
		var batchResults []string
		for i := uint32(0); i < options.Batch; i++ {
			var result string
			result = options.Prefix
			for j := uint32(0); j < options.Length; j++ {
				result += string(options.Pool[rand.Intn(len(options.Pool))])
			}
			batchResults = append(batchResults, result)
		}
		return batchResults, nil
	} else {
		var result string
		result = options.Prefix
		for i := uint32(0); i < options.Length; i++ {
			result += string(options.Pool[rand.Intn(len(options.Pool))])
		}

		return result, nil
	}
}
