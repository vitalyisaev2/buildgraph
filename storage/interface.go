package storage

import (
	"context"

	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/vcs"
)

// Storage is an abstraction layer above the particular SQL/NoSQL storages;
// it should implement all the methods required by front and graph layer;
type Storage interface {
	// PushEvent
	SavePushEvent(context.Context, vcs.PushEvent) error
	common.Service
}
