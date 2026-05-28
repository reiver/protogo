package dbsrv

import (
	"database/sql"
	"encoding/json"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"codeberg.org/reiver/go-log"
)

type SeedPerson struct {
	Name     string
	Title    string
	Company  string
	FediID   string
	Note     string
	Favorite bool
	Messages []SeedMessage
	Resumes  []SeedResume
}

type SeedMessage struct {
	FromMe    bool
	Sender    string
	Text      string
	Timestamp string
}

type SeedResume struct {
	Label string
	Data  interface{} // will be JSON-encoded
}

type SeedGroup struct {
	Name     string
	Members  []string
	Favorite bool
	Messages []SeedMessage
}

type SeedGig struct {
	Title       string
	Company     string
	Location    string
	Type        string
	Description string
	PostedBy    string
	Timestamp   string
}

type SeedProfile struct {
	Name    string
	Title   string
	Company string
	FediID  string
	Resumes []SeedResume
}

func SeedIfEmpty(logger log.Logger, db *sql.DB) error {
	log := logger.Begin()
	defer log.End()

	if nil == db {
		return erorr.Wrap(erorr.Error("nil db"), "failed to seed database")
	}

	var count int64
	err := db.QueryRow(`SELECT COUNT(*) FROM people`).Scan(&count)
	if nil != err {
		log.Error(field.S("failed to count people"), field.E(err))
		return erorr.Wrap(err, "failed to count people for seed check")
	}

	if 0 < count {
		log.Trace(field.S("database already has data, skipping seed"))
		return nil
	}

	log.Inform(field.S("seeding database with initial data"))

	return nil
}

func SeedPeople(logger log.Logger, db *sql.DB, people []SeedPerson) error {
	log := logger.Begin()
	defer log.End()

	for _, sp := range people {
		personID, err := InsertPerson(logger, db, PersonRow{
			Name:     sp.Name,
			Title:    sp.Title,
			Company:  sp.Company,
			FediID:   sp.FediID,
			Note:     sp.Note,
			Favorite: sp.Favorite,
		})
		if nil != err {
			return err
		}

		for _, sm := range sp.Messages {
			pid := personID
			_, err := InsertMessage(logger, db, MessageRow{
				PersonID:  &pid,
				FromMe:    sm.FromMe,
				Sender:    sm.Sender,
				Text:      sm.Text,
				Timestamp: sm.Timestamp,
			})
			if nil != err {
				return err
			}
		}

		for _, sr := range sp.Resumes {
			dataBytes, err := json.Marshal(sr.Data)
			if nil != err {
				log.Error(field.S("failed to marshal resume data"), field.E(err))
				return erorr.Wrap(err, "failed to marshal resume data")
			}

			pid := personID
			_, err = InsertResume(logger, db, ResumeRow{
				PersonID: &pid,
				Label:    sr.Label,
				Data:     string(dataBytes),
			})
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func SeedGroups(logger log.Logger, db *sql.DB, groups []SeedGroup) error {
	log := logger.Begin()
	defer log.End()

	for _, sg := range groups {
		groupID, err := InsertGroup(logger, db, GroupRow{
			Name:     sg.Name,
			Favorite: sg.Favorite,
			Members:  sg.Members,
		})
		if nil != err {
			return err
		}

		for _, sm := range sg.Messages {
			gid := groupID
			_, err := InsertMessage(logger, db, MessageRow{
				GroupID:   &gid,
				FromMe:    sm.FromMe,
				Sender:    sm.Sender,
				Text:      sm.Text,
				Timestamp: sm.Timestamp,
			})
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func SeedGigs(logger log.Logger, db *sql.DB, gigs []SeedGig) error {
	for _, sg := range gigs {
		_, err := InsertGig(logger, db, GigRow{
			Title:       sg.Title,
			Company:     sg.Company,
			Location:    sg.Location,
			Type:        sg.Type,
			Description: sg.Description,
			PostedBy:    sg.PostedBy,
			Timestamp:   sg.Timestamp,
		})
		if nil != err {
			return err
		}
	}

	return nil
}

func SeedProfileData(logger log.Logger, db *sql.DB, profile SeedProfile) error {
	log := logger.Begin()
	defer log.End()

	err := UpsertProfile(logger, db, ProfileRow{
		Name:    profile.Name,
		Title:   profile.Title,
		Company: profile.Company,
		FediID:  profile.FediID,
	})
	if nil != err {
		return err
	}

	for _, sr := range profile.Resumes {
		dataBytes, err := json.Marshal(sr.Data)
		if nil != err {
			log.Error(field.S("failed to marshal resume data"), field.E(err))
			return erorr.Wrap(err, "failed to marshal resume data")
		}

		var profileID int64 = 1
		_, err = InsertResume(logger, db, ResumeRow{
			ProfileID: &profileID,
			Label:     sr.Label,
			Data:      string(dataBytes),
		})
		if nil != err {
			return err
		}
	}

	return nil
}
