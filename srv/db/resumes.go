package dbsrv

import (
	"database/sql"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"codeberg.org/reiver/go-log"
)

type ResumeRow struct {
	ID        int64
	PersonID  *int64
	ProfileID *int64
	Label     string
	Data      string // JSON-encoded resume data
}

func LoadResumesForPerson(logger log.Logger, db *sql.DB, personID int64) ([]ResumeRow, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return nil, erorr.Wrap(erorr.Error("nil db"), "failed to load resumes for person")
	}

	rows, err := db.Query(`SELECT id, person_id, profile_id, label, data FROM resumes WHERE person_id = ? ORDER BY id`, personID)
	if nil != err {
		log.Error(field.S("failed to load resumes for person"), field.E(err), field.Int64("person_id", personID))
		return nil, erorr.Wrap(err, "failed to load resumes for person")
	}
	defer rows.Close()

	var resumes []ResumeRow
	for rows.Next() {
		var r ResumeRow
		err := rows.Scan(&r.ID, &r.PersonID, &r.ProfileID, &r.Label, &r.Data)
		if nil != err {
			return nil, erorr.Wrap(err, "failed to scan resume row")
		}
		resumes = append(resumes, r)
	}
	if err := rows.Err(); nil != err {
		return nil, erorr.Wrap(err, "failed to iterate resume rows")
	}

	return resumes, nil
}

func LoadResumesForProfile(logger log.Logger, db *sql.DB) ([]ResumeRow, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return nil, erorr.Wrap(erorr.Error("nil db"), "failed to load resumes for profile")
	}

	rows, err := db.Query(`SELECT id, person_id, profile_id, label, data FROM resumes WHERE profile_id = 1 ORDER BY id`)
	if nil != err {
		log.Error(field.S("failed to load resumes for profile"), field.E(err))
		return nil, erorr.Wrap(err, "failed to load resumes for profile")
	}
	defer rows.Close()

	var resumes []ResumeRow
	for rows.Next() {
		var r ResumeRow
		err := rows.Scan(&r.ID, &r.PersonID, &r.ProfileID, &r.Label, &r.Data)
		if nil != err {
			return nil, erorr.Wrap(err, "failed to scan resume row")
		}
		resumes = append(resumes, r)
	}
	if err := rows.Err(); nil != err {
		return nil, erorr.Wrap(err, "failed to iterate resume rows")
	}

	return resumes, nil
}

func InsertResume(logger log.Logger, db *sql.DB, r ResumeRow) (int64, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return 0, erorr.Wrap(erorr.Error("nil db"), "failed to insert resume")
	}

	result, err := db.Exec(
		`INSERT INTO resumes (person_id, profile_id, label, data) VALUES (?, ?, ?, ?)`,
		r.PersonID, r.ProfileID, r.Label, r.Data,
	)
	if nil != err {
		log.Error(field.S("failed to insert resume"), field.E(err), field.String("label", r.Label))
		return 0, erorr.Wrap(err, "failed to insert resume")
	}

	id, err := result.LastInsertId()
	if nil != err {
		return 0, erorr.Wrap(err, "failed to get last insert id for resume")
	}

	return id, nil
}
