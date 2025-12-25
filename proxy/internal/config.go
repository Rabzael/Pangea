package internal

import (
	"encoding/json"
	"os"
)

/** Config struct to hold configuration data */
type ConfigData struct {
	Http_port         int
	Cache_ttl_seconds int
	Blacklist         []string
	Log_filename      string
}

/** Configuration data */
var Config ConfigData

/** Read config from file and return a Config struct */
func ReadConfig(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		return err
	}

	return nil
}
