package config

import (
    "bitbucket.org/gismart/config"
    "bitbucket.org/gismart/config/source/file"
    "bitbucket.org/gismart/config/source/env"
    "bitbucket.org/gismart/config/source/flag"
)

const pth = "PATH/TO/CONFIG.JSON"

var Config = schema{}

func init() {
    loadConfig(&Config)
}

func loadConfig(cfg *schema) {
    // By default, the loader loads all the keys from the environment.
    // The loader can take other configuration source as parameters.
    loader := config.NewLoader(
        // IMPORTANT: sources should be provided in accordance
        // with the priority from the minor to the major
        file.New(pth),
        env.New(),
        flag.New(),
    )

    // Loading configuration
    err := loader.Load(cfg)
    if err != nil {
        panic("error to load config")
    }
}
