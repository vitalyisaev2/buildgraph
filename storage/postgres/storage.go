package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"

	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/storage"
	"github.com/vitalyisaev2/buildgraph/storage/postgres/migrations"
	"github.com/vitalyisaev2/buildgraph/vcs"
)

//var _ storage.Storage = (*defaultStorage)(nil)

var (
	readOnlyTransaction = &sql.TxOptions{ReadOnly: true}
)

type defaultStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

func (s *defaultStorage) SavePushEvent(ctx context.Context, event vcs.PushEvent) error {
	ex, err := s.makeExecutor(ctx, nil)
	if err != nil {
		return err
	}

	ex.saveProject(event.GetProject())
	ex.saveEvent(event)
	for _, commit := range event.GetCommits() {
		ex.saveAuthor(commit.GetAuthor())
		ex.saveCommit(commit, event)
	}
	return ex.finalize()
}

/*
func (s *defaultStorage) GetAuthor(ctx context.Context, name, email string) (storage.Author, error) {

	var result *author

	f := func(tx *sql.Tx) error {
		var id int
		err := tx.QueryRowContext(
			ctx,
			`SELECT id FROM vcs.authors WHERE name = $1 AND email = $2;`,
			name, email).Scan(&id)
		if err != nil {
			return err
		}
		result = &author{
			model: model{id: id},
			name:  name,
			email: email,
		}
		return nil
	}

	if err := s.performTransaction(ctx, f, readOnlyTransaction); err != nil {
		return nil, err
	}
	return result, nil
}
*/

func (s *defaultStorage) Stop() {
	if err := s.db.Close(); err != nil {
		s.logger.WithError(err).Error("Failed to close database")
	}
}

// performTransaction wraps routine that makes new transaction,
// executes it, than commits or rolls back
func (s *defaultStorage) makeExecutor(
	ctx context.Context,
	options *sql.TxOptions,
) (*executor, error) {

	// start tx
	tx, err := s.db.BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}

	// wrap transaction into executor
	ex := &executor{
		tx:     tx,
		logger: s.logger,
	}

	return ex, nil
}

func NewStorage(logger *logrus.Logger, cfg *config.PostgresConfig) (storage.Storage, error) {

	// Prepare migrations
	resource := bindata.Resource(
		migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		},
	)

	// Create migration driver
	driver, err := bindata.WithInstance(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare bindata driver: %v", err)
	}

	// Create new source with instance
	m, err := migrate.NewWithSourceInstance("go-bindata", driver, cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to prepare migration: %v", err)
	}

	// Migrate database
	logger.Debug("trying to migrate database")
	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return nil, fmt.Errorf("failed to migrate: %v", err)
		}
		logger.Debug("database is up-to-date")
	}

	// Prepare connection
	db, err := sql.Open("postgres", cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DB client: %v", err)
	}

	s := &defaultStorage{
		db:     db,
		logger: logger,
	}
	return s, nil
}
