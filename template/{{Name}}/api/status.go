package api

import (
    "fmt"
    "net/http"
)

func Status(resp http.ResponseWriter, req *http.Request) {
    fmt.Fprint(resp, "OK")
}
