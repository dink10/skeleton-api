package server

import (
    "fmt"
    "net/http"
    cfg "bitbucket.org/gismart/{{Name}}/config"
    log "bitbucket.org/gismart/{{Name}}/logger"
)

func Run() {
    host := cfg.Config.Server.Host
    port := cfg.Config.Server.Port
    address := fmt.Sprintf("%v:%v", host, port)
    router := runRoute()

    log.Logger.Fatal(http.ListenAndServe(address, router))
}
