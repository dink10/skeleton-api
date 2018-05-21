package main

import (
    "bitbucket.org/gismart/{{Name}}/config"
    "bitbucket.org/gismart/{{Name}}/logger"
    "bitbucket.org/gismart/{{Name}}/server"
)

var cfg = config.Config
var log = logger.Logger

func init() {
    log.Infof("application started on %v:%v", cfg.Server.Host, cfg.Server.Port)
}

func main() {
    server.Run()

}
