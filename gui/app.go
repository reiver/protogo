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

	people []Person
	groups []Group

	personClicks []widget.Clickable
	groupClicks  []widget.Clickable
	backClick    widget.Clickable

	homeList layout.List
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
	}
}

func (receiver *App) Layout(gtx layout.Context) layout.Dimensions {
	switch receiver.page {
	case PagePersonDetail:
		return receiver.layoutPersonDetail(gtx)
	case PageGroupDetail:
		return receiver.layoutGroupDetail(gtx)
	default:
		return receiver.layoutHome(gtx)
	}
}
