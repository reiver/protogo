package gui

type Person struct {
	Name    string
	Title   string
	Company string
	Note    string
	Resumes []Resume
}

type Resume struct {
	Label   string
	Content string
}

type Group struct {
	Name    string
	Members []string
}

func dummyPeople() []Person {
	return []Person{
		{
			Name:    "Alice Zhang",
			Title:   "Senior Backend Engineer",
			Company: "Meshwork Labs",
			Note:    "Met at GopherCon 2025. Interested in distributed systems and ActivityPub.",
			Resumes: []Resume{
				{
					Label:   "Software Engineer Resume",
					Content: "10 years experience in Go, distributed systems, and API design.",
				},
			},
		},
		{
			Name:    "Bob Okafor",
			Title:   "Product Manager",
			Company: "FediCorp",
			Note:    "Met at Fediverse Developer Summit. Working on decentralized identity.",
			Resumes: []Resume{},
		},
		{
			Name:    "Carol Reyes",
			Title:   "Freelance UX Designer",
			Company: "",
			Note:    "Met at Vancouver Tech Meetup. Does contract work for startups.",
			Resumes: []Resume{
				{
					Label:   "UX Design Portfolio",
					Content: "8 years of UX design for web and mobile. Specializes in B2B SaaS.",
				},
				{
					Label:   "Consulting CV",
					Content: "Independent UX consultant since 2021. Clients include fintech and healthtech.",
				},
			},
		},
		{
			Name:    "David Park",
			Title:   "CTO",
			Company: "OpenRelay Inc.",
			Note:    "Spoke on a panel about federation protocols. Very knowledgeable about Mastodon internals.",
			Resumes: []Resume{
				{
					Label:   "Executive Resume",
					Content: "15 years in engineering leadership. Founded two startups.",
				},
			},
		},
		{
			Name:    "Fatima Al-Rashid",
			Title:   "DevRel Engineer",
			Company: "Meshwork Labs",
			Note:    "Gave a great talk on onboarding developers to the Fediverse.",
			Resumes: []Resume{},
		},
	}
}

func dummyGroups() []Group {
	return []Group{
		{
			Name:    "Meshwork Labs + FediCorp",
			Members: []string{"Alice Zhang", "Bob Okafor", "Fatima Al-Rashid"},
		},
		{
			Name:    "Vancouver Tech Meetup",
			Members: []string{"Carol Reyes", "David Park"},
		},
	}
}
