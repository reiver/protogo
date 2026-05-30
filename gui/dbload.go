package gui

import (
	"encoding/json"

	dbsrv "protogo/srv/db"
	"protogo/srv/log"
)

func loadPeopleFromDB() []Person {
	logger := logsrv.Logger()
	db := dbsrv.ReadDB()

	personRows, err := dbsrv.LoadPeople(logger, db)
	if nil != err {
		return DummyPeople()
	}

	var people []Person
	for _, pr := range personRows {
		p := Person{
			DBID:     pr.ID,
			Name:     pr.Name,
			Title:    pr.Title,
			Company:  pr.Company,
			FediID:   pr.FediID,
			Note:     pr.Note,
			Favorite: pr.Favorite,
		}

		msgRows, err := dbsrv.LoadMessagesForPerson(logger, db, pr.ID)
		if nil == err {
			for _, mr := range msgRows {
				p.Messages = append(p.Messages, ChatMessage{
					FromMe:    mr.FromMe,
					Sender:    mr.Sender,
					Text:      mr.Text,
					Timestamp: mr.Timestamp,
				})
			}
		}

		resumeRows, err := dbsrv.LoadResumesForPerson(logger, db, pr.ID)
		if nil == err {
			for _, rr := range resumeRows {
				r := decodeResume(rr)
				p.Resumes = append(p.Resumes, r)
			}
		}

		people = append(people, p)
	}

	return people
}

func loadGroupsFromDB() []Group {
	logger := logsrv.Logger()
	db := dbsrv.ReadDB()

	groupRows, err := dbsrv.LoadGroups(logger, db)
	if nil != err {
		return DummyGroups()
	}

	var groups []Group
	for _, gr := range groupRows {
		g := Group{
			DBID:     gr.ID,
			Name:     gr.Name,
			Members:  gr.Members,
			Favorite: gr.Favorite,
		}

		msgRows, err := dbsrv.LoadMessagesForGroup(logger, db, gr.ID)
		if nil == err {
			for _, mr := range msgRows {
				g.Messages = append(g.Messages, ChatMessage{
					FromMe:    mr.FromMe,
					Sender:    mr.Sender,
					Text:      mr.Text,
					Timestamp: mr.Timestamp,
				})
			}
		}

		groups = append(groups, g)
	}

	return groups
}

func loadGigsFromDB() []Gig {
	logger := logsrv.Logger()
	db := dbsrv.ReadDB()

	gigRows, err := dbsrv.LoadGigs(logger, db)
	if nil != err {
		return DummyGigs()
	}

	var gigs []Gig
	for _, gr := range gigRows {
		gigs = append(gigs, Gig{
			Title:       gr.Title,
			Company:     gr.Company,
			Location:    gr.Location,
			Type:        gr.Type,
			Description: gr.Description,
			PostedBy:    gr.PostedBy,
			Timestamp:   gr.Timestamp,
		})
	}

	return gigs
}

func loadMeFromDB() Person {
	logger := logsrv.Logger()
	db := dbsrv.ReadDB()

	profileRow, found, err := dbsrv.LoadProfile(logger, db)
	if nil != err || !found {
		return DummyMe()
	}

	me := Person{
		Name:        profileRow.Name,
		Title:       profileRow.Title,
		Company:     profileRow.Company,
		FediID:      profileRow.FediID,
		SummaryHTML: profileRow.SummaryHTML,
		IconURL:     profileRow.IconURL,
		BannerURL:   profileRow.BannerURL,
		ProfileURL:  profileRow.ProfileURL,
	}

	resumeRows, err := dbsrv.LoadResumesForProfile(logger, db)
	if nil == err {
		for _, rr := range resumeRows {
			r := decodeResume(rr)
			me.Resumes = append(me.Resumes, r)
		}
	}

	return me
}

type resumeData struct {
	Basics       ResumeBasics        `json:"Basics"`
	Work         []ResumeWork        `json:"Work"`
	Volunteer    []ResumeVolunteer   `json:"Volunteer"`
	Education    []ResumeEducation   `json:"Education"`
	Awards       []ResumeAward       `json:"Awards"`
	Certificates []ResumeCertificate `json:"Certificates"`
	Publications []ResumePublication `json:"Publications"`
	Skills       []ResumeSkill       `json:"Skills"`
	Languages    []ResumeLanguage    `json:"Languages"`
	Interests    []ResumeInterest    `json:"Interests"`
	References   []ResumeReference   `json:"References"`
	Projects     []ResumeProject     `json:"Projects"`
}

func decodeResume(rr dbsrv.ResumeRow) Resume {
	r := Resume{
		Label: rr.Label,
	}

	if "" != rr.Data {
		var rd resumeData
		if nil == json.Unmarshal([]byte(rr.Data), &rd) {
			r.Basics = rd.Basics
			r.Work = rd.Work
			r.Volunteer = rd.Volunteer
			r.Education = rd.Education
			r.Awards = rd.Awards
			r.Certificates = rd.Certificates
			r.Publications = rd.Publications
			r.Skills = rd.Skills
			r.Languages = rd.Languages
			r.Interests = rd.Interests
			r.References = rd.References
			r.Projects = rd.Projects
		}
	}

	return r
}
