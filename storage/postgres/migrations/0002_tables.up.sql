CREATE TABLE vcs.authors (
    id SERIAL,
    name VARCHAR NOT NULL,
    email VARCHAR,
    UNIQUE (name, email),
    PRIMARY KEY ("id")
);


CREATE TABLE vcs.projects (
    id SERIAL,
    namespace  VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    UNIQUE (namespace, name),
    PRIMARY KEY ("id")
);
