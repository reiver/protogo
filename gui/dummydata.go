package gui

type Person struct {
	Name     string
	Title    string
	Company  string
	FediID   string
	Note     string
	Favorite bool
	Resumes  []Resume
	Messages []ChatMessage
}

type ChatMessage struct {
	FromMe    bool
	Sender    string
	Text      string
	Timestamp string
}

type Resume struct {
	Label string

	Basics       ResumeBasics
	Work         []ResumeWork
	Volunteer    []ResumeVolunteer
	Education    []ResumeEducation
	Awards       []ResumeAward
	Certificates []ResumeCertificate
	Publications []ResumePublication
	Skills       []ResumeSkill
	Languages    []ResumeLanguage
	Interests    []ResumeInterest
	References   []ResumeReference
	Projects     []ResumeProject
}

type ResumeBasics struct {
	Name     string
	Label    string
	Image    string
	Email    string
	Phone    string
	URL      string
	Summary  string
	Location ResumeLocation
	Profiles []ResumeProfile
}

type ResumeLocation struct {
	Address     string
	PostalCode  string
	City        string
	CountryCode string
	Region      string
}

type ResumeProfile struct {
	Network  string
	Username string
	URL      string
}

type ResumeWork struct {
	Name        string
	Location    string
	Description string
	Position    string
	URL         string
	StartDate   string
	EndDate     string
	Summary     string
	Highlights  []string
}

type ResumeVolunteer struct {
	Organization string
	Position     string
	URL          string
	StartDate    string
	EndDate      string
	Summary      string
	Highlights   []string
}

type ResumeEducation struct {
	Institution string
	URL         string
	Area        string
	StudyType   string
	StartDate   string
	EndDate     string
	Score       string
	Courses     []string
}

type ResumeAward struct {
	Title   string
	Date    string
	Awarder string
	Summary string
}

type ResumeCertificate struct {
	Name   string
	Date   string
	URL    string
	Issuer string
}

type ResumePublication struct {
	Name        string
	Publisher   string
	ReleaseDate string
	URL         string
	Summary     string
}

type ResumeSkill struct {
	Name     string
	Level    string
	Keywords []string
}

type ResumeLanguage struct {
	Language string
	Fluency  string
}

type ResumeInterest struct {
	Name     string
	Keywords []string
}

type ResumeReference struct {
	Name      string
	Reference string
}

type ResumeProject struct {
	Name        string
	Description string
	Highlights  []string
	Keywords    []string
	Roles       []string
	StartDate   string
	EndDate     string
	URL         string
	Entity      string
	Type        string
}

type Group struct {
	Name     string
	Members  []string
	Favorite bool
	Messages []ChatMessage
}

func dummyMe() Person {
	return Person{
		Name:    "Charles Iliya Krempeaux",
		Title:   "Software Engineer",
		Company: "ProToGo",
		FediID:  "@reiver@mastodon.social",
		Resumes: []Resume{
			{
				Label: "My Resume",
				Basics: ResumeBasics{
					Name:    "Charles Iliya Krempeaux",
					Label:   "Software Engineer",
					Email:   "charles@example.com",
					URL:     "https://changelog.ca",
					Summary: "Software engineer passionate about the Fediverse, decentralized systems, and building tools for the Social Web.",
					Location: ResumeLocation{
						City:        "Vancouver",
						CountryCode: "CA",
						Region:      "British Columbia",
					},
					Profiles: []ResumeProfile{
						{Network: "Fediverse", Username: "@reiver@mastodon.social", URL: "https://mastodon.social/@reiver"},
						{Network: "Codeberg", Username: "reiver", URL: "https://codeberg.org/reiver"},
						{Network: "GitHub", Username: "reiver", URL: "https://github.com/reiver"},
					},
				},
				Skills: []ResumeSkill{
					{Name: "Go", Level: "Expert", Keywords: []string{"concurrency", "networking", "CLI tools"}},
					{Name: "Fediverse", Level: "Expert", Keywords: []string{"ActivityPub", "NodeInfo", "WebFinger", "federation"}},
				},
			},
		},
	}
}

