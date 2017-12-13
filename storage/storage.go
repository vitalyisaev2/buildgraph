package storage

import (
	"context"

	"github.com/vitalyisaev2/buildgraph/common"
)

// Storage is an abstraction layer above the particular SQL/NoSQL storages;
// it should implement all the methods required by front and graph layer;
type Storage interface {

	// Author
	SaveAuthor(ctx context.Context, author Author) error
	GetAuthor(ctx context.Context, name, email string) (Author, error)

	// Project
	//SaveProject(ctx context.Context, project Project) error
	//GetProject(ctx context.Context, namespace, name string) (Project, error)

	common.Service
}
