package gui

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var resumeSectionColor color.NRGBA = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
var resumeLabelColor color.NRGBA = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}

func (receiver *App) layoutResumeDetail(gtx layout.Context) layout.Dimensions {
	if receiver.backClick.Clicked(gtx) {
		receiver.page = PagePersonDetail
	}

	if receiver.selectedPerson < 0 || receiver.selectedPerson >= len(receiver.people) {
		receiver.page = PageHome
		return layout.Dimensions{}
	}

	var person Person = receiver.people[receiver.selectedPerson]

	if receiver.selectedResume < 0 || receiver.selectedResume >= len(person.Resumes) {
		receiver.page = PagePersonDetail
		return layout.Dimensions{}
	}

	var resume Resume = person.Resumes[receiver.selectedResume]

	var widgets []layout.Widget

	widgets = appendResumeBasics(widgets, receiver.theme, resume.Basics)
	widgets = appendResumeWork(widgets, receiver.theme, resume.Work)
	widgets = appendResumeVolunteer(widgets, receiver.theme, resume.Volunteer)
	widgets = appendResumeEducation(widgets, receiver.theme, resume.Education)
	widgets = appendResumeAwards(widgets, receiver.theme, resume.Awards)
	widgets = appendResumeCertificates(widgets, receiver.theme, resume.Certificates)
	widgets = appendResumePublications(widgets, receiver.theme, resume.Publications)
	widgets = appendResumeSkills(widgets, receiver.theme, resume.Skills)
	widgets = appendResumeLanguages(widgets, receiver.theme, resume.Languages)
	widgets = appendResumeInterests(widgets, receiver.theme, resume.Interests)
	widgets = appendResumeReferences(widgets, receiver.theme, resume.References)
	widgets = appendResumeProjects(widgets, receiver.theme, resume.Projects)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, resume.Label, true)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.resumeList.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
				return widgets[index](gtx)
			})
		}),
	)
}

func resumeSectionHeader(th *material.Theme, title string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			label := material.Subtitle1(th, title)
			label.Color = resumeSectionColor
			return label.Layout(gtx)
		})
	}
}

func resumeCardItem(th *material.Theme, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layoutCard(gtx, w)
		})
	}
}

// --- Basics ---

func appendResumeBasics(widgets []layout.Widget, th *material.Theme, basics ResumeBasics) []layout.Widget {
	if "" == basics.Name && "" == basics.Summary {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Basics"))
	widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
		var items []layout.FlexChild

		if "" != basics.Name {
			items = append(items, resumeField(th, "Name", basics.Name))
		}
		if "" != basics.Label {
			items = append(items, resumeField(th, "Title", basics.Label))
		}
		if "" != basics.Email {
			items = append(items, resumeField(th, "Email", basics.Email))
		}
		if "" != basics.Phone {
			items = append(items, resumeField(th, "Phone", basics.Phone))
		}
		if "" != basics.URL {
			items = append(items, resumeField(th, "Website", basics.URL))
		}
		if "" != basics.Summary {
			items = append(items, resumeField(th, "Summary", basics.Summary))
		}

		var locationParts []string
		if "" != basics.Location.City {
			locationParts = append(locationParts, basics.Location.City)
		}
		if "" != basics.Location.Region {
			locationParts = append(locationParts, basics.Location.Region)
		}
		if "" != basics.Location.CountryCode {
			locationParts = append(locationParts, basics.Location.CountryCode)
		}
		if 0 < len(locationParts) {
			items = append(items, resumeField(th, "Location", strings.Join(locationParts, ", ")))
		}

		for _, profile := range basics.Profiles {
			var value string
			if "" != profile.URL {
				value = fmt.Sprintf("%s (%s)", profile.Username, profile.URL)
			} else {
				value = profile.Username
			}
			items = append(items, resumeField(th, profile.Network, value))
		}

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
	}))

	return widgets
}

// --- Work ---

