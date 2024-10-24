package config

import (
	"encoding/json"
	"os"

	"github.com/48club/rpc-watchdog/types"
)

var Config types.Config

func init() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
}
