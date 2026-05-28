package main

import (
	guipkg "protogo/gui"
	dbsrv "protogo/srv/db"
	"protogo/srv/log"
)

func seeddb() {
	logger := logsrv.Logger()
	db := dbsrv.WriteDB()

	err := dbsrv.SeedIfEmpty(logger, db)
	if nil != err {
		return
	}

	// Check if data already exists.
	var count int64
	err2 := db.QueryRow(`SELECT COUNT(*) FROM people`).Scan(&count)
	if nil != err2 {
		return
	}
	if 0 < count {
		return
	}

	// Seed people.
	dummyPeople := guipkg.DummyPeople()
	var seedPeople []dbsrv.SeedPerson
	for _, p := range dummyPeople {
		sp := dbsrv.SeedPerson{
			Name:     p.Name,
			Title:    p.Title,
			Company:  p.Company,
			FediID:   p.FediID,
			Note:     p.Note,
			Favorite: p.Favorite,
		}
		for _, m := range p.Messages {
			sp.Messages = append(sp.Messages, dbsrv.SeedMessage{
				FromMe:    m.FromMe,
				Sender:    m.Sender,
				Text:      m.Text,
				Timestamp: m.Timestamp,
			})
		}
		for _, r := range p.Resumes {
			sp.Resumes = append(sp.Resumes, dbsrv.SeedResume{
				Label: r.Label,
				Data:  resumeToMap(r),
			})
		}
		seedPeople = append(seedPeople, sp)
	}
	dbsrv.SeedPeople(logger, db, seedPeople)

	// Seed groups.
	dummyGroups := guipkg.DummyGroups()
	var seedGroups []dbsrv.SeedGroup
	for _, g := range dummyGroups {
		sg := dbsrv.SeedGroup{
			Name:     g.Name,
			Members:  g.Members,
			Favorite: g.Favorite,
		}
		for _, m := range g.Messages {
			sg.Messages = append(sg.Messages, dbsrv.SeedMessage{
				FromMe:    m.FromMe,
				Sender:    m.Sender,
				Text:      m.Text,
				Timestamp: m.Timestamp,
			})
		}
		seedGroups = append(seedGroups, sg)
	}
	dbsrv.SeedGroups(logger, db, seedGroups)

	// Seed gigs.
	dummyGigs := guipkg.DummyGigs()
	var seedGigs []dbsrv.SeedGig
	for _, g := range dummyGigs {
		seedGigs = append(seedGigs, dbsrv.SeedGig{
			Title:       g.Title,
			Company:     g.Company,
			Location:    g.Location,
			Type:        g.Type,
			Description: g.Description,
			PostedBy:    g.PostedBy,
			Timestamp:   g.Timestamp,
		})
	}
	dbsrv.SeedGigs(logger, db, seedGigs)

	// Seed profile.
	dummyMe := guipkg.DummyMe()
	sp := dbsrv.SeedProfile{
		Name:    dummyMe.Name,
		Title:   dummyMe.Title,
		Company: dummyMe.Company,
		FediID:  dummyMe.FediID,
	}
	for _, r := range dummyMe.Resumes {
		sp.Resumes = append(sp.Resumes, dbsrv.SeedResume{
			Label: r.Label,
			Data:  resumeToMap(r),
		})
	}
	dbsrv.SeedProfileData(logger, db, sp)
}

func resumeToMap(r guipkg.Resume) map[string]interface{} {
	return map[string]interface{}{
		"Basics":       r.Basics,
		"Work":         r.Work,
		"Volunteer":    r.Volunteer,
		"Education":    r.Education,
		"Awards":       r.Awards,
		"Certificates": r.Certificates,
		"Publications": r.Publications,
		"Skills":       r.Skills,
		"Languages":    r.Languages,
		"Interests":    r.Interests,
		"References":   r.References,
		"Projects":     r.Projects,
	}
}
