package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	appError "app/error"

	"github.com/dotenv-org/godotenvvault"
)

var (
	AppConfig  map[string]string
	configOnce sync.Once
)

func LoadConfig() {
	configOnce.Do(func() {
		LoadEnv()
		LoadJson()
	})
}

func LoadEnv() {
	var err error
	env := os.Getenv("env")
	if env != "" {
		err = godotenvvault.Overload(fmt.Sprintf(".%s.%s", "env", env))
	} else {
		err = godotenvvault.Overload()
	}
	CheckError(err, appError.NewApplicationError("Fail to load env"))
}

func LoadJson() {
	_, filename, _, _ := runtime.Caller(0)
	configDir := filepath.Join(filepath.Dir(filename), "/../config")
	data, err := os.ReadFile(fmt.Sprintf("%s/%s", configDir, "config.json"))
	CheckError(err, appError.NewApplicationError("Fail to load config json"))

	rawConfig := make(map[string]interface{})
	err = json.Unmarshal(data, &rawConfig)
	CheckError(err, appError.NewApplicationError("Fail to load config json"))

	appConfig := make(map[string]string)
	FlattenJson(rawConfig, "", appConfig)
	AppConfig = appConfig
}

func FlattenJson(data map[string]interface{}, prefix string, flatMap map[string]string) {
	for key, value := range data {
		// Create the full key path (e.g., "parent.child" for nested keys)
		fullKey := prefix + key

		// Check the type of value
		switch v := value.(type) {
		case map[string]interface{}: // If it's a nested object, recurse
			FlattenJson(v, fullKey+"_", flatMap)
		case string: // If it's a string, add it to the flat map
			flatMap[fullKey] = v
		case float64: // If it's a number, convert to string and add
			flatMap[fullKey] = fmt.Sprintf("%v", v)
		case bool: // If it's a boolean, convert to string and add
			flatMap[fullKey] = fmt.Sprintf("%v", v)
		default:
			GlobalLog.Printf("Unsupported data type for key %s: %T\n", fullKey, v)
		}
	}
}
