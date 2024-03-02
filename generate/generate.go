package generateapikey

import (
	crypto "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
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

// GenerateKeyOptions is the options used to generate an API key
type GenerateKeyOptions struct {
	// The length of the API key
	Length uint32

	// The characters used for the API key generation
	Pool string

	// A string prefix for the API key, followed by a period (.)
	Prefix string

	// The number of API keys to generate
	Batch uint32

	// Add dashes (-) to the API key or not
	Dashes bool

	// The method used to generate the API key (string, bytes, base32, base62, uuidv4, uuidv5)
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

// splitString splits a string into parts of given length
func splitString(s string, length int) []string {
	var parts []string
	for len(s) >= length {
		parts = append(parts, s[:length])
		s = s[length:]
	}
	if len(s) > 0 {
		parts = append(parts, s)
	}
	return parts
}

// base32Stringify converts a slice of integers to base32 string
func base32Stringify(numArr []int) string {
	base32Alphabet := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	var result strings.Builder
	for _, num := range numArr {
		result.WriteByte(base32Alphabet[num%32])
	}
	return result.String()
}

func generateBase32(options GenerateKeyOptions) (any, error) {
	if options.Pool != "" {
		return nil, fmt.Errorf("pool is not supported for base32 method")
	}
	if options.Length != 0 {
		return nil, fmt.Errorf("length is not supported for base32 method")
	}

	// create a new uuid
	uuid := uuid.New()
	// split the uuid into four parts
	uuidParts := strings.Split(uuid.String(), "-")
	println(uuid.String())

	// convert the uuid into 4 equally separate parts
	partsArr := []string{
		uuidParts[0],
		fmt.Sprintf("%s%s", uuidParts[1], uuidParts[2]),
		fmt.Sprintf("%s%s", uuidParts[3], uuidParts[4][:4]),
		uuidParts[4][4:],
	}

	if options.Batch > 1 {
		var batchResults []string

		for i := uint32(0); i < options.Batch; i++ {
			var apiKeyArr []string
			var finalKey string
			for _, value := range partsArr {
				// Get every two characters
				valueArr := splitString(value, 2)
				// Convert each value into a number
				var numArr []int
				for _, item := range valueArr {
					num, _ := hex.DecodeString(item)
					numArr = append(numArr, int(num[0]))
				}
				// Create the string
				apiKeyArr = append(apiKeyArr, base32Stringify(numArr))
			}

			// Check if we should add dashes
			apiKey := strings.Join(apiKeyArr, "-")
			if options.Dashes {
				finalKey = apiKey
			} else {
				finalKey = strings.ReplaceAll(apiKey, "-", "")
			}
			batchResults = append(batchResults, finalKey)
		}

		return batchResults, nil
	} else {
		// Iterate through each part and convert to base32
		var apiKeyArr []string
		for _, value := range partsArr {
			// Get every two characters
			valueArr := splitString(value, 2)
			// Convert each value into a number
			var numArr []int
			for _, item := range valueArr {
				num, _ := hex.DecodeString(item)
				numArr = append(numArr, int(num[0]))
			}
			// Create the string
			apiKeyArr = append(apiKeyArr, base32Stringify(numArr))
		}

		// Check if we should add dashes
		apiKey := strings.Join(apiKeyArr, "-")
		if options.Dashes {
			return apiKey, nil
		}

		return strings.ReplaceAll(apiKey, "-", ""), nil
	}
}

// GenerateAPIKey generates an API key based on the options provided
func (*Generate) GenerateAPIKey(options GenerateKeyOptions) (any, error) {
	switch options.Method {
	case StringOption:
		return generateString(options)
	case BytesOption:
		return generateByte(options)
	case Base32Option:
		return generateBase32(options)
	default:
		return nil, fmt.Errorf("unsupported method %s", options.Method)
	}
}