func appendResumeWork(widgets []layout.Widget, th *material.Theme, work []ResumeWork) []layout.Widget {
	if 0 == len(work) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Work"))

	for i := range work {
		var entry ResumeWork = work[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			var title string = entry.Position
			if "" != entry.Name {
				title = fmt.Sprintf("%s — %s", entry.Position, entry.Name)
			}
			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, title).Layout(gtx)
			}))
			items = append(items, resumeDateRange(th, entry.StartDate, entry.EndDate))

			if "" != entry.Summary {
				items = append(items, resumeField(th, "", entry.Summary))
			}
			for _, h := range entry.Highlights {
				items = append(items, resumeBullet(th, h))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- Volunteer ---

func appendResumeVolunteer(widgets []layout.Widget, th *material.Theme, volunteer []ResumeVolunteer) []layout.Widget {
	if 0 == len(volunteer) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Volunteer"))

	for i := range volunteer {
		var entry ResumeVolunteer = volunteer[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			var title string = entry.Position
			if "" != entry.Organization {
				title = fmt.Sprintf("%s — %s", entry.Position, entry.Organization)
			}
			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, title).Layout(gtx)
			}))
			items = append(items, resumeDateRange(th, entry.StartDate, entry.EndDate))

			if "" != entry.Summary {
				items = append(items, resumeField(th, "", entry.Summary))
			}
			for _, h := range entry.Highlights {
				items = append(items, resumeBullet(th, h))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- Education ---

func appendResumeEducation(widgets []layout.Widget, th *material.Theme, education []ResumeEducation) []layout.Widget {
	if 0 == len(education) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Education"))

	for i := range education {
		var entry ResumeEducation = education[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			var title string = entry.Institution
			if "" != entry.StudyType && "" != entry.Area {
				title = fmt.Sprintf("%s — %s, %s", entry.Institution, entry.StudyType, entry.Area)
			} else if "" != entry.Area {
				title = fmt.Sprintf("%s — %s", entry.Institution, entry.Area)
			}
			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, title).Layout(gtx)
			}))
			items = append(items, resumeDateRange(th, entry.StartDate, entry.EndDate))

			if "" != entry.Score {
				items = append(items, resumeField(th, "Score", entry.Score))
			}
			for _, c := range entry.Courses {
				items = append(items, resumeBullet(th, c))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- Awards ---

func appendResumeAwards(widgets []layout.Widget, th *material.Theme, awards []ResumeAward) []layout.Widget {
	if 0 == len(awards) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Awards"))

	for i := range awards {
		var entry ResumeAward = awards[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, entry.Title).Layout(gtx)
			}))

			if "" != entry.Awarder || "" != entry.Date {
				var detail string
				if "" != entry.Awarder && "" != entry.Date {
					detail = fmt.Sprintf("%s — %s", entry.Awarder, entry.Date)
				} else if "" != entry.Awarder {
					detail = entry.Awarder
				} else {
					detail = entry.Date
				}
				items = append(items, resumeFieldGray(th, detail))
			}

			if "" != entry.Summary {
				items = append(items, resumeField(th, "", entry.Summary))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- Certificates ---

func appendResumeCertificates(widgets []layout.Widget, th *material.Theme, certs []ResumeCertificate) []layout.Widget {
	if 0 == len(certs) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Certificates"))

	for i := range certs {
		var entry ResumeCertificate = certs[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, entry.Name).Layout(gtx)
			}))

			if "" != entry.Issuer {
				items = append(items, resumeField(th, "Issuer", entry.Issuer))
			}
			if "" != entry.Date {
				items = append(items, resumeField(th, "Date", entry.Date))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- Publications ---

func appendResumePublications(widgets []layout.Widget, th *material.Theme, pubs []ResumePublication) []layout.Widget {
	if 0 == len(pubs) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Publications"))

	for i := range pubs {
		var entry ResumePublication = pubs[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, entry.Name).Layout(gtx)
			}))

			if "" != entry.Publisher || "" != entry.ReleaseDate {
				var detail string
				if "" != entry.Publisher && "" != entry.ReleaseDate {
					detail = fmt.Sprintf("%s — %s", entry.Publisher, entry.ReleaseDate)
				} else if "" != entry.Publisher {
					detail = entry.Publisher
				} else {
					detail = entry.ReleaseDate
				}
				items = append(items, resumeFieldGray(th, detail))
			}

			if "" != entry.Summary {
				items = append(items, resumeField(th, "", entry.Summary))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- Skills ---

func appendResumeSkills(widgets []layout.Widget, th *material.Theme, skills []ResumeSkill) []layout.Widget {
	if 0 == len(skills) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Skills"))

	for i := range skills {
		var entry ResumeSkill = skills[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			var title string = entry.Name
			if "" != entry.Level {
				title = fmt.Sprintf("%s — %s", entry.Name, entry.Level)
			}
			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, title).Layout(gtx)
			}))

			if 0 < len(entry.Keywords) {
				items = append(items, resumeFieldGray(th, strings.Join(entry.Keywords, ", ")))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- Languages ---

func appendResumeLanguages(widgets []layout.Widget, th *material.Theme, languages []ResumeLanguage) []layout.Widget {
	if 0 == len(languages) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Languages"))
	widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
		var items []layout.FlexChild
		for _, lang := range languages {
			var value string = lang.Language
			if "" != lang.Fluency {
				value = fmt.Sprintf("%s — %s", lang.Language, lang.Fluency)
			}
			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body2(th, value).Layout(gtx)
			}))
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
	}))

	return widgets
}

