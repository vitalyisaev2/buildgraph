package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"

	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/storage"
	"github.com/vitalyisaev2/buildgraph/storage/postgres/migrations"
)

var _ storage.Storage = (*storageImpl)(nil)

var (
	readOnlyTransaction = &sql.TxOptions{ReadOnly: true}
)

type storageImpl struct {
	db *sql.DB
}

type transaction func(*sql.Tx) error

func (s *storageImpl) performTransaction(
	ctx context.Context,
	method transaction,
	options *sql.TxOptions,
) error {
	// start tx
	tx, err := s.db.BeginTx(ctx, options)
	if err != nil {
		return err
	}

	// either commit, or rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// apply method
	err = method(tx)
	return err
}

func (s *storageImpl) SaveAuthor(ctx context.Context, author storage.Author) error {

	f := func(tx *sql.Tx) error {
		var id int

		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO vcs.authors(name, email) VALUES ($1, $2) ON CONFLICT DO NOTHING;`,
			author.GetName(), author.GetEmail())
		if err != nil {
			return err
		}

		err = tx.QueryRowContext(
			ctx,
			`SELECT id from vcs.authors WHERE name = $1 AND email = $2;`,
			author.GetName(), author.GetEmail()).Scan(&id)
		if err != nil {
			return err
		}

		author.SetID(id)
		return nil
	}

	return s.performTransaction(ctx, f, nil)
}

func (s *storageImpl) GetAuthor(ctx context.Context, name, email string) (storage.Author, error) {

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

func (s *storageImpl) Close() error { return s.Close() }

func NewStorage(cfg *config.PostgresConfig) (storage.Storage, error) {
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
		return nil, fmt.Errorf("Failed to prepare bindata driver: %v", err)
	}

	// Create new source with instance
	m, err := migrate.NewWithSourceInstance("go-bindata", driver, cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare migration: %v", err)
	}

	// Migrate database
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("Failed to migrate: %v", err)
	}

	// Prepare connection
	db, err := sql.Open("postgres", cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize DB client: %v", err)
	}

	s := &storageImpl{
		db: db,
	}
	return s, nil
}
