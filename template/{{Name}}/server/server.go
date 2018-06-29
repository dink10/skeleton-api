package server

import (
    "fmt"
    "net/http"
    cfg "bitbucket.org/gismart/{{Name}}/config"
    log "github.com/sirupsen/logrus"
)

func Run(datadogWebServiceName string) {
    host := cfg.Config.Server.Host
    port := cfg.Config.Server.Port
    address := fmt.Sprintf("%v:%v", host, port)
    router := runRoute(datadogWebServiceName)

    log.Fatal(http.ListenAndServe(address, router))
}
