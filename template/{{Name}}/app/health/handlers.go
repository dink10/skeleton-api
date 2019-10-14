package health

import (
	"fmt"
	"net/http"
	"sort"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/gismart/{{Name}}/config"
	"bitbucket.org/gismart/{{Name}}/database"
	"bitbucket.org/gismart/{{Name}}/migration/postgres"
)

type healthResponse struct {
	ServiceName string  `json:"service-name"`
	Host        string  `json:"host"`
	Port        int     `json:"port"`
	Error       *string `json:"error,omitempty"`
}

// Create godoc
// @Description Liveness probe
// @Tags health
// @Success 200
// @Failure 503 {object} render.ErrorResponse
// @Router /health [get]
func Health(w http.ResponseWriter, _ *http.Request) {
	resp := gatherHealthData()
	for _, r := range resp {
		if r.Error != nil {
			log.Debugf("%s %s:%d is unavailable due to: %v", r.ServiceName, r.Host, r.Port, r.Error)
			w.WriteHeader(http.StatusServiceUnavailable)
			break
		}
	}

	regMigrations := postgres.GetRegisteredMigrations()
	sort.SliceStable(regMigrations, func(i, j int) bool {
		return regMigrations[i].Version < regMigrations[j].Version
	})

	if len(regMigrations) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	lastVersion, err := postgres.GetLastVersion()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lastRegisteredVersion := regMigrations[len(regMigrations)-1].Version

	if lastRegisteredVersion != lastVersion {
		log.Warn(fmt.Sprintf("Incompatible migration version: current version - %+v, registered version - %+v", lastVersion, lastRegisteredVersion))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func gatherHealthData() []healthResponse {
	resp := make([]healthResponse, 0)
	postgresResp := healthResponse{
		ServiceName: "postgres",
		Host:        config.Config.Postgres.Host,
		Port:        config.Config.Postgres.Port,
	}

	if err := database.PostgresPing(); err != nil {
		errString := err.Error()
		postgresResp.Error = &errString
	} else {
		postgresResp.Error = nil
	}

	resp = append(resp, postgresResp)

	return resp
}
