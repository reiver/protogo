package gui

type Person struct {
	DBID        int64
	Name        string
	Title       string
	Company     string
	FediID      string
	Note        string
	SummaryHTML string
	IconURL     string
	BannerURL   string
	ProfileURL  string
	Favorite    bool
	Resumes     []Resume
	Messages    []ChatMessage
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

type Gig struct {
	Title       string
	Company     string
	Location    string
	Type        string // "Full-time", "Contract", "Freelance", etc.
	Description string
	PostedBy    string
	Timestamp   string
}

type Group struct {
	DBID     int64
	Name     string
	Members  []string
	Favorite bool
	Messages []ChatMessage
}
