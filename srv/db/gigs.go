package dbsrv

import (
	"database/sql"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"codeberg.org/reiver/go-log"
)

type GigRow struct {
	ID          int64
	Title       string
	Company     string
	Location    string
	Type        string
	Description string
	PostedBy    string
	Timestamp   string
}

func LoadGigs(logger log.Logger, db *sql.DB) ([]GigRow, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return nil, erorr.Wrap(erorr.Error("nil db"), "failed to load gigs")
	}

	rows, err := db.Query(`SELECT id, title, company, location, type, description, posted_by, timestamp FROM gigs ORDER BY id`)
	if nil != err {
		log.Error(field.S("failed to load gigs"), field.E(err))
		return nil, erorr.Wrap(err, "failed to load gigs")
	}
	defer rows.Close()

	var gigs []GigRow
	for rows.Next() {
		var g GigRow
		err := rows.Scan(&g.ID, &g.Title, &g.Company, &g.Location, &g.Type, &g.Description, &g.PostedBy, &g.Timestamp)
		if nil != err {
			log.Error(field.S("failed to scan gig row"), field.E(err))
			return nil, erorr.Wrap(err, "failed to scan gig row")
		}
		gigs = append(gigs, g)
	}
	if err := rows.Err(); nil != err {
		return nil, erorr.Wrap(err, "failed to iterate gig rows")
	}

	return gigs, nil
}

func InsertGig(logger log.Logger, db *sql.DB, g GigRow) (int64, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return 0, erorr.Wrap(erorr.Error("nil db"), "failed to insert gig")
	}

	result, err := db.Exec(
		`INSERT INTO gigs (title, company, location, type, description, posted_by, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		g.Title, g.Company, g.Location, g.Type, g.Description, g.PostedBy, g.Timestamp,
	)
	if nil != err {
		log.Error(field.S("failed to insert gig"), field.E(err), field.String("title", g.Title))
		return 0, erorr.Wrap(err, "failed to insert gig")
	}

	id, err := result.LastInsertId()
	if nil != err {
		return 0, erorr.Wrap(err, "failed to get last insert id for gig")
	}

	return id, nil
}
