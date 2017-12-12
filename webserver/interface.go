package webserver

import (
	"net/http"

	"github.com/vitalyisaev2/buildgraph/common"
)

type Webserver interface {
	common.Service
	GitlabPushEvent(http.ResponseWriter, *http.Request)
}
