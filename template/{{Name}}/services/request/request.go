package request

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"bitbucket.org/gismart/ddtracer"
	"bitbucket.org/gismart/{{Name}}/config"

	log "github.com/sirupsen/logrus"
)

func MakeRequest(ctx context.Context, u *url.URL, method string, body []byte) (int, []byte, error) {
	defer ddtracer.TraceExternalQuery(ctx, string(body), u.String(), method).Finish()

	client := http.Client{
		Timeout: time.Duration(config.Config.HTTPTimeout) * time.Second,
	}

	resp, err := client.Do(&http.Request{
		Method: method,
		URL:    u,
		Body:   ioutil.NopCloser(bytes.NewReader(body)),
		Header: http.Header{"Content-Type": {"application/json"}},
	})
	if err != nil {
		return 0, nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()

	return resp.StatusCode, respBody, nil
}
