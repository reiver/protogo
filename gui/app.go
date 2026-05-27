package gui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type App struct {
	theme *material.Theme

	page           Page
	selectedPerson int
	selectedGroup  int
	selectedResume int

	people []Person
	groups []Group

	personClicks []widget.Clickable
	groupClicks  []widget.Clickable
	resumeClicks []widget.Clickable
	backClick    widget.Clickable
	chatClick    widget.Clickable
	sendClick    widget.Clickable
	chatEditor   widget.Editor

	homeList   layout.List
	resumeList layout.List
	chatList   layout.List
}

func newApp() *App {
	var people []Person = dummyPeople()
	var groups []Group = dummyGroups()

	return &App{
		theme: material.NewTheme(),

		page: PageHome,

		people: people,
		groups: groups,

		personClicks: make([]widget.Clickable, len(people)),
		groupClicks:  make([]widget.Clickable, len(groups)),

		homeList: layout.List{
			Axis: layout.Vertical,
		},
		resumeList: layout.List{
			Axis: layout.Vertical,
		},
		chatEditor: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		chatList: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.End,
		},
	}
}

func (receiver *App) Layout(gtx layout.Context) layout.Dimensions {
	switch receiver.page {
	case PagePersonDetail:
		return receiver.layoutPersonDetail(gtx)
	case PageGroupDetail:
		return receiver.layoutGroupDetail(gtx)
	case PageResumeDetail:
		return receiver.layoutResumeDetail(gtx)
	case PageChat:
		return receiver.layoutChat(gtx)
	default:
		return receiver.layoutHome(gtx)
	}
}
