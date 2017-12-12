package common

import (
	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/storage"
)

type Service interface {
	Start()
	Stop()
}

type Collection struct {
	Logger    common.Logger
	WebServer common.Service
	Storage   storage.Storage
}