// --- Interests ---

func appendResumeInterests(widgets []layout.Widget, th *material.Theme, interests []ResumeInterest) []layout.Widget {
	if 0 == len(interests) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Interests"))

	for i := range interests {
		var entry ResumeInterest = interests[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, entry.Name).Layout(gtx)
			}))

			if 0 < len(entry.Keywords) {
				items = append(items, resumeFieldGray(th, strings.Join(entry.Keywords, ", ")))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- References ---

func appendResumeReferences(widgets []layout.Widget, th *material.Theme, refs []ResumeReference) []layout.Widget {
	if 0 == len(refs) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "References"))

	for i := range refs {
		var entry ResumeReference = refs[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, entry.Name).Layout(gtx)
			}))

			if "" != entry.Reference {
				items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, fmt.Sprintf("%q", entry.Reference))
					lbl.Color = color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xFF}
					return lbl.Layout(gtx)
				}))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- Projects ---

func appendResumeProjects(widgets []layout.Widget, th *material.Theme, projects []ResumeProject) []layout.Widget {
	if 0 == len(projects) {
		return widgets
	}

	widgets = append(widgets, resumeSectionHeader(th, "Projects"))

	for i := range projects {
		var entry ResumeProject = projects[i]
		widgets = append(widgets, resumeCardItem(th, func(gtx layout.Context) layout.Dimensions {
			var items []layout.FlexChild

			items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Subtitle2(th, entry.Name).Layout(gtx)
			}))

			if "" != entry.StartDate || "" != entry.EndDate {
				items = append(items, resumeDateRange(th, entry.StartDate, entry.EndDate))
			}

			if "" != entry.Description {
				items = append(items, resumeField(th, "", entry.Description))
			}

			if 0 < len(entry.Roles) {
				items = append(items, resumeField(th, "Roles", strings.Join(entry.Roles, ", ")))
			}

			if "" != entry.URL {
				items = append(items, resumeField(th, "URL", entry.URL))
			}

			for _, h := range entry.Highlights {
				items = append(items, resumeBullet(th, h))
			}

			if 0 < len(entry.Keywords) {
				items = append(items, resumeFieldGray(th, strings.Join(entry.Keywords, ", ")))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
		}))
	}

	return widgets
}

// --- Helpers ---

func resumeField(th *material.Theme, label string, value string) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if "" == label {
				return material.Body2(th, value).Layout(gtx)
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Caption(th, label)
					lbl.Color = resumeLabelColor
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Body2(th, value).Layout(gtx)
				}),
			)
		})
	})
}

func resumeFieldGray(th *material.Theme, value string) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(th, value)
			lbl.Color = resumeLabelColor
			return lbl.Layout(gtx)
		})
	})
}

func resumeDateRange(th *material.Theme, startDate string, endDate string) layout.FlexChild {
	var value string
	switch {
	case "" != startDate && "" != endDate:
		value = fmt.Sprintf("%s — %s", startDate, endDate)
	case "" != startDate:
		value = fmt.Sprintf("%s — Present", startDate)
	case "" != endDate:
		value = endDate
	default:
		return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		})
	}
	return resumeFieldGray(th, value)
}

func resumeBullet(th *material.Theme, text string) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(1), Bottom: unit.Dp(1), Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(th, fmt.Sprintf("- %s", text))
			lbl.Color = color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xFF}
			return lbl.Layout(gtx)
		})
	})
}
