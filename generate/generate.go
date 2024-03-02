package generateapikey

import (
	crypto "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand"
)

type Generate struct{}

func Initialize() *Generate {
	return &Generate{}
}

type GenerateKeyMethod string

const (
	StringOption GenerateKeyMethod = "string"
	BytesOption  GenerateKeyMethod = "bytes"
	Base32Option GenerateKeyMethod = "base32"
	Base62Option GenerateKeyMethod = "base62"
	UUIDV4Option GenerateKeyMethod = "uuidv4"
	UUIDV5Option GenerateKeyMethod = "uuidv5"
)

type GenerateKeyOptions struct {
	Length uint32
	Pool   string
	Prefix string
	Batch  uint32
	Dashes bool
	Method GenerateKeyMethod
}

func generateString(options GenerateKeyOptions) (any, error) {
	if options.Pool == "" {
		options.Pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~+/"
	}
	if options.Length == 0 {
		options.Length = 36
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

func generateRandomBytes(length uint32) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := crypto.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func generateByte(options GenerateKeyOptions) (any, error) {
	if options.Pool != "" {
		return nil, fmt.Errorf("pool is not supported for byte method")
	}
	if options.Length == 0 {
		options.Length = 36
	}
	if options.Prefix == "" {
		options.Prefix = ""
	} else {
		options.Prefix = fmt.Sprintf("%s.", options.Prefix)
	}

	if options.Batch > 1 {
		var batchResults []string
		for i := uint32(0); i < options.Batch; i++ {
			bytes, err := generateRandomBytes(options.Length)
			if err != nil {
				return nil, err
			}
			batchResults = append(batchResults, options.Prefix+hex.EncodeToString(bytes))
		}
		return batchResults, nil
	} else {
		var result string
		bytes, err := generateRandomBytes(options.Length)
		if err != nil {
			return nil, err
		}
		result = options.Prefix + hex.EncodeToString(bytes)
		return result, nil
	}
}

func (*Generate) GenerateString(options GenerateKeyOptions) (any, error) {
	switch options.Method {
	case StringOption:
		return generateString(options)
	case BytesOption:
		return generateByte(options)
	case Base32Option:
		// return generateBase32(options)
	default:
		return nil, fmt.Errorf("unsupported method %s", options.Method)
	}
}
