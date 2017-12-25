package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/vcs"
)

type step func() error

type executor struct {
	tx     *sql.Tx
	logger *logrus.Logger
	steps  []step
}

func (ex *executor) saveProject(project vcs.Project) {
	f := func() error {
		_, err := ex.tx.Exec(
			`INSERT INTO vcs.projects(name, namespace, http_url)
			VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			project.GetName(), project.GetNamespace(), project.GetHTTPURL(),
		)
		if err != nil {
			return err
		}

		var id common.ObjectID
		err = ex.tx.QueryRow(
			`SELECT id from vcs.projects WHERE name = $1 AND namespace = $2;`,
			project.GetName(), project.GetNamespace()).Scan(&id)
		if err != nil {
			return err
		}

		project.SetObjectID(id)

		return nil
	}

	ex.addStep(f)
}

func (ex *executor) saveEvent(event vcs.PushEvent) {
	f := func() error {
		var id common.ObjectID
		err := ex.tx.QueryRow(
			`INSERT INTO vcs.events(project_id) VALUES ($1) RETURNING ID`,
			event.GetProject().GetObjectID(),
		).Scan(&id)

		if err != nil {
			return err
		}

		event.SetObjectID(id)
		return nil
	}

	ex.addStep(f)
}

func (ex *executor) saveCommit(
	commit vcs.Commit,
	event vcs.PushEvent) {

	f := func() error {
		var id common.ObjectID
		err := ex.tx.QueryRow(
			`INSERT INTO vcs.commits(
				hash, message, time, url,
				added, modified, removed,
				project_id, author_id, event_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			commit.GetHash(),
			commit.GetMessage(),
			pq.FormatTimestamp(commit.GetTimestamp()),
			commit.GetURL(),
			pq.Array(commit.GetAdded()),
			pq.Array(commit.GetModified()),
			pq.Array(commit.GetRemoved()),
			event.GetProject().GetObjectID(),
			commit.GetAuthor().GetObjectID(),
			event.GetObjectID(),
		).Scan(&id)

		if err != nil {
			return err
		}

		commit.SetObjectID(id)
		return nil
	}

	ex.addStep(f)
}

func (ex *executor) saveAuthor(author vcs.Author) {

	f := func() error {
		var id int

		_, err := ex.tx.Exec(
			`INSERT INTO vcs.authors(name, email) VALUES ($1, $2) ON CONFLICT DO NOTHING;`,
			author.GetName(), author.GetEmail())
		if err != nil {
			return err
		}

		err = ex.tx.QueryRow(
			`SELECT id from vcs.authors WHERE name = $1 AND email = $2;`,
			author.GetName(), author.GetEmail()).Scan(&id)
		if err != nil {
			return err
		}

		author.SetObjectID(id)
		return nil
	}

	ex.addStep(f)
}

func (ex *executor) addStep(f step) {
	ex.steps = append(ex.steps, f)
}

func (ex *executor) finalize() error {
	var err error

	// either commit, or rollback on exit
	defer func() {
		if err != nil {
			err = ex.tx.Rollback()
		} else {
			err = ex.tx.Commit()
		}
	}()

	// walks through stored steps and
	for i, s := range ex.steps {
		if err = s(); err != nil {
			ex.logger.WithError(err).WithField("step", i).Error("transaction error")
			break
		}
	}

	return err
}
