package main

import (
    "bitbucket.org/gismart/{{Name}}/config"
    "bitbucket.org/gismart/{{Name}}/server"
    "bitbucket.org/gismart/{{Name}}/logger"
)


func main() {
    config.InitConfig()
    logger.InitLogger(config.Config.Logger.LogLevel, config.Config.Logger.SentryDSN)

    server.Run(config.Config.Server.Host, config.Config.Server.Port)
}

