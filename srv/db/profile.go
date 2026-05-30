package dbsrv

import (
	"database/sql"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"codeberg.org/reiver/go-log"
)

type ProfileRow struct {
	Name        string
	Title       string
	Company     string
	FediID      string
	SummaryHTML string
	IconURL     string
	BannerURL   string
	ProfileURL  string
}

func LoadProfile(logger log.Logger, db *sql.DB) (ProfileRow, bool, error) {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return ProfileRow{}, false, erorr.Wrap(erorr.Error("nil db"), "failed to load profile")
	}

	var p ProfileRow
	err := db.QueryRow(`SELECT name, title, company, fedi_id, summary_html, icon_url, banner_url, profile_url FROM profile WHERE id = 1`).Scan(&p.Name, &p.Title, &p.Company, &p.FediID, &p.SummaryHTML, &p.IconURL, &p.BannerURL, &p.ProfileURL)
	if err == sql.ErrNoRows {
		return ProfileRow{}, false, nil
	}
	if nil != err {
		log.Error(field.S("failed to load profile"), field.E(err))
		return ProfileRow{}, false, erorr.Wrap(err, "failed to load profile")
	}

	return p, true, nil
}

func UpdateProfileFediID(logger log.Logger, db *sql.DB, fediID string) error {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return erorr.Wrap(erorr.Error("nil db"), "failed to update profile fedi_id")
	}

	_, err := db.Exec(`UPDATE profile SET fedi_id = ? WHERE id = 1`, fediID)
	if nil != err {
		log.Error(field.S("failed to update profile fedi_id"), field.E(err), field.String("fedi_id", fediID))
		return erorr.Wrap(err, "failed to update profile fedi_id")
	}

	return nil
}

func UpdateProfileFromActor(logger log.Logger, db *sql.DB, name string, summaryHTML string, iconURL string, bannerURL string, profileURL string) error {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return erorr.Wrap(erorr.Error("nil db"), "failed to update profile from actor")
	}

	_, err := db.Exec(
		`UPDATE profile SET name = ?, summary_html = ?, icon_url = ?, banner_url = ?, profile_url = ? WHERE id = 1`,
		name, summaryHTML, iconURL, bannerURL, profileURL,
	)
	if nil != err {
		log.Error(field.S("failed to update profile from actor"), field.E(err))
		return erorr.Wrap(err, "failed to update profile from actor")
	}

	return nil
}

func UpsertProfile(logger log.Logger, db *sql.DB, p ProfileRow) error {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return erorr.Wrap(erorr.Error("nil db"), "failed to upsert profile")
	}

	_, err := db.Exec(
		`INSERT INTO profile (id, name, title, company, fedi_id, summary_html, icon_url, banner_url, profile_url) VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET name = ?, title = ?, company = ?, fedi_id = ?, summary_html = ?, icon_url = ?, banner_url = ?, profile_url = ?`,
		p.Name, p.Title, p.Company, p.FediID, p.SummaryHTML, p.IconURL, p.BannerURL, p.ProfileURL,
		p.Name, p.Title, p.Company, p.FediID, p.SummaryHTML, p.IconURL, p.BannerURL, p.ProfileURL,
	)
	if nil != err {
		log.Error(field.S("failed to upsert profile"), field.E(err))
		return erorr.Wrap(err, "failed to upsert profile")
	}

	return nil
}
