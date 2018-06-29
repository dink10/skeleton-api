package main

import (
    "bitbucket.org/gismart/{{Name}}/config"
    "bitbucket.org/gismart/{{Name}}/server"
    log "github.com/sirupsen/logrus"
    _ "bitbucket.org/gismart/{{Name}}/logger"
    "os"
    "fmt"
)

var cfg = config.Config

func init() {
    log.Infof("application started on %v:%v", cfg.Server.Host, cfg.Server.Port)
}

func main() {

    datadogEnv := os.Getenv("DATADOG_ENV")
    if datadogEnv == "" {
        datadogEnv = "stage"
    }

    datadogWebServiceName := fmt.Sprintf("{{projectName}}web-%s", datadogEnv)

    server.Run(datadogWebServiceName)

}
