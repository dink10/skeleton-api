package app

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"bitbucket.org/gismart/{{Name}}/config"
)

type healhResponse struct {
	ServiceName string  `json:"service-name"`
	Host        string  `json:"host"`
	Port        int     `json:"port"`
	Error       *string `json:"error,omitempty"`
}

var cfg = config.Config

// Health returns brief status based on HealthCheck results
func Health(w http.ResponseWriter, r *http.Request) {
	resp := gatherHealthData()
	for _, r := range resp {
		if r.Error != nil {
			log.Debugf("%s %s:%s is unavailable due to: %s", r.ServiceName, r.Host, r.Port, r.Error)
			w.WriteHeader(http.StatusServiceUnavailable)
			break
		}
	}

	w.WriteHeader(http.StatusOK)
}

func gatherHealthData() []healhResponse {
	// check all third party services availability here
	resp := make([]healhResponse, 0)

	return resp
}
