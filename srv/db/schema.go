package dbsrv

import (
	"database/sql"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"codeberg.org/reiver/go-log"
)

const codeCreate string = `
CREATE TABLE IF NOT EXISTS people (
	id INTEGER PRIMARY KEY,
	when_created INTEGER DEFAULT (strftime('%s','now')),

	name     TEXT NOT NULL,
	title    TEXT NOT NULL DEFAULT '',
	company  TEXT NOT NULL DEFAULT '',
	fedi_id  TEXT NOT NULL DEFAULT '',
	note     TEXT NOT NULL DEFAULT '',
	favorite BOOLEAN NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS groups (
	id INTEGER PRIMARY KEY,
	when_created INTEGER DEFAULT (strftime('%s','now')),

	name     TEXT NOT NULL,
	favorite BOOLEAN NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS group_members (
	id INTEGER PRIMARY KEY,

	group_id    INTEGER NOT NULL,
	person_name TEXT    NOT NULL,

	UNIQUE (group_id, person_name),

	FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages (
	id INTEGER PRIMARY KEY,
	when_created INTEGER DEFAULT (strftime('%s','now')),

	person_id INTEGER,
	group_id  INTEGER,

	from_me   BOOLEAN NOT NULL DEFAULT 0,
	sender    TEXT    NOT NULL DEFAULT '',
	text      TEXT    NOT NULL DEFAULT '',
	timestamp TEXT    NOT NULL DEFAULT '',

	FOREIGN KEY (person_id) REFERENCES people(id) ON DELETE CASCADE,
	FOREIGN KEY (group_id)  REFERENCES groups(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS gigs (
	id INTEGER PRIMARY KEY,
	when_created INTEGER DEFAULT (strftime('%s','now')),

	title       TEXT NOT NULL,
	company     TEXT NOT NULL DEFAULT '',
	location    TEXT NOT NULL DEFAULT '',
	type        TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',
	posted_by   TEXT NOT NULL DEFAULT '',
	timestamp   TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS profile (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	when_created INTEGER DEFAULT (strftime('%s','now')),

	name         TEXT NOT NULL DEFAULT '',
	title        TEXT NOT NULL DEFAULT '',
	company      TEXT NOT NULL DEFAULT '',
	fedi_id      TEXT NOT NULL DEFAULT '',
	summary_html TEXT NOT NULL DEFAULT '',
	icon_url     TEXT NOT NULL DEFAULT '',
	banner_url   TEXT NOT NULL DEFAULT '',
	profile_url  TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS resumes (
	id INTEGER PRIMARY KEY,
	when_created INTEGER DEFAULT (strftime('%s','now')),

	person_id  INTEGER,
	profile_id INTEGER,

	label   TEXT NOT NULL DEFAULT '',
	data    TEXT NOT NULL DEFAULT '',

	FOREIGN KEY (person_id)  REFERENCES people(id) ON DELETE CASCADE,
	FOREIGN KEY (profile_id) REFERENCES profile(id) ON DELETE CASCADE
);
`

func CreateDatabase(logger log.Logger, db *sql.DB) error {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		log.Alert(field.S("nil db"))
		var err error = erorr.Error("nil db")
		err = erorr.Wrap(err, "failed to create database table(s)")
		return err
	}

	_, err := db.Exec(codeCreate)
	if nil != err {
		const msg string = "failed to create database table(s)"
		log.Error(
			field.S(msg),
			field.E(err),
		)
		err = erorr.Wrap(err, msg)
		return err
	}

	return nil
}
