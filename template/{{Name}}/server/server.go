package server

import (
    "fmt"
    "net/http"
    cfg "bitbucket.org/gismart/{{Name}}/config"
    log "github.com/sirupsen/logrus"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Run() {
    host := cfg.Config.Server.Host
    port := cfg.Config.Server.Port
    address := fmt.Sprintf("%v:%v", host, port)
    router := runRoute()

    if cfg.Config.Logger.DataDogAgentAddr != "" &&
        cfg.Config.Logger.DataDogEnv != "" {
        tracer.Start(tracer.WithAgentAddr(cfg.Config.Logger.DataDogAgentAddr))
        defer tracer.Stop()
    }

    log.Fatal(http.ListenAndServe(address, router))
}