func dummyPeople() []Person {
	return []Person{
		{
			Name:     "Alice Zhang",
			Title:    "Senior Backend Engineer",
			Company:  "Meshwork Labs",
			FediID:   "@alice@meshwork.social",
			Favorite: true,
			Note:     "Met at GopherCon 2025. Interested in distributed systems and ActivityPub.",
			Resumes: []Resume{
				{
					Label: "Software Engineer Resume",
					Basics: ResumeBasics{
						Name:    "Alice Zhang",
						Label:   "Senior Backend Engineer",
						Email:   "alice@example.com",
						Phone:   "+1-604-555-0101",
						URL:     "https://alicezhang.dev",
						Summary: "Backend engineer with 10 years of experience building distributed systems, APIs, and federation protocols in Go.",
						Location: ResumeLocation{
							City:        "Vancouver",
							CountryCode: "CA",
							Region:      "British Columbia",
						},
						Profiles: []ResumeProfile{
							{Network: "Fediverse", Username: "@alice@meshwork.social", URL: "https://meshwork.social/@alice"},
							{Network: "GitHub", Username: "alicezhang", URL: "https://github.com/alicezhang"},
						},
					},
					Work: []ResumeWork{
						{
							Name:      "Meshwork Labs",
							Position:  "Senior Backend Engineer",
							URL:       "https://meshwork.example.com",
							StartDate: "2022-03",
							Summary:   "Lead engineer on the federation team, building ActivityPub-based services in Go.",
							Highlights: []string{
								"Designed and built a federated messaging system handling 50k messages/day",
								"Implemented HTTP Signatures and WebFinger for inter-server authentication",
							},
						},
						{
							Name:      "DataFlow Inc.",
							Position:  "Backend Engineer",
							StartDate: "2018-06",
							EndDate:   "2022-02",
							Summary:   "Built data pipeline services and REST APIs.",
							Highlights: []string{
								"Reduced API response latency by 40% through caching redesign",
							},
						},
					},
					Education: []ResumeEducation{
						{
							Institution: "University of British Columbia",
							Area:        "Computer Science",
							StudyType:   "B.Sc.",
							StartDate:   "2012-09",
							EndDate:     "2016-05",
						},
					},
					Skills: []ResumeSkill{
						{Name: "Go", Level: "Expert", Keywords: []string{"concurrency", "net/http", "testing"}},
						{Name: "Distributed Systems", Level: "Advanced", Keywords: []string{"ActivityPub", "federation", "consensus"}},
						{Name: "Databases", Level: "Advanced", Keywords: []string{"PostgreSQL", "Redis", "SQLite"}},
					},
					Languages: []ResumeLanguage{
						{Language: "English", Fluency: "Native"},
						{Language: "Mandarin", Fluency: "Fluent"},
					},
					Projects: []ResumeProject{
						{
							Name:        "go-activitypub",
							Description: "Open-source Go library for ActivityPub server implementations.",
							URL:         "https://github.com/alicezhang/go-activitypub",
							Roles:       []string{"Creator", "Maintainer"},
							Keywords:    []string{"Go", "ActivityPub", "Federation"},
						},
					},
				},
			},
			Messages: []ChatMessage{
				{FromMe: false, Text: "Hey! Great meeting you at GopherCon. Love your talk on federation.", Timestamp: "2025-11-15 14:32"},
				{FromMe: true, Text: "Thanks Alice! Your work on ActivityPub at Meshwork is really impressive.", Timestamp: "2025-11-15 14:35"},
				{FromMe: false, Text: "We should chat more about the WebFinger stuff you mentioned.", Timestamp: "2025-11-15 14:37"},
				{FromMe: true, Text: "Definitely. I've been working on a Go library for it.", Timestamp: "2025-11-15 14:40"},
				{FromMe: false, Text: "Oh nice! Send me the repo when you get a chance.", Timestamp: "2025-11-15 14:41"},
			},
		},
		{
			Name:    "Bob Okafor",
			Title:   "Product Manager",
			Company: "FediCorp",
			FediID:  "@bob@fedicorp.example",
			Note:    "Met at Fediverse Developer Summit. Working on decentralized identity.",
			Resumes: []Resume{},
			Messages: []ChatMessage{
				{FromMe: false, Text: "Hey, good to connect! Are you going to the Fediverse summit next month?", Timestamp: "2025-12-01 09:15"},
				{FromMe: true, Text: "Wouldn't miss it. Are you presenting?", Timestamp: "2025-12-01 10:22"},
				{FromMe: false, Text: "Yeah, doing a session on decentralized identity. Would love your feedback on the slides.", Timestamp: "2025-12-01 10:30"},
				{FromMe: true, Text: "Happy to take a look. Send them over whenever.", Timestamp: "2025-12-01 10:33"},
			},
		},
		{
			Name:    "Carol Reyes",
			Title:   "Freelance UX Designer",
			Company: "",
			Note:    "Met at Vancouver Tech Meetup. Does contract work for startups.",
			Resumes: []Resume{
				{
					Label: "UX Design Portfolio",
					Basics: ResumeBasics{
						Name:    "Carol Reyes",
						Label:   "UX Designer",
						Email:   "carol@example.com",
						URL:     "https://carolreyes.design",
						Summary: "UX designer with 8 years of experience in web and mobile. Specializes in B2B SaaS products.",
						Location: ResumeLocation{
							City:        "Vancouver",
							CountryCode: "CA",
							Region:      "British Columbia",
						},
					},
					Work: []ResumeWork{
						{
							Name:      "Various Startups",
							Position:  "Freelance UX Designer",
							StartDate: "2019-01",
							Summary:   "Contract UX design for early-stage startups in fintech and healthtech.",
							Highlights: []string{
								"Redesigned onboarding flow for a fintech app, improving completion rate by 35%",
								"Created design system used across 3 product teams",
							},
						},
					},
					Skills: []ResumeSkill{
						{Name: "UX Design", Level: "Expert", Keywords: []string{"Figma", "user research", "prototyping"}},
						{Name: "Frontend", Level: "Intermediate", Keywords: []string{"HTML", "CSS", "React"}},
					},
				},
				{
					Label: "Consulting CV",
					Basics: ResumeBasics{
						Name:    "Carol Reyes",
						Label:   "Independent UX Consultant",
						Email:   "carol@example.com",
						Summary: "Independent UX consultant since 2021. Clients include fintech and healthtech companies.",
					},
					Work: []ResumeWork{
						{
							Name:      "Self-Employed",
							Position:  "UX Consultant",
							StartDate: "2021-01",
							Summary:   "Providing UX strategy, audits, and design services to B2B SaaS companies.",
						},
					},
					Certificates: []ResumeCertificate{
						{Name: "Google UX Design Certificate", Date: "2021-06", Issuer: "Google"},
					},
					References: []ResumeReference{
						{Name: "Jordan Lee", Reference: "Carol transformed our product experience. Her research-driven approach led to measurable improvements in user retention."},
					},
				},
			},
		},
		{
			Name:     "David Park",
			Title:    "CTO",
			Company:  "OpenRelay Inc.",
			FediID:   "@david@openrelay.social",
			Favorite: true,
			Note:     "Spoke on a panel about federation protocols. Very knowledgeable about Mastodon internals.",
			Resumes: []Resume{
				{
					Label: "Executive Resume",
					Basics: ResumeBasics{
						Name:    "David Park",
						Label:   "CTO & Engineering Leader",
						Email:   "david@openrelay.example.com",
						URL:     "https://davidpark.io",
						Summary: "Engineering leader with 15 years of experience. Founded two startups. Deep expertise in federation protocols and open-source infrastructure.",
						Location: ResumeLocation{
							City:        "Toronto",
							CountryCode: "CA",
							Region:      "Ontario",
						},
						Profiles: []ResumeProfile{
							{Network: "Fediverse", Username: "@david@openrelay.social", URL: "https://openrelay.social/@david"},
						},
					},
					Work: []ResumeWork{
						{
							Name:      "OpenRelay Inc.",
							Position:  "CTO",
							StartDate: "2020-01",
							Summary:   "Leading engineering for an open-source federation relay service.",
							Highlights: []string{
								"Grew engineering team from 2 to 14",
								"Architected relay system handling 1M+ activities/day",
							},
						},
						{
							Name:      "SocialMesh (acquired)",
							Position:  "Co-Founder & CTO",
							StartDate: "2016-03",
							EndDate:   "2019-12",
							Summary:   "Built a decentralized social networking platform.",
						},
					},
					Education: []ResumeEducation{
						{
							Institution: "University of Toronto",
							Area:        "Software Engineering",
							StudyType:   "M.Eng.",
							StartDate:   "2009-09",
							EndDate:     "2011-05",
						},
					},
					Awards: []ResumeAward{
						{Title: "Top 40 Under 40 in Canadian Tech", Date: "2023", Awarder: "TechTO"},
					},
					Publications: []ResumePublication{
						{
							Name:        "Scaling Federation: Lessons from Building a Relay Network",
							Publisher:   "ACM Queue",
							ReleaseDate: "2024-03",
							Summary:     "Practical lessons on scaling ActivityPub relay infrastructure.",
						},
					},
					Skills: []ResumeSkill{
						{Name: "Engineering Leadership", Level: "Expert", Keywords: []string{"team building", "architecture", "strategy"}},
						{Name: "Federation Protocols", Level: "Expert", Keywords: []string{"ActivityPub", "Mastodon", "relays"}},
						{Name: "Go", Level: "Advanced", Keywords: []string{"microservices", "performance"}},
					},
					Languages: []ResumeLanguage{
						{Language: "English", Fluency: "Native"},
						{Language: "Korean", Fluency: "Conversational"},
					},
				},
			},
		},
		{
			Name:    "Fatima Al-Rashid",
			Title:   "DevRel Engineer",
			Company: "Meshwork Labs",
			FediID:  "@fatima@meshwork.social",
			Note:    "Gave a great talk on onboarding developers to the Fediverse.",
			Resumes: []Resume{},
		},
	}
}

