CREATE TABLE vcs.projects (
    id SERIAL PRIMARY KEY,
    namespace TEXT NOT NULL,
    name TEXT NOT NULL,
    http_url TEXT NOT NULL,
    UNIQUE (namespace, name)
);

CREATE TABLE vcs.events (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES vcs.projects(id)
);


CREATE TABLE vcs.authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT,
    UNIQUE (name, email)
);

CREATE TABLE vcs.commits (
    id SERIAL PRIMARY KEY,
    hash TEXT NOT NULL,
    message TEXT NOT NULL,
    time TIMESTAMP NOT NULL,
    url TEXT NOT NULL,
    added TEXT[],
    modified TEXT[],
    removed TEXT[],
    project_id INTEGER REFERENCES vcs.projects(id),
    author_id INTEGER REFERENCES vcs.authors(id),
    event_id INTEGER REFERENCES vcs.events(id),
    UNIQUE(hash, project_id)
);
/*
*/
