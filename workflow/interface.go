package workflow

import "github.com/vitalyisaev2/buildgraph/vcs"

type Manager interface {
	RegisterVCSPushEvent(vcs.PushEvent) error
}