func dummyGroups() []Group {
	return []Group{
		{
			Name:     "Meshwork Labs + FediCorp",
			Members:  []string{"Alice Zhang", "Bob Okafor", "Fatima Al-Rashid"},
			Favorite: true,
			Messages: []ChatMessage{
				{Sender: "Alice Zhang", Text: "Hey everyone, should we sync on the federation API this week?", Timestamp: "2025-12-10 09:00"},
				{Sender: "Bob Okafor", Text: "Works for me. Thursday afternoon?", Timestamp: "2025-12-10 09:15"},
				{Sender: "Fatima Al-Rashid", Text: "Thursday is good. I can share the developer onboarding doc beforehand.", Timestamp: "2025-12-10 09:22"},
				{FromMe: true, Text: "Thursday works. Looking forward to it.", Timestamp: "2025-12-10 09:30"},
				{Sender: "Alice Zhang", Text: "Great, I'll send a calendar invite.", Timestamp: "2025-12-10 09:35"},
			},
		},
		{
			Name:    "Vancouver Tech Meetup",
			Members: []string{"Carol Reyes", "David Park"},
			Messages: []ChatMessage{
				{Sender: "David Park", Text: "Anyone going to the meetup next Tuesday?", Timestamp: "2025-12-05 18:00"},
				{Sender: "Carol Reyes", Text: "I'll be there! Hoping to see the talk on Wasm.", Timestamp: "2025-12-05 18:10"},
				{FromMe: true, Text: "Count me in.", Timestamp: "2025-12-05 18:15"},
			},
		},
	}
}
